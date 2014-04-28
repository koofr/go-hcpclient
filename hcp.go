package hcpclient

import (
	"fmt"
	"git.koofr.lan/go-httpclient.git"
	"git.koofr.lan/go-ioutils.git"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Hcp struct {
	HTTPClient *httpclient.HTTPClient
	Prefix     string
}

func NewHcp(endpoint string, user string, key string, prefix string) (hcp *Hcp, err error) {
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}

	if prefix != "" {
		if !strings.HasSuffix(prefix, "/") {
			prefix += "/"
		}

		endpoint += prefix
	}

	u, err := url.Parse(endpoint)

	if err != nil {
		return
	}

	client := httpclient.Insecure()
	client.BaseURL = u
	client.Headers.Set("Authorization", AuthHeader(user, key))

	hcp = &Hcp{
		HTTPClient: client,
		Prefix:     prefix,
	}

	return
}

func (h *Hcp) Request(req *httpclient.RequestData) (response *http.Response, err error) {
	return h.HTTPClient.Request(req)
}

func (h *Hcp) Path(path string) string {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}

	return path
}

func (h *Hcp) ObjectInfo(path string) (info *HcpObject, err error) {
	res, err := h.Request(&httpclient.RequestData{
		Method:         "HEAD",
		Path:           h.Path(path),
		ExpectedStatus: []int{http.StatusOK},
		RespConsume:    true,
	})

	if err != nil {
		return
	}

	info = HcpObjectFromHeaders(path, res.Header)

	return
}

func (h *Hcp) GetObject(path string, span *ioutils.FileSpan) (obj *HcpObject, err error) {
	req := httpclient.RequestData{
		Method:         "GET",
		Path:           h.Path(path),
		ExpectedStatus: []int{http.StatusOK, http.StatusPartialContent},
	}

	if span != nil {
		req.Headers = make(http.Header)
		req.Headers.Set("Range", fmt.Sprintf("bytes=%d-%d", span.Start, span.End))
	}

	res, err := h.Request(&req)

	if err != nil {
		return
	}

	obj = HcpObjectFromHeaders(path, res.Header)
	obj.Reader = res.Body

	return
}

func (h *Hcp) PutObject(path string, reader io.Reader) (err error) {
	_, err = h.Request(&httpclient.RequestData{
		Method:         "PUT",
		Path:           h.Path(path),
		ReqReader:      reader,
		ExpectedStatus: []int{http.StatusCreated},
		RespConsume:    true,
	})

	return
}

func (h *Hcp) CreateDirectory(path string) (err error) {
	params := make(url.Values)
	params.Set("type", "directory")

	_, err = h.Request(&httpclient.RequestData{
		Method:         "PUT",
		Path:           h.Path(path),
		Params:         params,
		ExpectedStatus: []int{http.StatusCreated},
		RespConsume:    true,
	})

	return
}

func (h *Hcp) CopyObject(fromPath string, toPath string) (err error) {
	host := h.HTTPClient.BaseURL.Host
	hostParts := strings.Split(host, ".")
	namespace := hostParts[0]
	tenant := hostParts[1]

	// Yes, that's right! Double escape!
	escapedFromPath := httpclient.EscapePath(fromPath)
	escapedFromPath = httpclient.EscapePath(escapedFromPath)

	copySource := fmt.Sprintf("%s.%s/%s%s", namespace, tenant, h.Prefix, escapedFromPath)

	headers := make(http.Header)
	headers.Set("X-HCP-CopySource", copySource)

	_, err = h.Request(&httpclient.RequestData{
		Method:         "PUT",
		Path:           h.Path(toPath),
		Headers:        headers,
		ExpectedStatus: []int{http.StatusOK},
		RespConsume:    true,
	})

	return
}

func (h *Hcp) DeleteObject(path string) (err error) {
	_, err = h.Request(&httpclient.RequestData{
		Method:         "DELETE",
		Path:           h.Path(path),
		ExpectedStatus: []int{http.StatusOK},
		RespConsume:    true,
	})

	return
}
