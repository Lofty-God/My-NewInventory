package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	var err error
	a.Initialise(DbUser, DbPassword, "test")
	if err != nil {
		log.Fatal("error occured while initialising the database")
	}
	createTable()
	m.Run()

}

func createTable() {
	createTableQuery := `create table IF NOT EXISTS products(
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255),
		quantity int,
		price float(5,2),
		PRIMARY KEY(id)
	);`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("insert into products(name, quantity, price)values('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}
func (a *App) clearTable() {
	a.DB.Exec("delete from products")
	a.DB.Exec("alter table products AUTO_INCREMENT=1")
	log.Println("clearTable")

}

func TestGetProduct(t *testing.T) {
	a.clearTable()
	addProduct("chair", 220, 200)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}
func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatusCode int) {
	if expectedStatusCode != actualStatusCode {
		t.Errorf("expected status: %v, actual:%v", expectedStatusCode, actualStatusCode)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	a.clearTable()
	var product = []byte(`{"name":"chair", "quantity":250, "price":350}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("content-type", "application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}

	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "chair" {
		t.Errorf("expected value:%v, got: %v ", "chair", m["name"])
	}
	log.Printf("%T", m["quantity"])
	if m["quantity"] != 250.00 {
		t.Errorf("expected:%v, Got: %v", 250.00, m["quantity"])
	}

}

func TestDeleteProduct(t *testing.T) {
	a.clearTable()
	addProduct("divider", 50, 125)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/products/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusAccepted, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}

func TestUpdateProduct(t *testing.T) {
	a.clearTable()
	addProduct("Book", 2, 45)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"name":"Book", "quantity":5, "price":60}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("content-type", "application/json")
	response = sendRequest(req)

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
		t.Errorf("expected id: %v, Got: %v", oldValue["id"], newValue["id"])
	}

	if oldValue["name"] != newValue["name"] {
		t.Errorf("expected name: %v, Got: %v", oldValue["name"], newValue["name"])
	}

	if oldValue["quantity"] == newValue["quantity"] {
		t.Errorf("expected quuantity: %v, Got: %v", newValue["quantity"], oldValue["quantity"])
	}

	if oldValue["price"] == newValue["price"] {
		t.Errorf("expected price: %v, Got: %v", newValue["price"], oldValue["price"])
	}

}
