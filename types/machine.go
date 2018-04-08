package types

import (
	"database/sql"

	"github.com/pkg/errors"
)

type Machine struct {
	ID   int64
	Name string
}

func CreateMachinesTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS machines (name TEXT NOT NULL );")
	if err != nil {
		panic(errors.Wrap(err, "Could not create machines table"))
	}
}

func (m *Machine) Insert(db *sql.DB) Machine {
	res, err := db.Exec("INSERT INTO machines VALUES (?);", m.Name)
	if err != nil {
		panic(errors.Wrapf(err, "Could not insert machine: %s", m.Name))
	}
	id, err := res.LastInsertId()
	if err != nil {
		panic(errors.Wrapf(err, "Could not retrieve ID for machine: %s", m.Name))
	}
	return Machine{id, m.Name}
}

func GetAllMachines(db *sql.DB) []Machine {
	rows, err := db.Query("SELECT rowid, * FROM machines;")
	if err != nil {
		panic(errors.Wrap(err, "Could not query machines table"))
	}
	var machines []Machine
	for rows.Next() {
		machines = append(machines, *scanMachine(rows))
	}
	return machines
}

func GetMachineByID(id int64, db *sql.DB) *Machine {
	return scanMachine(db.QueryRow("SELECT rowid, * FROM machines WHERE rowid = ?;", id))
}

func scanMachine(row scanner) *Machine {
	machine := &Machine{}
	err := row.Scan(&machine.ID, &machine.Name)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		panic(errors.Wrap(err, "Could not scan machine from row"))
	}
	return machine
}
