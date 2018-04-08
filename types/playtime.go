package types

import (
	"database/sql"
	"math"
	"time"

	"github.com/pkg/errors"
	"github.com/waelbendhia/tariffs-app/database"
	"github.com/waelbendhia/tariffs-app/panicers"
)

// Playtime struct represent a period someone played on a Machine with a set
// Tariff
type Playtime struct {
	ID      int64
	Tariff  Tariff
	Machine Machine
	Start   time.Time
	End     *time.Time
}

func createPlaytimesTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS playtimes (
			start DATETIME DEFAULT CURRENT_TIMESTAMP,
			end DATETIME DEFAULT NULL,
			tariff_id INTEGER NOT NULL,
			machine_id INTEGER NOT NULL,
			FOREIGN KEY(tariff_id) REFERENCES tariff(rowid)
			FOREIGN KEY(machine_ID) REFERENCES machines(rowid)
		);`)
	panicers.WrapAndPanicIfErr(err, "Could not create Playtimes table")
}

// GetPlaytimeByID from db. Panics on failure.
func GetPlaytimeByID(id int64, db *sql.DB) *Playtime {
	pt := scanPlaytime(
		db.QueryRow("SELECT rowid, * FROM playtimes WHERE rowid = ?;", id),
	)
	t := GetTariffByID(pt.Tariff.ID, db)
	panicers.PanicOn(
		t == nil,
		errors.Errorf("Could not retrieve tariff for playtime: %d", pt.ID),
	)

	pt.Tariff = *t
	m := GetMachineByID(pt.Machine.ID, db)
	panicers.PanicOn(
		t == nil,
		errors.Errorf("Could not retrieve machine for playtime: %d", pt.ID),
	)

	pt.Machine = *m
	return pt
}

// Insert playtime in DB, panics on failure.
func (p *Playtime) Insert(db *sql.DB) (Playtime, error) {
	res, err := db.Exec(
		"INSERT INTO playtimes (tariff_id, machine_id) VALUES (?, ?);",
		p.Tariff.ID,
		p.Machine.ID,
	)
	if database.IsConstraintViolationError(err) {
		return Playtime{},
			errors.Wrapf(
				err,
				"Machine with ID %d or Tariff with ID %d does not exists",
				p.Machine.ID,
				p.Tariff.ID,
			)
	}
	panicers.WrapAndPanicIfErr(err, "Could not insert playtime with machine id: %d")

	id, err := res.LastInsertId()
	panicers.WrapAndPanicIfErr(err, "Could not retrieve ID after playtime insert")

	return Playtime{
		ID:      id,
		Start:   time.Now(),
		Machine: p.Machine,
		Tariff:  p.Tariff,
	}, nil
}

// EndPlaytime sets p's End time to now, panics on failure
func (p *Playtime) EndPlaytime(db *sql.DB) Playtime {
	_, err := db.Exec(
		"UPDATE playtimes SET end = CURRENT_TIMESTAMP WHERE rowid = ?;",
		p.ID,
	)
	panicers.WrapAndPanicIfErr(err, "Could not update end for playtime")

	pt := GetPlaytimeByID(p.ID, db)
	panicers.PanicOn(
		pt == nil,
		errors.New("Could not retrieve playtime after insert"),
	)
	return *pt
}

// CalculatePrice for this playtime with its tariff. If Playtime does not have
// an end time then it will use the current time to calculate.
func (p *Playtime) CalculatePrice() int64 {
	var end time.Time
	if p.End != nil {
		end = *p.End
	} else {
		end = time.Now()
	}
	ppu := int64(p.Tariff.PricePerUnit)
	dur := int64(end.Sub(p.Start))
	return int64(math.Ceil(float64(ppu*dur) / float64(p.Tariff.UnitSize)))
}

func scanPlaytime(row scanner) *Playtime {
	pt := &Playtime{}
	err := row.Scan(
		&pt.ID,
		&pt.Start,
		pt.End,
		&pt.Tariff.ID,
		&pt.Machine.ID,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	panicers.WrapAndPanicIfErr(err, "Could not scan playtime from row")

	return pt
}
