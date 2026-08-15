package abstract

import "net/http"

type Imeter interface {
	Init(mux *http.ServeMux) error
	Increment() error
	Get() (uint, error)
}
