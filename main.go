package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Product struct {
	Id int
	Nama string
	Harga int
	Stok int
}

var Products = []Product{
	{1,"Baju", 2008, 12},
	{2,"Kemeja", 2009, 20},
	{3,"Jeans", 2010, 10},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)  {
		fmt.Fprintln(w, "Hello World!")
	})
	http.HandleFunc("/product", func (w http.ResponseWriter, r *http.Request)  {
		idParam := r.URL.Path[len("/product/"):]//menentukan parameter id
		id, err := strconv.Atoi(idParam)
		if err != nil {
			http.Error(w, "Invalid Product Id", http.StatusBadRequest)
		}

		var foundIndex = -1
		for i, p := range Products {
			if p.Id == id {
				foundIndex = i
				break
			}
		}

		if foundIndex == -1 {
			http.Error(w, "Id not found", http.StatusNotFound)
		}

		if r.Method == "GET" {
			res, err := json.Marshal(Products)
		if err != nil {
			http.Error(w, "gagal koneksi ke Json", http.StatusInternalServerError)
		}
		w.Write(res)
		w.Header().Set("Content-Type","applicaiton/json")
		return

		} else if r.Method == "PUT"{
			var updateProduct Product
			err := json.NewDecoder(r.Body).Decode(&updateProduct)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if updateProduct.Id <= -1 || updateProduct.Nama == "" || updateProduct.Harga == 0 || updateProduct.Stok == 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprint(w, "Invalid Product data")
				return
			}
			Products[foundIndex] = updateProduct; // menggunakan found index untuk mengupdate
			msg := map[string]string{
				"Message": "Product Updated",
			}
			res, err := json.Marshal(msg)
			if err != nil {
				http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
				return
			}
			w.Write(res)

		} else if r.Method == "DELETE" {
			_ = append(Products[:foundIndex], Products[foundIndex+1:]...)
			msg := map[string]string{
				"Message": "Product Deleted",
			}
			res, err := json.Marshal(msg)
			if err != nil {
				http.Error(w, "gagal Konvesi Json", http.StatusInternalServerError)
				return
			}
			w.Write(res)
			
		} else {
			http.Error(w, "Merhod tidak diizinkan", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/products", func (w http.ResponseWriter, r *http.Request)  {
		if r.Method == "GET" {
			res, err := json.Marshal(Products)
			if err != nil {
				http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
				return
			}
			w.Write(res)
			w.Header().Set("Content-Type", "application/json")
			return
		} else if r.Method == "POST" {
			var product Product
			err := json.NewDecoder(r.Body).Decode(&product)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if product.Id <= 0 || product.Nama == "" || product.Harga <= 0 || product.Stok <= 0 {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "Invalid product data")
				return
			}
			Products = append(Products, product)
			w.WriteHeader(http.StatusCreated)
			msg := map[string]string{
				"Message": "Product Created",
			}
			res, err := json.Marshal(msg)
			if err != nil {
				http.Error(w, "Gagal Konversi Json", http.StatusInternalServerError)
				return
			}
			w.Write(res)
		} else {
			http.Error(w, "method tidak diizinkan", http.StatusMethodNotAllowed)
		}
	})
	fmt.Print("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}