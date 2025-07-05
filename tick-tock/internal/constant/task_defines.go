package constant

type TaskDefineStatus int32

const (
	TaskDefineUnknow TaskDefineStatus = iota
	TaskDefineActive
	TaskDefineInactive
)

func (s TaskDefineStatus) ToInt32() int32 {
	return int32(s)
}
