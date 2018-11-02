package model

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	// Registers the sqlite3 db driver
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

// Model controls all the data flow into and out of the db layer
type Model struct {
	db *sqlx.DB
}

// New creates a new fully initialized Model
func New() (Model, error) {
	model := Model{}

	file := "abv.sqlite"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.Create(file)
	}

	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		return model, err
	}

	model.db = db
	model.CreateTablesIfNeeded()
	return model, nil
}

// CreateTablesIfNeeded ensures that the db has the necessary tables
func (m Model) CreateTablesIfNeeded() {
	m.db.Exec("create table if not exists Drinks (barcode integer primary key, brand varchar(255), name varchar(255), abv real, ibu real, type varchar(255), date integer)")
	m.db.Exec("create table if not exists Input (id integer primary key, barcode integer, quantity integer, date integer)")
	m.db.Exec("create table if not exists Output (id integer primary key, barcode integer, quantity integer, date integer)")
}

// Drink stores information about an available beverage
type Drink struct {
	Barcode int
	Brand   string
	Name    string
	Abv     float64
	Ibu     int
	Type    string
	Date    int
}

// DrinkEntry defines quantities of drinks for transactions
type DrinkEntry struct {
	Barcode  int
	Quantity int
	Date     int
}

//TODO DeleteDrink
//TODO UpdateDrink

// BarcodeExists checks if a barcode is already in the database
func (m Model) BarcodeExists(bc int) (bool, error) {
	var barcodes []int
	if err := m.db.Select(&barcodes, "select barcode from Drinks"); err != nil {
		return false, err
	}
	for _, code := range barcodes {
		if code == bc {
			return true, nil
		}
	}
	return false, nil
}

// CreateDrink adds an entry to the Drinks table, returning the id
func (m Model) CreateDrink(d Drink) (int, error) {
	now := time.Now().Unix()
	res, err := m.db.Exec(
		"insert into Drinks (barcode, brand, name, abv, ibu, type, date) Values (?, ?, ?, ?, ?, ?, ?)", d.Barcode, d.Brand, d.Name, d.Abv, d.Ibu, d.Type, now)
	if err != nil {
		return -1, err
	}
	return getID(res)
}

// GetAllStoredDrinks returns every saved Drink row in the database
func (m Model) GetAllStoredDrinks() ([]Drink, error) {
	var drinks []Drink
	err := m.db.Select(&drinks, "select * from Drinks")
	return drinks, err
}

// InputDrinks adds an entry to the Input table, returning the id
func (m Model) InputDrinks(d DrinkEntry) (int, error) {
	now := time.Now().Unix()
	res, err := m.db.Exec(
		"insert into Input (barcode, quantity, date) Values (?, ?, ?)", d.Barcode, d.Quantity, now)
	if err != nil {
		return -1, err
	}
	return getID(res)
}

// OutputDrinks adds an entry to the Output table, returning the id
func (m Model) OutputDrinks(d DrinkEntry) (int, error) {
	now := time.Now().Unix()
	res, err := m.db.Exec(
		"insert into Output (barcode, quantity, date) Values (?, ?, ?)", d.Barcode, d.Quantity, now)
	if err != nil {
		return -1, err
	}
	return getID(res)
}

func getID(result sql.Result) (int, error) {
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}
