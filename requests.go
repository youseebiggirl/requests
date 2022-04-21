package requests

import (
	"encoding/json"
	"io"
	"net/http"
)

type Options func(*requests)

func WithHeaders(header http.Header) Options {
	return func(r *requests) {
		r.header = header
	}
}

func WithCookie(cookie string) Options {
	return func(r *requests) {
		r.cookie = cookie
	}
}

type requests struct {
	url    string
	header http.Header
	cookie string
}

var defaultRequests = &requests{}

func (r *requests) init(url string, op []Options) {
	if r == nil {
		return
	}
	if r.header == nil {
		r.header = make(http.Header)
	}
	r.url = url
	for _, option := range op {
		option(r)
	}
}

func Get(url string, op ...Options) *result {
	defaultRequests.init(url, op)
	return get(defaultRequests)
}

func get(r *requests) *result {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		panic(err)
	}
	if len(defaultRequests.header) > 0 {
		req.Header = defaultRequests.header.Clone()
	}
	if defaultRequests.cookie != "" {
		req.Header.Set("Cookie", defaultRequests.cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	return &result{resp: *resp}
}

type result struct {
	resp http.Response
}

func (r *result) Text() string {
	b, err := readResponseBody[string](r.resp)
	if err != nil {
		panic(err)
	}
	return b
}

func (r *result) Unmarshal(obj any) {
	b, err := readResponseBody[[]byte](r.resp)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &obj); err != nil {
		panic(err)
	}
}

func (r *result) Map() (m map[string]any) {
	b, err := readResponseBody[[]byte](r.resp)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(b, &m); err != nil {
		panic(err)
	}
	return
}

func (r *result) StatusCode() int {
	return r.resp.StatusCode
}

func (r *result) StatusText() string {
	return http.StatusText(r.resp.StatusCode)
}

// readResponseBody 从 resp 中读取 body，并根据指定的泛型参数，转换成对应的类型，比如:
// readResponseBody[string](resp)，会返回一个 string 类型的 body
func readResponseBody[T []byte|string|json.RawMessage](resp http.Response) (t T, err error) {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	t = T(b)
	return
}
