package app

import (
	"database/sql"
	"log"

	"github.com/andlabs/ui"
	"github.com/waelbendhia/tariffs-app/app/elements"
	"github.com/waelbendhia/tariffs-app/database"
	"github.com/waelbendhia/tariffs-app/types"
)

type App struct {
	db *sql.DB
}

func Init() *App {
	log.Println("Opening database")
	db := database.OpenSQLite("./app.db")
	log.Println("Initializing database")
	types.InitializeDB(db)

	return &App{db}
}

func (a *App) Close() {
	log.Println("Closing database")
	types.CloseRemaining(a.db)
	if err := a.db.Close(); err != nil {
		log.Println("Failed to close DB")
	}
}

func Start() {
	app := Init()

	err := ui.Main(func() {
		w := elements.MainWindow()
		w.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			app.Close()
			return true
		})
		t := elements.TariffInput(types.GetLatestTariff(app.db), func(t *types.Tariff, err error) {
			if err != nil {
				log.Println(err)
			} else {
				log.Println(t)
				log.Println(t.Insert(app.db))
			}
		})
		w.SetChild(t)
		w.Show()
	})

	if err != nil {
		panic(err)
	}
}
