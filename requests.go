package requests

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
)

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
}

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

func WithData(data io.Reader) Options {
	return func(r *requests) {
		r.data = data
	}
}

type requests struct {
	method string
	url    string
	header http.Header
	cookie string
	data   io.Reader
}

var requestsPool = sync.Pool{New: func() any { return &requests{} }}

func (r *requests) init(url, method string, op []Options) {
	if r == nil {
		r = &requests{}
	}
	if r.header == nil {
		r.header = make(http.Header)
	}
	r.url = url
	r.method = method
	for _, option := range op {
		option(r)
	}
}

// reset 清空 r，主要用于 put 到 sync.pool 来复用
func (r *requests) reset() {
	if r == nil {
		return
	}
	for k := range r.header {
		delete(r.header, k)
	}
	r.method = ""
	r.url = ""
	r.data = nil
	r.cookie = ""
}

func GET(url string, op ...Options) *result {
	r := requestsPool.Get().(*requests)
	r.init(url, "GET", op)
	return get(r)
}

func get(r *requests) *result {
	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		log.Fatalf("create GET request error: %v\n", err)
	}
	if len(r.header) > 0 {
		req.Header = r.header.Clone()
	}
	if r.cookie != "" {
		req.Header.Set("Cookie", r.cookie)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("do GET Request error: %v\n", err)
	}
	r.reset()
	requestsPool.Put(r)
	return &result{resp: *resp}
}

func POST(url string, op ...Options) *result {
	r := requestsPool.Get().(*requests)
	r.init(url, "POST", op)
	return post(r)
}

func post(r *requests) *result {
	req, err := http.NewRequest("POST", r.url, r.data)
	if err != nil {
		log.Fatalf("create POST request error: %v\n", err)
	}
	if len(r.header) > 0 {
		req.Header = r.header.Clone()
	}
	if r.cookie != "" {
		req.Header.Set("Cookie", r.cookie)
	}
	//log.Println(req.URL.Scheme, req.Method)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("do POST Request error: %v\n", err)
	}
	r.reset()
	requestsPool.Put(r)
	return &result{resp: *resp}
}

type result struct {
	resp http.Response
}

func (r *result) Text() string {
	return readResponseBody[string](r.resp)
}

func (r *result) Unmarshal(obj any) {
	b := readResponseBody[[]byte](r.resp)
	if err := json.Unmarshal(b, &obj); err != nil {
		log.Fatalf("unmarshal body to obj error: %v\n", err)
	}
}

func (r *result) Map() (m map[string]any) {
	b := readResponseBody[[]byte](r.resp)
	if err := json.Unmarshal(b, &m); err != nil {
		log.Fatalf("unmarshal body to map[string]any error: %v\n", err)
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
func readResponseBody[T []byte | string | json.RawMessage](resp http.Response) (t T) {
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("read response body error: %v\n", err)
	}
	t = T(b)
	return
}
