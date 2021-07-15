package handlers

import (
	"go-service/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
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

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		reg := regexp.MustCompile(`/([0-9]+)`)
		g := reg.FindAllStringSubmatch(r.URL.Path, -1)

		if len(g) != 1 {
			http.Error(rw, "[ERROR] Invalid URL", http.StatusBadRequest)
			return
		}

		if len(g[0]) != 2 {
			http.Error(rw, "[ERROR] Invalid URL", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "[ERROR] Failed to convert id", http.StatusInternalServerError)
			return
		}

		p.updateProduct(id, rw, r)
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
	p.l.Println("Post a new product")
	prod := data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		p.l.Println("[ERROR] Decode json failed")
		http.Error(rw, "[ERROR] Decode json failed", http.StatusBadRequest)
		return
	}
	data.AddProduct(&prod)
}

func (p *Product) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Update a Product")
	prod := data.Product{}
	err := prod.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "[ERROR] Faild to read request", http.StatusBadRequest)
		return
	}
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrorNotFound {
		http.Error(rw, "Product not found", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(rw, "[ERROR] Failed to update product", http.StatusInternalServerError)
		return
	}
}
