package types

import (
	"database/sql"
	"math"
	"time"

	"github.com/pkg/errors"
)

type Playtime struct {
	ID      int64
	Tariff  Tariff
	Machine Machine
	Start   time.Time
	End     *time.Time
}

func CreatePlaytimesTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS playtimes (
			start DATETIME DEFAULT CURRENT_TIMESTAMP,
			end DATETIME DEFAULT NULL,
			tariff_id INTEGER NOT NULL,
			machine_id INTEGER NOT NULL,
			FOREIGN KEY(tariff_id) REFERENCES tariff(rowid)
			FOREIGN KEY(machine_ID) REFERENCES machines(rowid)
		);`)
	if err != nil {
		panic(errors.Wrap(err, "Could not create Playtimes table"))
	}
}

func GetPlaytimeByID(id int64, db *sql.DB) *Playtime {
	pt := scanPlaytime(
		db.QueryRow("SELECT rowid, * FROM playtimes WHERE rowid = ?;", id),
	)
	t := GetTariffByID(pt.Tariff.ID, db)
	if t == nil {
		panic(
			errors.Errorf("Could not retrieve tariff for playtime: %d", pt.ID),
		)
	}
	pt.Tariff = *t
	m := GetMachineByID(pt.Machine.ID, db)
	if t == nil {
		panic(
			errors.Errorf("Could not retrieve machine for playtime: %d", pt.ID),
		)
	}
	pt.Machine = *m
	return pt
}

func (p *Playtime) Insert(db *sql.DB) Playtime {
	res, err := db.Exec(
		"INSERT INTO playtimes (machine_id) VALUES (?);",
		p.Machine.ID,
	)
	if err != nil {
		panic(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}
	return Playtime{
		ID:      id,
		Start:   time.Now(),
		Machine: p.Machine,
		Tariff:  p.Tariff,
	}
}

func (p *Playtime) EndTime(db *sql.DB) Playtime {
	_, err := db.Exec(
		"UPDATE playtimes SET end = CURRENT_TIMESTAMP WHERE rowid = ?;",
		p.ID,
	)
	if err != nil {
		panic(err)
	}
	pt := GetPlaytimeByID(p.ID, db)
	if pt == nil {
		panic("Could not retrieve playtime after insert")
	}
	return *pt
}

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
		pt.EndTime,
		&pt.Tariff.ID,
		&pt.Machine.ID,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(err)
	}
	return pt
}
