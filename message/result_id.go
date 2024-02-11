package message

type ResultId int32

const (
	Result_Success ResultId = 0

	Result_AlreaedRegistered ResultId = 1

	Result_SmartItemNotFound ResultId = 1001
)
