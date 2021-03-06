package requests

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

var cookie = ""

func TestGET(t *testing.T) {
	r := GET("https://weibo.com/ajax/favorites/all_fav?page=1", WithCookie(cookie))
	fmt.Printf("%v %v\n", r.StatusCode(), r.StatusText())
	text := r.Text()
	fmt.Println(text)

	m := r.Map()
	data, ok := m["data"].([]any)
	if ok {
		for _, v := range data {
			id := v.(map[string]any)["idstr"]
			fmt.Printf("%T %v\n", id, id)
			//s := strconv.FormatFloat(id.(float64), 'f', 0, 64)
			//fmt.Println(s)
			break
		}
	}
}

func TestPOST(t *testing.T) {
	j := `{"id":"4760653908937563"}`
	// b, err := json.Marshal(j)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	b := []byte(j)
	//data := bytes.NewBuffer(b)

	r := POST("https://weibo.com/ajax/statuses/createFavorites?",
		WithCookie(cookie),
		WithHeaders(http.Header{
			"x-xsrf-token": []string{""},
			//"content-type": []string{"application/json"},
		}),
		WithJson(b),
	)
	fmt.Printf("%v %v\n", r.StatusCode(), r.StatusText())
	text := r.Text()
	fmt.Println(text)
}

func TestPOST1(t *testing.T) {
	j := `{"id":"4760653908937563"}`
	// b, err := json.Marshal(j)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	data := bytes.NewBuffer([]byte(j))
	r, err := http.NewRequest("POST",
		"https://weibo.com/ajax/statuses/createFavorites?",
		data)
	if err != nil {
		t.Fatal(err)
	}
	r.Header.Set("Cookie", cookie)
	r.Header.Set("x-xsrf-token", "vXiIhmjvwyPCygYO6k9Lusdf")
	r.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("resp.Status: %v\n", resp.Status)
}
