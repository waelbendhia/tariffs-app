package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/waelbendhia/tariffs-app/types"
)

func main() {

	db := openSQLite("./foo.db")
	types.CreateTariffsTable(db)
	t1 := types.Tariff{PricePerUnit: 3, UnitSize: 3}
	t1.Insert(db)
	t2 := types.Tariff{PricePerUnit: 4, UnitSize: 4}
	t2.Insert(db)
	t3 := types.GetTariffLatest(db)
	log.Print(t3)
	log.Print(types.GetAllTariffs(db))
	// Initialize GTK without parsing any command line arguments.
	gtk.Init(nil)

	// Create a new toplevel window, set its title, and connect it to the
	// "destroy" signal to exit the GTK main loop when it is destroyed.
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Simple Example")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a new label widget to show in the window.
	l, err := gtk.LabelNew("Hello world!")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}

	// Add the label to the window.
	win.Add(l)

	// Set the default window size.
	win.SetDefaultSize(800, 600)

	// Recursively show all widgets contained in this window.
	win.ShowAll()

	// Begin executing the GTK main loop.  This blocks until
	// gtk.MainQuit() is run.
	gtk.Main()
}
func openSQLite(path string) *sql.DB {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		panic(err)
	}
	return db
}
