package elements

import (
	"sync"

	"github.com/andlabs/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

func newMachinesElement(app machineCRUDER) *ui.Box {
	var (
		hasChildren     = false
		hasChildrenLock sync.Mutex
		rootBox         = ui.NewVerticalBox()
		label           = ui.NewLabel("Machines:")
		// Add machine elements
		addMachineBox     = ui.NewHorizontalBox()
		addMachineLabel   = ui.NewLabel("Ajouter:")
		addMachineInput   = ui.NewEntry()
		addMachineConfirm = ui.NewButton("Confirmer")
		seperator         = ui.NewHorizontalSeparator()
		// Machines list
		machinesListHolder = ui.NewHorizontalBox()
	)

	// Set up add machine elements
	addMachineBox.Append(addMachineLabel, false)
	addMachineBox.Append(addMachineInput, true)
	addMachineBox.Append(addMachineConfirm, false)

	addMachineBox.SetPadded(true)

	addMachineInput.OnChanged(func(i *ui.Entry) {
		if i.Text() != "" {
			addMachineConfirm.Enable()
		} else {
			addMachineConfirm.Disable()
		}
	})

	addMachineConfirm.Disable()

	addMachineConfirm.OnClicked(func(_ *ui.Button) {
		app.AddMachine(types.Machine{Name: addMachineInput.Text()})
		addMachineInput.SetText("")
		addMachineConfirm.Disable()
		redrawMachines(app, machinesListHolder, &hasChildrenLock, &hasChildren)
	})

	redrawMachines(app, machinesListHolder, &hasChildrenLock, &hasChildren)

	rootBox.Append(label, false)
	rootBox.Append(addMachineBox, false)
	rootBox.Append(seperator, false)
	rootBox.Append(machinesListHolder, true)

	rootBox.SetPadded(true)

	return rootBox
}

func redrawMachines(
	app machineCRUDER,
	machinesListHolder *ui.Box,
	hasChildrenLock *sync.Mutex,
	hasChildren *bool,
) {
	machinesListBox := ui.NewVerticalBox()
	for _, m := range app.GetMachines() {
		var (
			mBox    = ui.NewHorizontalBox()
			mLabel  = ui.NewLabel(m.Name)
			mDelete = ui.NewButton("Suprimer")
		)
		mBox.Append(mLabel, true)
		mBox.Append(mDelete, false)
		mDelete.OnClicked(func(_ *ui.Button) {
			app.DeleteMachine(m)
			redrawMachines(app, machinesListHolder, hasChildrenLock, hasChildren)
		})
		machinesListBox.Append(mBox, false)
	}
	machinesListBox.SetPadded(true)
	hasChildrenLock.Lock()
	if *hasChildren {
		machinesListHolder.Delete(0)
	} else {
		*hasChildren = true
	}
	machinesListHolder.Append(machinesListBox, true)
	hasChildrenLock.Unlock()
}
