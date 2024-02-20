package hook

import (
	"net/http"
	"sync"
)

type Hook struct {
	http.Handler
	lock sync.Mutex
}

type WriteReset interface {
	Write(p []byte) (n int, err error)
	Reset()
}

func NewHook(ob WriteReset) *Hook {
	r := http.NewServeMux()

	t := &Hook{}

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = ""
		t.lock.Lock()
		ob.Reset()
		r.Write(ob)
		defer t.lock.Unlock()
	})

	t.Handler = r
	return t
}
