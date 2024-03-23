package internal

import "strings"

const (
	CLRF = "\r\n"
)

type Request struct {
	Method      string
	HttpVersion string
	Headers     map[string]string
	Body        string
	Url         string
}

func ParseRequest(buffer []byte) Request {
	raw := string(buffer)
	meta_headers := strings.Split(raw, CLRF+CLRF)[0]
	body := strings.TrimSpace(strings.Split(raw, CLRF+CLRF)[0])

	meta := strings.Split(strings.Split(meta_headers, CLRF)[0], " ")
	headers_raw := strings.Split(meta_headers, CLRF)[1:]
	headers := make(map[string]string, 0)
	for _, h := range headers_raw {
		header := strings.Split(h, ":")
		if len(header) != 2 {
			continue
		}
	}

	return Request{
		Method:      meta[0],
		Url:         meta[1],
		HttpVersion: meta[2],
		Headers:     headers,
		Body:        body,
	}
}
