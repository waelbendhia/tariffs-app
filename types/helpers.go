package types

import "database/sql"

// InitializeDB creates entity tables in DB
func InitializeDB(db *sql.DB) {
	createMachinesTable(db)
	createPlaytimesTable(db)
	createTariffsTable(db)
}

// CloseRemaining will end all running playtimes
func CloseRemaining(db *sql.DB) {
	EndAllOpenPlaytimes(db)
}
