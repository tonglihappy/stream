package util

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type RespChecker interface {
	AssertRespHeaderCheck(t *testing.T, resp *http.Response)
	AssertRespBodyCheck(t *testing.T, resp *http.Response)
}

type HttpReq struct {
	Url       string
	Method    string
	Headers   map[string]string
	Follow302 bool
	Host      string
	Timeout   time.Duration
	Body      io.Reader
}

type CaseResReq struct {
	Uri     string
	Method  string
	Timeout time.Duration
	Body    io.Reader
}

type HttpRespCheckInfo struct {
	StatusCode  int
	RespHeaders map[string]string
	Size        int
}

func NewCaseResRes() *CaseResReq {
	treq := &CaseResReq{}

	treq.Method = "GET"
	return treq
}

func NewHttpReq() *HttpReq {
	treq := &HttpReq{}

	treq.Method = "GET"
	treq.Headers = make(map[string]string)
	treq.Follow302 = false

	return treq
}

func NewHttpChecker() *HttpRespCheckInfo {
	check := &HttpRespCheckInfo{}

	check.StatusCode = 200
	check.RespHeaders = make(map[string]string)

	return check
}

func (tcheck *HttpRespCheckInfo) AssertRespHeaderCheck(t *testing.T, resp *http.Response) {
	if tcheck.StatusCode != 0 && resp.StatusCode != tcheck.StatusCode {
		TFatal(t, fmt.Sprintf("expect status code %d, got %d", tcheck.StatusCode, resp.StatusCode))
	}

	for k, v := range tcheck.RespHeaders {
		headers, ok := resp.Header[k]
		if !ok {
			TFatal(t, fmt.Sprintf("response header: %s not found!", k))
		}

		if len(headers) > 1 {
			TFatal(t, fmt.Sprintf("get multi headers, %s:%s", k, headers))
		}

		//t.Logf("header: %s, expect: %s, got: %s", k, v, headers[0])
		if headers[0] != v {
			TFatal(t, fmt.Sprintf("response header: %s not equal %s", headers[0], v))
		}
	}

	if tcheck.Size > 0 {
		AssertReadSize(t, tcheck.Size, resp.Body)
	}
}

func (tcheck *HttpRespCheckInfo) AssertRespBodyCheck(t *testing.T, resp *http.Response) {
	return
}

func AssertNormalHdl(t *testing.T, url string, size int) string {

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("url: %s, err: %s", url, err)
	}

	defer resp.Body.Close()

	streamReader := bufio.NewReader(resp.Body)
	streamSize, err := streamReader.Discard(size)

	if err != nil {
		return fmt.Sprintf("url: %s, err: %s", url, err)
	}

	if streamSize != size {
		return fmt.Sprintf("url: %s, read size: %d not satisfy %d",
			url, streamSize, size)
	}

	return "ok"
}

func AssertOkRequest(t *testing.T, url string, size int) {
	ret := AssertNormalHdl(t, url, size)
	if ret != "ok" {
		TFatal(t, fmt.Sprintf("url: %s assert ok failed: %s", url, ret))
	}
}

func AssertFailRequest(t *testing.T, url string, size int) {
	ret := AssertNormalHdl(t, url, size)
	if ret == "ok" {
		TFatal(t, fmt.Sprintf("url: %s assert failed, got ok"))
	}
}

func AssertMultiNormalHdl(t *testing.T, urls []string, size int, ok bool) {
	ch := make(chan string)
	for i := 0; i < len(urls); i++ {

		url := urls[i]
		//t.Logf("test url: %s", url)

		go func() {
			ch <- AssertNormalHdl(t, url, size)
		}()
	}

	var rets []string
	for i := 0; i < len(urls); i++ {
		rets = append(rets, <-ch)
	}

	for i := 0; i < len(rets); i++ {
		if rets[i] != "ok" && ok {
			TFatal(t, rets[i])
		}
		if rets[i] == "ok" && !ok {
			TFatal(t, "test failed")
		}
	}
}

func AssertReadSize(t *testing.T, size int, body io.ReadCloser) {
	streamReader := bufio.NewReader(body)
	streamSize, err := streamReader.Discard(size)

	if err != nil {
		TFatal(t, fmt.Sprintf("assert read special size error: %s", err))
	}

	if streamSize != size {
		TFatal(t, fmt.Sprintf("assert read special size error, expect size: %d got %d", size, streamSize))
	}

	return
}

func AssertTestCaseResultFromApiServer(t *testing.T, treq *CaseResReq) *http.Response {
	client := &http.Client{}
	client.Timeout = treq.Timeout

	url := "http://" + ApiServer + treq.Uri

	req, err := http.NewRequest(treq.Method, url, treq.Body)
	//	fmt.Println(url)
	if err != nil {
		TFatal(t, fmt.Sprintf("%s", err))
	}

	resp, err := client.Do(req)
	if err != nil {
		TFatal(t, fmt.Sprintf("%s", err))
	}
	return resp
}

func AssertWebSocketRequest(t *testing.T) {

}

func AssertHttpRequest(t *testing.T, treq *HttpReq, tcheck RespChecker) {
	client := &http.Client{}
	client.Timeout = treq.Timeout

	if !treq.Follow302 {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}
	//	fmt.Println(treq.Url)
	req, err := http.NewRequest(treq.Method, treq.Url, treq.Body)
	if err != nil {
		TFatal(t, fmt.Sprintf("%s", err))
	}

	if treq.Host != "" {
		req.Host = treq.Host
	}

	for k, v := range treq.Headers {
		req.Header.Add(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		TFatal(t, fmt.Sprintf("%s", err))
	}
	defer resp.Body.Close()

	tcheck.AssertRespHeaderCheck(t, resp)
	tcheck.AssertRespBodyCheck(t, resp)

	return
}

func GetHttpBody(url string) string {
	client := http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	return string(body)
}
