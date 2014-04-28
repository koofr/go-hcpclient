package hcpclient

import (
	"io"
	"net/http"
	"strconv"
	"time"
)

type HcpObject struct {
	Name     string
	Length   int64
	Modified *time.Time
	Reader   io.ReadCloser
}

func HcpObjectFromHeaders(path string, headers http.Header) *HcpObject {
	contentLength, _ := strconv.ParseInt(headers.Get("Content-Length"), 10, 0)
	m, _ := strconv.ParseFloat(headers.Get("X-HCP-ChangeTimeMilliseconds"), 64)
	modified := time.Unix(0, int64(m)*1000*1000)

	return &HcpObject{
		Name:     path,
		Length:   contentLength,
		Modified: &modified,
	}
}
