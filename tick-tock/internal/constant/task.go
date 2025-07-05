package constant

type TaskStatus int32

const (
	TaskInit TaskStatus = iota
	TaskIng
	TaskSuccess
	TaskFail
)

func (s TaskStatus) ToInt32() int32 {
	return int32(s)
}
