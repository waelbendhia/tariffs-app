package types

import (
	"database/sql"
	"fmt"
	"time"

	panicers "github.com/waelbendhia/tariffs-app/panicers"
)

// Tariff is a unitary price for playtime
type Tariff struct {
	ID           int64
	PricePerUnit int64
	UnitSize     time.Duration
	CreatedAt    time.Time
}

func createTariffsTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS tariffs (
			price_per_unit INTEGER NOT NULL,
			unit_size INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
	)
	panicers.WrapAndPanicIfErr(err, "Could not create tariffs table")
}

// GetAllTariffs from db, panics on failure
func GetAllTariffs(db *sql.DB) []Tariff {
	rows, err := db.Query("SELECT rowid, * FROM tariffs;")
	panicers.WrapAndPanicIfErr(err, "Could not query tariffs table")

	defer func() {
		panicers.WrapAndPanicIfErr(rows.Close(), "Error while closing rows")
	}()

	var ts []Tariff
	for rows.Next() {
		ts = append(ts, *scanTariff(rows))
	}
	panicers.WrapAndPanicIfErr(rows.Err(), "Error while querying for machines")
	return ts
}

// GetLatestTariff from db, if none are found returns nil, panics on failure
func GetLatestTariff(db *sql.DB) *Tariff {
	return scanTariff(db.QueryRow("SELECT rowid, * FROM tariffs ORDER BY created_at DESC;"))
}

// GetTariffByID from db, panics on failure
func GetTariffByID(id int64, db *sql.DB) *Tariff {
	return scanTariff(db.QueryRow("SELECT rowid, * FROM tariffs WHERE rowid = ?;", id))
}

// Insert t in db, panics on failure
func (t *Tariff) Insert(db *sql.DB) Tariff {
	res, err := db.Exec(
		`INSERT INTO 
			tariffs (price_per_unit, unit_size)
			VALUES (?, ?);`,
		t.PricePerUnit,
		t.UnitSize,
	)
	panicers.WrapAndPanicIfErr(err, "Could not to insert tariff")

	id, err := res.LastInsertId()
	panicers.WrapAndPanicIfErr(err, "Could not retrieve ID after insert")

	return Tariff{id, t.PricePerUnit, t.UnitSize, t.CreatedAt}
}

func (t *Tariff) Equals(o interface{}) bool {
	if t == nil || o == nil {
		return t == nil && o == nil
	}

	switch v := o.(type) {
	case *Tariff:
		if v == nil {
			return false
		}
		return t.UnitSize == v.UnitSize && t.PricePerUnit == v.PricePerUnit
	case Tariff:
		return t.UnitSize == v.UnitSize && t.PricePerUnit == v.PricePerUnit
	default:
		return false
	}
}

func scanTariff(row scanner) *Tariff {
	t := &Tariff{}
	err := row.Scan(&t.ID, &t.PricePerUnit, &t.UnitSize, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil
	}
	panicers.WrapAndPanicIfErr(err, "Could not scan tariff from row")

	return t
}

func (t *Tariff) String() string {
	var (
		unitInTimeStr int
		timeStr       string
	)
	switch {
	case t.UnitSize%time.Hour == 0:
		unitInTimeStr = int(t.UnitSize / time.Hour)
		if unitInTimeStr > 1 {
			timeStr = "heure"
		} else {
			timeStr = "heures"
		}
	case t.UnitSize%time.Minute == 0:
		unitInTimeStr = int(t.UnitSize / time.Minute)
		if unitInTimeStr > 1 {
			timeStr = "minutes"
		} else {
			timeStr = "minutes"
		}
	default:
		unitInTimeStr = int(t.UnitSize / time.Second)
		if unitInTimeStr > 1 {
			timeStr = "minutes"
		} else {
			timeStr = "minutes"
		}
	}
	return fmt.Sprintf(
		"%d Millimes par %d %s",
		t.PricePerUnit,
		unitInTimeStr,
		timeStr,
	)
}
