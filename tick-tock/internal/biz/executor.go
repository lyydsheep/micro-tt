package biz

import (
	"context"
	"encoding/json"
	"tick-tock/internal/conf"
	"tick-tock/internal/constant"
	"tick-tock/pkg/log"
	"time"
)

type ExecutorUsecase struct {
	conf           *conf.Data
	taskDefineRepo TaskDefineRepo
	taskRepo       TaskRepo
}

func NewExecutorUsecase(conf *conf.Data, taskDefineRepo TaskDefineRepo, taskRepo TaskRepo) *ExecutorUsecase {
	return &ExecutorUsecase{
		conf:           conf,
		taskDefineRepo: taskDefineRepo,
		taskRepo:       taskRepo,
	}
}

func (uc *ExecutorUsecase) Work(ctx context.Context, task *Task) {
	// 查询 task_define，获取具体执行的动作
	taskDefine, err := uc.taskDefineRepo.GetTaskDefineByTID(ctx, task.Tid)
	if err != nil {
		log.Error(ctx, "get task define error.", "error", err, "tID", task.Tid)
		uc.fail(ctx, task.Tid, task.RunTime)
		return
	}
	// 执行
	costTime, outPut, err := uc.execute(ctx, taskDefine)
	if err != nil {
		log.Error(ctx, "execute task error.", "error", err, "tID", task.Tid)
		uc.fail(ctx, task.Tid, task.RunTime)
		return
	}

	// 保存结果并设置任务状态
	uc.success(ctx, task.Tid, task.RunTime, costTime, outPut)
}

func (uc *ExecutorUsecase) execute(ctx context.Context, taskDefine *TaskDefine) (costTime int64, outPut string, err error) {
	// 执行
	var notifyHTTPParam NotifyHTTPParam
	if err := json.Unmarshal([]byte(taskDefine.NotifyHTTPParam), &notifyHTTPParam); err != nil {
		log.Error(ctx, "unmarshal notify http param error.", "error", err, "tID", taskDefine.Tid, "notifyParams", taskDefine.NotifyHTTPParam)
		return 0, "", err
	}
	startTime := time.Now()
	var resp map[string]any
	if err := Do(ctx, notifyHTTPParam.URL, notifyHTTPParam.Method, notifyHTTPParam.Body, notifyHTTPParam.Headers, &resp); err != nil {
		log.Error(ctx, "do request error.", "error", err, "tID", taskDefine.Tid, "url", notifyHTTPParam.URL, "method", notifyHTTPParam.Method, "headers", notifyHTTPParam.Headers, "body", notifyHTTPParam.Body)
		return 0, "", err
	}

	// 序列化 resp
	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Error(ctx, "marshal response error.", "error", err, "tID", taskDefine.Tid, "resp", resp)
		return 0, "", err
	}

	return time.Now().Sub(startTime).Milliseconds(), string(bytes), nil
}

func (uc *ExecutorUsecase) fail(ctx context.Context, tID string, runTime time.Time) {
	if err := uc.taskRepo.UpdateByTIDAndRuntime(ctx, tID, runTime, map[string]any{"status": constant.TaskFail.ToInt32(), "update_time": time.Now().UTC()}); err != nil {
		log.Error(ctx, "update task status error.", "error", err, "tID", tID, "runTime", runTime)
	}
}

func (uc *ExecutorUsecase) success(ctx context.Context, tID string, runTime time.Time, costTime int64, outPut string) {
	// 更新任务状态
	if err := uc.taskRepo.UpdateByTIDAndRuntime(ctx, tID, runTime, map[string]any{"status": constant.TaskSuccess.ToInt32(),
		"cost_time": costTime, "output": outPut, "update_time": time.Now().UTC()}); err != nil {
		log.Error(ctx, "update task status error.", "error", err, "tID", tID, "runTime", runTime)
	}
}
