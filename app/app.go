package app

import (
	"database/sql"
	"log"

	"github.com/andlabs/ui"
	"github.com/waelbendhia/tariffs-app/app/elements"
	"github.com/waelbendhia/tariffs-app/database"
	"github.com/waelbendhia/tariffs-app/types"
)

type App struct{ db *sql.DB }

func initApp() *App {
	log.Println("Opening database")
	db := database.OpenSQLite("./app.db")
	log.Println("Initializing database")
	types.InitializeDB(db)

	return &App{db}
}

func (a *App) close() {
	log.Println("Closing database")
	types.CloseRemaining(a.db)
	if err := a.db.Close(); err != nil {
		log.Println("Failed to close DB")
	}
}

func Start() {
	app := initApp()
	err := ui.Main(func() {
		w := elements.MainWindow(app)
		w.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			app.close()
			return true
		})
		w.Show()
	})

	if err != nil {
		panic(err)
	}
}

func (a *App) GetTariff() *types.Tariff {
	return types.GetLatestTariff(a.db)
}

func (a *App) SetTariff(t types.Tariff) types.Tariff {
	log.Printf("Updating tarrif to: %s", t.String())
	return t.Insert(a.db)
}

func (a *App) AddMachine(m types.Machine) types.Machine {
	return m.Insert(a.db)
}

func (a *App) GetMachines() []types.Machine {
	return types.GetAllMachines(a.db)
}

func (a *App) DeleteMachine(m types.Machine) types.Machine {
	return m.Delete(a.db)
}

func (a *App) UpdateMachine(m types.Machine) types.Machine {
	return m.Update(a.db)
}
