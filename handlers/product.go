package handlers

import (
	"context"
	"fmt"
	"go-service/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Product struct {
	l *log.Logger
}

func NewProduct(l *log.Logger) *Product {
	return &Product{l}
}

func (p *Product) GetProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Get list Product")
	lp := data.GetProduct()
	err := lp.ToJSON(rw)
	if err != nil {
		p.l.Println(rw, "Failed to marchall object", http.StatusInternalServerError)
		return
	}
}

func (p *Product) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Post a new product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Product) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Update a Product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Failed to convert id", http.StatusBadRequest)
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

type KeyProduct struct{}

func (p *Product) MiddlewareValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Error to read request", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			http.Error(rw, fmt.Sprintf("Invalid input: %s\n", err), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
