package types

import (
	"database/sql"
	"time"
)

type Tariff struct {
	ID           int64
	PricePerUnit int64
	UnitSize     time.Duration
	CreatedAt    time.Time
}

func CreateTariffsTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS tariffs (
			price_per_unit INTEGER NOT NULL,
			unit_size INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	)
	if err != nil {
		panic(err)
	}
}

func GetAllTariffs(db *sql.DB) []Tariff {
	rows, err := db.Query("SELECT rowid, * FROM tariffs;")
	if err != nil {
		panic(err)
	}
	var ts []Tariff
	for rows.Next() {
		ts = append(ts, *scanTariff(rows))
	}
	return ts
}

func GetTariffLatest(db *sql.DB) *Tariff {
	return scanTariff(db.QueryRow("SELECT rowid, * FROM tariffs ORDER BY created_at DESC;"))
}

func GetTariffByID(id int64, db *sql.DB) *Tariff {
	return scanTariff(db.QueryRow("SELECT rowid, * FROM tariffs WHERE rowid = ?;", id))
}

func (t *Tariff) Insert(db *sql.DB) Tariff {
	res, err := db.Exec(
		`INSERT INTO 
			tariffs (price_per_unit, unit_size)
			VALUES (?, ?);`,
		t.PricePerUnit,
		t.UnitSize,
	)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return Tariff{id, t.PricePerUnit, t.UnitSize, t.CreatedAt}
}

func scanTariff(row scanner) *Tariff {
	t := &Tariff{}
	err := row.Scan(&t.ID, &t.PricePerUnit, &t.UnitSize, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return t
}
