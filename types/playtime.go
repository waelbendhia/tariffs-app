package types

import (
	"database/sql"
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
		db.QueryRow(
			`SELECT
				p.rowid, p.*,
				t.rowid, t.*,
				m.rowid, m.*
			FROM 
				playtimes p 
				INNER JOIN
					tariffs t  ON p.tariff_id = t.rowid 
				INNER JOIN
					machines m ON p.machine_id = m.rowid 
			WHERE p.rowid = ?;`,
			id,
		),
	)
	return pt
}

func GetOpenPlaytimeByMachineID(db *sql.DB, id int64) *Playtime {
	for _, p := range GetPlaytimeByMachineID(db, id) {
		if p.End == nil {
			return &p
		}
	}
	return nil
}

func GetPlaytimeByMachineID(db *sql.DB, id int64) []Playtime {
	rows, err := db.Query(
		`SELECT
			p.rowid, p.*,
			t.rowid, t.*,
			m.rowid, m.*
		FROM 
			playtimes p
		INNER JOIN tariffs t  ON p.tariff_id = t.rowid 
		INNER JOIN machines m ON p.machine_id = m.rowid 
		WHERE machine_id = ?;`,
		id,
	)

	panicers.WrapAndPanicIfErr(err, "Could not get playtimes by machine id: %d", id)

	defer func() {
		panicers.WrapAndPanicIfErr(rows.Close(), "Error while closing rows")
	}()

	var pts []Playtime
	for rows.Next() {
		pts = append(pts, *scanPlaytime(rows))
	}

	panicers.WrapAndPanicIfErr(rows.Err(), "Error while querying for playtimes")

	return pts
}

func GetPlaytimes(
	db *sql.DB,
	machineID *int64,
	minDate *time.Time,
	maxDate *time.Time,
) []Playtime {
	rows, err := db.Query(
		`SELECT
				p.rowid, p.*,
				t.rowid, t.*,
				m.rowid, m.*
			FROM
				playtimes p
			INNER JOIN tariffs t  ON p.tariff_id = t.rowid
			INNER JOIN machines m ON p.machine_id = m.rowid
			WHERE
				p.end IS NOT NULL
			AND
				CASE WHEN :machineID IS NOT NULL 
					THEN :machineID = p.machine_id
					ELSE 1
				END
			AND
				CASE WHEN :minDate IS NOT NULL
					THEN :minDate < p.start
					ELSE 1
				END
			AND 
				CASE WHEN :maxDate IS NOT NULL
				  THEN :maxDate > p.end
					ELSE 1
				END;`,
		sql.Named("machineID", machineID),
		sql.Named("minDate", minDate),
		sql.Named("maxDate", maxDate),
	)
	panicers.WrapAndPanicIfErr(err, "Could not search playtimes")

	defer func() {
		panicers.WrapAndPanicIfErr(rows.Close(), "Error while closing rows")
	}()

	var pts []Playtime
	for rows.Next() {
		pts = append(pts, *scanPlaytime(rows))
	}

	panicers.WrapAndPanicIfErr(rows.Err(), "Error while querying for playtimes")
	return pts
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

// EndAllOpenPlaytimes panics on failure
func EndAllOpenPlaytimes(db *sql.DB) {
	_, err := db.Exec(
		"UPDATE playtimes SET end = CURRENT_TIMESTAMP WHERE end IS NULL;",
	)

	panicers.WrapAndPanicIfErr(err, "Failed to end all open playtimes")
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
	var (
		ppu = int64(p.Tariff.PricePerUnit)
		dur = int64(end.Sub(p.Start) / p.Tariff.UnitSize)
	)
	if end.Sub(p.Start)%p.Tariff.UnitSize > 0 {
		dur++
	}
	return ppu * dur
}

func scanPlaytime(row scanner) *Playtime {
	var delFlag int
	// var endTime pq.NullTime
	pt := &Playtime{}
	err := row.Scan(
		&pt.ID,
		&pt.Start,
		&pt.End,
		&pt.Tariff.ID,
		&pt.Machine.ID,
		&pt.Tariff.ID,
		&pt.Tariff.PricePerUnit,
		&pt.Tariff.UnitSize,
		&pt.Tariff.CreatedAt,
		&pt.Machine.ID,
		&pt.Machine.Name,
		&delFlag,
	)
	if err == sql.ErrNoRows {
		return nil
	}
	panicers.WrapAndPanicIfErr(err, "Could not scan playtime from row")
	pt.Machine.deleted = delFlag != 0
	// if endTime.Valid {
	// 	pt.End = &endTime.Time
	// 	}else{
	// 		pt.End = &endTime.Time
	// }
	return pt
}
