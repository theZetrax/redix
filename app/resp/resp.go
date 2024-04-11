package resp

const (
	CRLF = "\r\n"
)

type RespHandler interface {
	Process() []byte
	String() string
}
type Resp struct {
	Data []byte
}

func HandleResp(data []byte) (resp RespHandler, size int) {
	if len(data) == 0 {
		return NewSimpleString(data), -1
	}

	switch data[0] {
	case '+':
		return NewSimpleString(data), -1
	case '*':
		return NewArray(data), GetArraySize(data)
	case '$':
		return NewBulkString(data), GetBulkStringSize(data)
	default:
		return NewSimpleString(data), -1
	}
}
