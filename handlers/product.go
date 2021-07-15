package handlers

import (
	"go-service/data"
	"log"
	"net/http"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProduct(rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Product) getProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Get list Product")
	lp := data.GetProduct()
	err := lp.ToJSON(rw)
	if err != nil {
		p.l.Println(rw, "Failed to marchall object", http.StatusInternalServerError)
		return
	}
}

func (p *Product) addProduct(rw http.ResponseWriter, r *http.Request) {

}
