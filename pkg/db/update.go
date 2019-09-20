package db

import (
	"database/sql"
	"time"
)

// ClearInputTable deletes all stocking records
func (m *Bundle) ClearInputTable() error {
	_, err := m.db.Exec("delete from Input")
	return err
}

// ClearOutputTable deletes all serving records
func (m *Bundle) ClearOutputTable() error {
	_, err := m.db.Exec("delete from Output")
	return err
}

// CreateChili adds an entry to the Drinks table, returning the id
func (m *Bundle) CreateChili(chili Chili) (int, error) {

	res, err := m.db.Exec(`
		insert into Chilis (
		Name            ,
		Origin          ,
		Family          ,
		Heat            ,
		Taste           ,
		Info             ,
		Color           ,
		Size            ,
		Fruit           ,
		Plant_Form       ,
		Plant_Size       ,
		Guide           ,
		Image           ,
		Seeds_Per_Request ,
		Date       ) Values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		chili.Name,
		chili.Origin,
		chili.Family,
		chili.Heat,
		chili.Taste,
		chili.Info,
		chili.Color,
		chili.Size,
		chili.Fruit,
		chili.Plant_Form,
		chili.Plant_Size,
		chili.Guide,
		chili.Image,
		chili.Seeds_Per_Request,
		time.Now())
	if err != nil {
		return -1, err
	}
	return getID(res)
}

// DeleteDrink removes an entry from the Drinks table using its barcode
func (m *Bundle) DeleteDrink(bc string) error {
	_, err := m.db.Exec("delete from Drinks where barcode = ?", bc)
	return err
}

// InputDrinks adds an entry to the Input table, returning the id
func (m *Bundle) InputDrinks(d ChiliEntry) (int, error) {
	now := time.Now().Unix()
	res, err := m.db.Exec(
		"insert into Input (barcode, quantity, date) Values (?, ?, ?)", d.Barcode, d.Quantity, now)
	if err != nil {
		return -1, err
	}
	return getID(res)
}

// UndoInputDrinks removes an entry from the Input table by id
func (m *Bundle) UndoInputDrinks(id int) error {
	_, err := m.db.Exec("delete from Input where id = ?", id)
	return err
}

// OutputDrinks adds an entry to the Output table, returning the id
func (m *Bundle) OutputDrinks(d ChiliEntry) (int, error) {
	now := time.Now().Unix()
	res, err := m.db.Exec(
		"insert into Output (barcode, quantity, date) Values (?, ?, ?)", d.Barcode, d.Quantity, now)
	if err != nil {
		return -1, err
	}
	return getID(res)
}

// UndoOutputDrinks removes an entry from the Output table by id
func (m *Bundle) UndoOutputDrinks(id int) error {
	_, err := m.db.Exec("delete from Output where id = ?", id)
	return err
}

func getID(result sql.Result) (int, error) {
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(id), nil
}
