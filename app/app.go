package app

import (
	"database/sql"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/therecipe/qt/widgets"
	"github.com/waelbendhia/tariffs-app/app/elements"
	"github.com/waelbendhia/tariffs-app/database"
	"github.com/waelbendhia/tariffs-app/panicers"
	"github.com/waelbendhia/tariffs-app/types"
)

// App holds the database for the application and methods for interacting with
// said database
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

// Start the app, close when receive on done
func Start(done chan struct{}) {
	app := initApp()

	appUI := widgets.NewQApplication(0, nil)

	mainWindow := elements.MainWindow(app)
	go func() {
		<-done
		mainWindow.DeleteLater()
	}()
	mainWindow.Show()

	appUI.Exec()

	app.close()
}

// GetTariff if set returns nil otherwise
func (a *App) GetTariff() *types.Tariff {
	return types.GetLatestTariff(a.db)
}

// SetTariff to t
func (a *App) SetTariff(t types.Tariff) types.Tariff {
	return t.Insert(a.db)
}

// AddMachine to database
func (a *App) AddMachine(m types.Machine) types.Machine {
	return m.Insert(a.db)
}

// GetMachines that are not deleted
func (a *App) GetMachines() []types.Machine {
	return types.GetAllMachines(a.db)
}

// DeleteMachine set machine m to deleted
func (a *App) DeleteMachine(m types.Machine) types.Machine {
	return m.Delete(a.db)
}

// End playtime with id
func (a *App) End(id int64) types.Playtime {
	pt := types.Playtime{ID: id}
	return pt.EndPlaytime(a.db)
}

// GetOpenPlayTime for machine with id
func (a *App) GetOpenPlayTime(id int64) *types.Playtime {
	return types.GetOpenPlaytimeByMachineID(a.db, id)
}

// Start a playtime for machine with id
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

// SearchPlaytimes matching criteria
func (a *App) SearchPlaytimes(
	mID *int64,
	minDate *time.Time,
	maxDate *time.Time,
) []types.Playtime {
	return types.GetPlaytimes(a.db, mID, minDate, maxDate)
}
