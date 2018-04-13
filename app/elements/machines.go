package elements

import (
	"fmt"
	"time"

	"github.com/visualfc/goqt/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

func newMachinesElement(app machineCRUDERTimer) *ui.QGroupBox {
	var (
		rootBox, rootBoxLayout = newVGroupBoxWithTitle("Machines:")
		// Add machine elements
		addMachineBox, addMachineBoxLayout = newHBox()
		addMachineLabel                    = newLabelWithText("Ajouter:")
		addMachineInput                    = ui.NewPlainTextEdit()
		addMachineConfirm                  = ui.NewPushButtonWithTextParent("Confirmer", nil)
		// Machines list
		machineScroller                              = ui.NewScrollArea()
		machinesListHolder, machinesListHolderLayout = newVBox()
		// Helper functions
		toggleAddButton = func() {
			if addMachineInput.ToPlainText() != "" {
				addMachineConfirm.SetEnabled(true)
			} else {
				addMachineConfirm.SetEnabled(false)
			}
		}
		insertMachine = func(m types.Machine) {
			var (
				mBox, mBoxLayout = newHBox()
				playtime         = make(chan *types.Playtime, 1)
				mLabel           = newLabelWithText(m.Name)
				mTimer           = newLabelWithText("")
				mSpacer          = ui.NewSpacerItem(
					0, 0,
					ui.QSizePolicy_Expanding,
					ui.QSizePolicy_Expanding,
				)
				mStartTimer   = ui.NewPushButtonWithTextParent("Commencer", nil)
				mEndTimer     = ui.NewPushButtonWithTextParent("Arreter", nil)
				mDelete       = ui.NewPushButtonWithTextParent("Suprimer", nil)
				toggleButtons = func() {
					pt := <-playtime
					mEndTimer.SetEnabled(pt != nil)
					mStartTimer.SetEnabled(pt == nil)
					mDelete.SetEnabled(pt == nil)
					playtime <- pt
				}
			)

			mLabel.SetStyleSheet(`QLabel {font-weight: 900}`)
			mBoxLayout.AddWidget(mLabel)

			mStartTimer.OnClicked(func() {
				<-playtime
				pt := app.Start(m.ID)
				playtime <- &pt
				mTimer.SetText("0s")
				toggleButtons()
				go func() {
					for _ = range time.Tick(time.Second) {
						pt := <-playtime
						playtime <- pt
						if pt == nil {
							mTimer.SetText("")
							return
						}
						mTimer.SetText(
							time.Since(pt.Start).Truncate(time.Second).String() +
								fmt.Sprintf(" %d Millimes", pt.CalculatePrice()),
						)
					}
				}()
			})

			mBoxLayout.AddWidget(mTimer)
			mBoxLayout.AddSpacerItem(mSpacer)
			mBoxLayout.AddWidget(mStartTimer)

			mEndTimer.OnClicked(func() {
				pt := <-playtime
				app.End(pt.ID)
				playtime <- nil
				toggleButtons()
			})

			mBoxLayout.AddWidget(mEndTimer)

			mDelete.OnClicked(func() {
				<-playtime
				close(playtime)
				app.DeleteMachine(m)
				mBox.Delete()
			})
			playtime <- app.GetOpenPlayTime(m.ID)

			mBoxLayout.AddWidget(mDelete)
			toggleButtons()
			mBox.SetFixedHeight(52)
			machinesListHolderLayout.AddWidget(mBox)
		}
		addMachine = func() {
			if addMachineInput.ToPlainText() != "" {
				m := app.AddMachine(types.Machine{Name: addMachineInput.ToPlainText()})
				addMachineInput.SetPlainText("")
				toggleAddButton()
				insertMachine(m)
			}
		}
	)

	// Set up add machine elements
	addMachineBoxLayout.AddWidget(addMachineLabel)
	addMachineBoxLayout.AddWidget(addMachineInput)
	addMachineBoxLayout.AddWidget(addMachineConfirm)

	addMachineInput.OnTextChanged(toggleAddButton)

	toggleAddButton()

	addMachineInput.InstallEventFilter(newSubmitOnEnterFilter(addMachine))
	addMachineInput.SetTabChangesFocus(true)
	addMachineConfirm.OnClicked(addMachine)

	for _, m := range app.GetMachines() {
		insertMachine(m)
	}

	machinesListHolderLayout.SetAlignment(ui.Qt_AlignTop)
	machineScroller.SetWidget(machinesListHolder)
	machineScroller.SetWidgetResizable(true)
	addMachineBox.SetMaximumHeight(inputHeight)

	rootBoxLayout.AddWidget(addMachineBox)
	rootBoxLayout.AddWidget(machineScroller)

	return rootBox
}
