package db

import (
	"log"
	"os"
	"time"

	"ChiliOverFlow/pkg/config"
	"github.com/jmoiron/sqlx"

	// Registers the sqlite3 db driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Bundle controls all the data flow into and out of the db layer
type Bundle struct {
	db   *sqlx.DB
	conf *viper.Viper
}

// New creates a new fully initialized Bundle
func New() (Bundle, error) {
	model := Bundle{}
	conf, err := config.New()
	if err != nil {
		return model, err
	}
	model.conf = conf

	configPath, _ := homedir.Expand((conf.GetString("configPath")))
	file := configPath + "/chilioverflow.sqlite"

	if _, err := os.Stat(file); os.IsNotExist(err) {
		os.Create(file)
	}

	db, err := sqlx.Open("sqlite3", file)
	if err != nil {
		return model, err
	}

	model.db = db
	model.CreateTablesIfNeeded()
	chilis, err := model.GetAllStored()

	if len(chilis) == 0 {
		csvChilis := ReadCSV()
		for _, chili := range csvChilis {
			log.Print("Loading ", chili.Name, "\n")
			_, err := model.CreateChili(chili)
			if err != nil {
				log.Print("Error: ", err)
			}

		}
	}
	return model, nil
}

// CreateTablesIfNeeded ensures that the db has the necessary tables
func (m *Bundle) CreateTablesIfNeeded() {
	m.db.Exec(`
create table if not exists Chilis (
id integer primary key,
name varchar(255),
family varchar(255),
origin varchar(255),
heat integer,
taste text,
info text,
color text,
size text,
fruit text,
plant_form text,
plant_size text,
guide text,
image varchar(255),
seeds_per_request integer,
date integer)
`)
	m.db.Exec(`
create table if not exists Input (
id integer primary key,
chili integer,
producer text,
quantity integer,
price real,
date integer
use_before integer)
`)
	m.db.Exec(`
create table if not exists Output (
id integer primary key,
chili integer ,
reciever text,
quantity integer,
date integer)
`)
}

// DateRange is an inclusive range of dates
type DateRange struct {
	Start time.Time
	End   time.Time
}

// Chili stores information about available seeds
type Chili struct {
	Id                int
	Name              string
	Origin            string
	Family            string
	Heat              string
	Taste             string
	Info              string
	Color             string
	Size              string
	Fruit             string
	Plant_Form        string
	Plant_Size        string
	Guide             string
	Image             string
	Seeds_Per_Request int
	Date              string
}

// ChiliEntry defines quantities of seeds for transactions
type ChiliEntry struct {
	Barcode  string
	Quantity int
	Date     time.Time
}

// StockedChili is an extension of Chili with an additional field for quantity
type StockedChili struct {
	Chili
	Quantity int
}
