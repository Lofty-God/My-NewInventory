package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) Initialise(DbUser string, DbPassword string, DBName string) {
	var err error
	connectionString := fmt.Sprintf("%v:%v@(127.0.0.1:3306)/%v", DbUser, DbPassword, DBName)
	app.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalln(nil)
	}
	app.Router = mux.NewRouter().StrictSlash(true)
	app.handleRoutes()
	fmt.Println(nil)

}
func (app *App) sendError(w http.ResponseWriter, StatusCode int, err string) {
	error_message := map[string]string{"error": err}
	app.sendResponse(w, StatusCode, error_message)

}
func (app *App) sendResponse(w http.ResponseWriter, StatusCode int, Payload interface{}) {
	response, _ := json.Marshal(Payload)
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(StatusCode)
	w.Write(response)

}
func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))

}
func (app *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.sendError(w, http.StatusBadRequest, "invalid payload requested")
		return
	}
	err = p.createProduct(app.DB)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.sendResponse(w, http.StatusCreated, p)

}

func (app *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.sendError(w, http.StatusBadRequest, "invalid payload requested")
		return
	}

	p.Id = key
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = p.updateProduct(app.DB)
	if err != nil {
		log.Fatal(err)
	}
	app.sendResponse(w, http.StatusAccepted, p)

}
func (app *App) deleteProducts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.sendError(w, http.StatusBadRequest, "invalid id requested")
		return
	}

	p := product{Id: key}
	err = p.deleteProducts(app.DB)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.sendResponse(w, http.StatusAccepted, p)

}

func (app *App) getProducts(w http.ResponseWriter, r *http.Request) {
	products, err := getProducts(app.DB)
	if err != nil {
		app.sendError(w, http.StatusInternalServerError, err.Error())
		return
	}
	app.sendResponse(w, http.StatusAccepted, products)
}
func (app *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, err := strconv.Atoi(vars["id"])
	if err != nil {
		app.sendError(w, http.StatusBadRequest, "invalid product id")
		return
	}
	p := product{Id: key}
	err = p.getProduct(app.DB)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			app.sendError(w, http.StatusNotFound, "product not found")
		default:
			app.sendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	app.sendResponse(w, http.StatusOK, p)
}

func (app *App) handleRoutes() {
	app.Router.HandleFunc("/products", app.getProducts).Methods("Get")
	app.Router.HandleFunc("/product/{id}", app.getProduct).Methods("Get")
	app.Router.HandleFunc("/product", app.createProduct).Methods("Post")
	app.Router.HandleFunc("/product/{id}", app.updateProduct).Methods("Put")
	app.Router.HandleFunc("/products/{id}", app.deleteProducts).Methods("Delete")

}
