package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("insert into products(name, quantity, price) values('%v',%v, %v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	p.Id = int(id)
	return nil

}
func (p *product) updateProduct(db *sql.DB) error {
	query := fmt.Sprintf("update products set name='%v', quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.Id)
	result, err := db.Exec(query)
	if err != nil {
		log.Fatalln(err)
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row exist")
	}
	return err
}
func (p *product) deleteProducts(db *sql.DB) error {
	query := fmt.Sprintf("delete from products where id=%v", p.Id)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such row exists")
	}
	return err
}
func getProducts(db *sql.DB) ([]product, error) {
	query := "select id, name, quantity, price from products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil

}
func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("select name, quantity, price from products where id=%v", p.Id)
	row := db.QueryRow(query)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil

}
