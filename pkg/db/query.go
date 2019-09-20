package db

import (
	"database/sql"
)

// BarcodeExists checks if a barcode is already in the database
func (m *Bundle) BarcodeExists(bc string) (bool, error) {
	var barcode string
	err := m.db.Get(&barcode, "select barcode from Drinks where barcode = ? limit 1", bc)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if barcode == bc {
		return true, nil
	}
	return false, nil
}

// GetAllStored returns every saved Chili row in the database
func (m *Bundle) GetAllStored() ([]Chili, error) {
	var drinks []Chili
	err := m.db.Select(&drinks, "select * from Chilis")
	return drinks, err
}

// GetCountByBarcode returns the total number of currently stocked beers with a specific barcode
func (m *Bundle) GetCountByName(name string) (int, error) {
	var input, output int
	if err := m.db.Get(&input, "select case when sum(quantity)is null then 0 else sum(quantity) end quantity from Input where barcode = ?", name); err != nil {
		return -1, err
	}
	if err := m.db.Get(&output, "select case when sum(quantity) is null then 0 else sum(quantity) end quantity from Output where barcode = ?", name); err != nil {
		return -1, err
	}

	return input - output, nil
}

// GetDrinkByBarcode returns all stored information about a drink based on its barcode
func (m *Bundle) GetDrinkByBarcode(bc string) (Chili, error) {
	var d Chili
	err := m.db.Get(&d, "select * from Drinks where barcode = ?", bc)
	return d, err
}

// GetInventoryTotalVariety returns the total number of beer varieties in stock
func (m *Bundle) GetInventoryTotalVariety() ([]StockedChili, error) {
	var result []StockedChili

	sql := `
select Chilis.*,
  case
    when Input.InputQuantity is null then 0
    when Output.OutputQuantity is null then Input.InputQuantity
    else (Input.InputQuantity - Output.OutputQuantity)
  end as quantity
from Chilis

left join (
  select chili, sum(quantity) as InputQuantity
  from Input
  group by chili
) as Input
on Chilis.ID = Input.Chili

left join (
  select chili, sum(quantity) as OutputQuantity
  from Output
  group by chili
) as Output
on Chilis.ID = Output.chili

where quantity > 0 order by Chilis.Name`

	err := m.db.Select(&result, sql)
	return result, err
}

// GetInventory returns every drink with at least one quantity in stock, sorted by Type
func (m *Bundle) GetInventory() ([]Chili, error) {
	return m.GetAllStored()
}
// GetInputWithinDateRange returns every drink inputted within a date range, inclusive
func (m *Bundle) GetInputWithinDateRange(dates DateRange) (result []StockedChili, err error) {
	sql := `
select A.*,
  case
    when C.InputQuantity is null then 0
    else C.InputQuantity
  end as quantity
from Drinks as A

left join (
  select barcode, sum(quantity) as InputQuantity
  from Input as O where O.Date >= ? and O.Date <= ?
  group by barcode
) as C
on A.Barcode = C.Barcode

where quantity > 0
order by A.Brand
`
	err = m.db.Select(&result, sql, dates.Start, dates.End)
	return result, err
}

// GetOutputWithinDateRange returns every drink served within a date range, inclusive
func (m *Bundle) GetOutputWithinDateRange(dates DateRange) (result []StockedChili, err error) {
	sql := `
select A.*,
  case
    when C.OutputQuantity is null then 0
    else C.OutputQuantity
  end as quantity
from Drinks as A

left join (
  select barcode, sum(quantity) as OutputQuantity
  from Output as O where O.Date >= ? and O.Date <= ?
  group by barcode
) as C
on A.Barcode = C.Barcode

where quantity > 0
order by A.Brand
`
	err = m.db.Select(&result, sql, dates.Start, dates.End)
	return result, err
}
