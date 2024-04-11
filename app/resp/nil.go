package resp

func EncodeNil() []byte {
	return []byte("$-1" + CRLF)
}
