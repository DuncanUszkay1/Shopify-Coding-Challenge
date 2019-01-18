// template from tutorial: https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo
package main

import (
  "encoding/json"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/mongodb/mongo-go-driver/mongo"
)

type Product struct {
  Uuid            string
  Title           string
  Price           string
  Inventorycount  string
}

// our main function
func main() {
  setup_db()

  router := mux.NewRouter()
  router.HandleFunc("/products", GetProducts).Methods("GET")
  router.HandleFunc("/products/{uuid}", GetProduct).Methods("GET")
  log.Fatal(http.ListenAndServe(":8000", router))
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
  products, err := get_records("products")
  if err == mongo.ErrNilDocument {
    w.WriteHeader(404)
  } else if err != nil {
    w.WriteHeader(500)
  } else {
    json.NewEncoder(w).Encode(products)
  }
}
func GetProduct(w http.ResponseWriter, r *http.Request) {
  product, err := get_record("products", mux.Vars(r)["uuid"])
  if err == mongo.ErrNoDocuments {
    w.WriteHeader(404)
  } else if err != nil {
    w.WriteHeader(500)
  } else {
    json.NewEncoder(w).Encode(product)
  }
}
