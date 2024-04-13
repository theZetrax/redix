package resp

func EncodeSimpleError(message string) []byte {
	return []byte("-" + message + CRLF)
}
