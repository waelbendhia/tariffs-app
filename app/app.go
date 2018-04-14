package app

import (
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/visualfc/goqt/ui"
	"github.com/waelbendhia/tariffs-app/app/elements"
	"github.com/waelbendhia/tariffs-app/database"
	"github.com/waelbendhia/tariffs-app/panicers"
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

func Start(done chan bool) {
	app := initApp()

	ui.Run(func() {
		mainWindow := elements.MainWindow(app)
		go func() {
			<-done
			mainWindow.Delete()
		}()
		mainWindow.Show()
	})

	app.close()
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

func (a *App) End(id int64) types.Playtime {
	pt := types.Playtime{ID: id}
	return pt.EndPlaytime(a.db)
}

func (a *App) GetOpenPlayTime(id int64) *types.Playtime {
	return types.GetOpenPlaytimeByMachineID(a.db, id)
}

func (a *App) Start(machineID int64) types.Playtime {
	t := a.GetTariff()
	panicers.PanicOn(
		t == nil,
		errors.New("Tried to start play session with no tariff"),
	)
	pt := types.Playtime{}
	pt.Machine.ID = machineID
	pt.Tariff = *t
	ptI, err := pt.Insert(a.db)
	panicers.WrapAndPanicIfErr(err, "Could not insert pt: %v", pt)
	return ptI
}

func (a *App) SearchPlaytimes(
	mID *int64,
	minDate *time.Time,
	maxDate *time.Time,
) []types.Playtime {
	return types.GetPlaytimes(a.db, mID, minDate, maxDate)
}
