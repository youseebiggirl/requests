package requests

import (
	"fmt"
	"testing"
)

var cookie = ""

func TestGet(t *testing.T) {
	r := Get("https://weibo.com/ajax/favorites/all_fav?page=1", WithCookie(cookie))
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
