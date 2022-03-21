package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Header struct {
	Key   string
	Value string
}

type RequestOption struct {
	Timeout time.Duration
}
type Request struct {
	Option  *RequestOption
	Url     string
	Headers []*Header
	Body    interface{}
}

type ModOption func(option *RequestOption)

func NewRequest(option *RequestOption) *Request {
	if option == nil {
		option = &RequestOption{}
	}
	return &Request{
		Option: option,
	}
}

// 检测HTTP请求的结果是否失败
func CheckResponseIsFail(r *http.Response) bool {
	if http.StatusOK <= r.StatusCode && r.StatusCode < http.StatusMultipleChoices {
		return false
	}
	return true
}

func (r *Request) Do(method string) ([]byte, error) {
	var client = &http.Client{}
	if r.Option.Timeout > 0 {
		client.Timeout = r.Option.Timeout
	} else {
		client.Timeout = time.Second * 30
	}
	var (
		js  []byte = nil
		err error
	)

	if r.Body != nil {
		js, err = json.Marshal(r.Body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, r.Url, bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if r.Headers != nil {
		for _, h := range r.Headers {
			req.Header.Add(h.Key, h.Value)
		}
	}

	rsps, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer rsps.Body.Close()
	body, err := ioutil.ReadAll(rsps.Body)
	if ok := CheckResponseIsFail(rsps); ok {
		return nil, fmt.Errorf("%s", body)
	}
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (r *Request) SetHeader(key, val string) {
	r.Headers = append(r.Headers, &Header{Key: key, Value: val})
}

func (r *Request) Get(rawURL string) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	u.RawQuery = query.Encode()
	r.Url = u.String()
	body, err := r.Do("GET")
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (r *Request) Post(rawURL string, data interface{}) ([]byte, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	u.RawQuery = query.Encode()

	r.Url = u.String()
	r.Body = data
	body, err := r.Do("POST")
	if err != nil {
		return nil, err
	}
	return body, nil
}
