-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS task_defines
(
    id                bigint AUTO_INCREMENT PRIMARY KEY,
    tid               varchar(36) NOT NULL COMMENT '任务唯一标识（UUID）',
    app               varchar(64) NOT NULL DEFAULT '' COMMENT '所属应用标识',
    name              varchar(32) NOT NULL DEFAULT '' COMMENT '任务名称',
    status            tinyint     NOT NULL DEFAULT 1 COMMENT '状态：1-active, 2-inactive',
    cron              varchar(64) NOT NULL COMMENT 'Cron表达式（标准5/6字段）',
    notify_http_param json NULL COMMENT '回调参数（JSON格式，如{"url":"","method":"POST","headers":{}}）',
    create_time       timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    update_time       timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_task_defines_tid (tid)
) COMMENT '任务定义表';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE task_defines;
-- +goose StatementEnd
