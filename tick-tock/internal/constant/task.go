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

const (
	BucketPrefixFormat = "2006-01-02 15:04:05"
)
