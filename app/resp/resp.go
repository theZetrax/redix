package resp

const (
	CRLF = "\r\n"
)

type RespHandler interface {
	Process() []byte
}
type Resp struct {
	Data []byte
}

func HandleResp(data []byte) RespHandler {
	if len(data) == 0 {
		return NewSimpleString(data)
	}

	switch data[0] {
	case '+':
		return NewSimpleString(data)
	default:
		return NewSimpleString(data)
	}
}
