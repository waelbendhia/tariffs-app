package types

import (
	"database/sql"

	panicers "github.com/waelbendhia/tariffs-app/panicers"
)

// Machine struct
type Machine struct {
	ID      int64
	Name    string
	deleted bool
}

func createMachinesTable(db *sql.DB) {
	_, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS 
			machines (
				name TEXT NOT NULL,
				deleted INTEGER NOT NULL DEFAULT 0
			);`,
	)
	panicers.WrapAndPanicIfErr(err, "Could not create machines table")
}

// Insert Machine m in db, panics on failure.
func (m *Machine) Insert(db *sql.DB) Machine {
	res, err := db.Exec("INSERT INTO machines (name) VALUES (?);", m.Name)
	panicers.WrapAndPanicIfErr(err, "Could not insert machine: %s", m.Name)

	id, err := res.LastInsertId()
	panicers.WrapAndPanicIfErr(err, "Could not retrieve ID for machine: %s", m.Name)

	return Machine{id, m.Name, false}
}

// Delete machine m from db, panics on failure
func (m *Machine) Delete(db *sql.DB) Machine {
	_, err := db.Exec("UPDATE machines SET deleted = 1 WHERE rowid = ?;", m.ID)
	panicers.WrapAndPanicIfErr(err, "Could not delete machine: %d", m.ID)
	return Machine{m.ID, m.Name, true}
}

// Update machine m in db, panics on failure
func (m *Machine) Update(db *sql.DB) Machine {
	_, err := db.Exec("UPDATE machines SET name = ? WHERE rowid = ?;", m.Name)
	panicers.WrapAndPanicIfErr(err, "Could not update machine: %d", m.ID)
	return Machine{m.ID, m.Name, true}
}

// GetAllMachines from db, panics on failure.
func GetAllMachines(db *sql.DB) []Machine {
	rows, err := db.Query("SELECT rowid, * FROM machines WHERE deleted = 0;")
	panicers.WrapAndPanicIfErr(err, "Could not query machines table")

	defer func() {
		panicers.WrapAndPanicIfErr(rows.Close(), "Error while closing rows")
	}()

	var machines []Machine

	for rows.Next() {
		machines = append(machines, *scanMachine(rows))
	}
	panicers.WrapAndPanicIfErr(rows.Err(), "Error while querying for machines")
	return machines
}

// GetMachineByID from db, if not found returns nil. Panics on failure.
func GetMachineByID(id int64, db *sql.DB) *Machine {
	return scanMachine(
		db.QueryRow(
			"SELECT rowid, * FROM machines WHERE rowid = ? AND deleted = 0;",
			id,
		),
	)
}

func scanMachine(row scanner) *Machine {
	var (
		machine = &Machine{}
		delFlag int
	)
	err := row.Scan(&machine.ID, &machine.Name, &delFlag)
	if err == sql.ErrNoRows {
		return nil
	}
	panicers.WrapAndPanicIfErr(err, "Could not scan machine from row")
	machine.deleted = delFlag != 0
	return machine
}
