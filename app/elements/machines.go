package elements

import (
	"fmt"
	"time"

	"github.com/visualfc/goqt/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

func newMachinesElement(app machineCRUDERTimer) (*ui.QGroupBox, func(t *types.Tariff)) {
	var (
		// root widget and layout
		rootBox, rootBoxLayout = newVGroupBoxWithTitle("Machines:")
		// hasTariff holds whether or not the app has a tariff set up
		hasTariff    = make(chan bool, 1)
		getHasTariff = func() bool {
			b := <-hasTariff
			hasTariff <- b
			return b
		}
		setHasTariff = func(b bool) {
			<-hasTariff
			hasTariff <- b
		}
		// Add machine elements
		addMachineBox, addMachineBoxLayout = newHBox()
		addMachineLabel                    = newLabelWithText("Ajouter:")
		addMachineInput                    = ui.NewPlainTextEdit()
		addMachineConfirm                  = newButton("Confirmer")
		// Machines list
		machineScroller                              = ui.NewScrollArea()
		machinesListHolder, machinesListHolderLayout = newVBox()
		// toggleAddButton sets the enabled state of a the add button depending on
		// the machine name input's value
		toggleAddButton = func() {
			addMachineConfirm.SetEnabled(addMachineInput.ToPlainText() != "")
		}
		// startButtons holds a slice of references to the startTimer buttons
		// to all machines
		startButtons = make(chan []*ui.QPushButton, 1)
		// withStartButtons provides a close to access startButtons' value
		withStartButtons = func(f func([]*ui.QPushButton) []*ui.QPushButton) {
			startButtons <- f(<-startButtons)
		}
		insertMachine = func(m types.Machine) {
			var (
				// box to hold all machine UI elements
				mBox, mBoxLayout = newHBox()
				// playtime will hold the current running playtime for this machine
				playtimeChan = make(chan *types.Playtime, 1)
				// convenience functions for accessing the current playtime
				withPlayTime = func(f func(*types.Playtime) *types.Playtime) {
					pt := <-playtimeChan
					playtimeChan <- f(pt)
				}
				getPlayTime = func() *types.Playtime {
					pt := <-playtimeChan
					playtimeChan <- pt
					return pt
				}
				// label and timer
				mLabel, mTimer = newLabelWithText(m.Name), newLabelWithText("")
				mSpacer        = ui.NewSpacerItem(
					0, 0,
					ui.QSizePolicy_Expanding,
					ui.QSizePolicy_Expanding,
				)
				// buttons
				mStartTimer   = newButton("Commencer")
				mEndTimer     = newButton("Arreter")
				mDelete       = newButton("Suprimer")
				toggleButtons = func() {
					pt := getPlayTime()
					mEndTimer.SetEnabled(pt != nil)
					mStartTimer.SetEnabled(getHasTariff() && pt == nil)
					mDelete.SetEnabled(pt == nil)
				}
			)
			// initialize playtime
			playtimeChan <- app.GetOpenPlayTime(m.ID)
			// insert all widgets
			mLabel.SetStyleSheet(`QLabel {font-weight: 900}`)
			mBoxLayout.AddWidget(mLabel)
			mBoxLayout.AddWidget(mTimer)
			mBoxLayout.AddSpacerItem(mSpacer)
			mBoxLayout.AddWidget(mStartTimer)
			mBoxLayout.AddWidget(mEndTimer)
			mBoxLayout.AddWidget(mDelete)

			mStartTimer.SetEnabled(getHasTariff() && getPlayTime() == nil)
			toggleButtons()
			mBox.SetFixedHeight(inputHeight)

			// set up button actions
			mStartTimer.OnClicked(func() {
				withPlayTime(func(_ *types.Playtime) *types.Playtime {
					pt := app.Start(m.ID)
					return &pt
				})
				mTimer.SetText("0s")
				toggleButtons()
				go func() {
					done := make(chan struct{}, 1)
					ticker := time.NewTicker(time.Second)
					for {
						select {
						case <-ticker.C:
							pt := getPlayTime()
							timerString := ""
							if pt == nil {
								ticker.Stop()
								done <- struct{}{}
							} else {
								timerString = time.Since(pt.Start).Truncate(time.Second).String() +
									fmt.Sprintf(" %d Millimes", pt.CalculatePrice())
							}
							mTimer.SetText(timerString)
						case <-done:
							close(done)
							return
						}
					}
				}()
			})
			mEndTimer.OnClicked(func() {
				withPlayTime(func(pt *types.Playtime) *types.Playtime {
					app.End(pt.ID)
					return nil
				})
				toggleButtons()
			})
			mDelete.OnClicked(func() {
				app.DeleteMachine(m)
				mBox.Delete()
			})

			mBox.OnDestroyed(func() {
				close(playtimeChan)
				withStartButtons(func(bs []*ui.QPushButton) []*ui.QPushButton {
					for i, b := range bs {
						if b == mStartTimer {
							bs = append(bs[:i], bs[i+1:]...)
							break
						}
					}
					return bs
				})
			})
			withStartButtons(func(bs []*ui.QPushButton) []*ui.QPushButton {
				bs = append(bs, mStartTimer)
				return bs
			})
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
	// We assume the app has no tariff specifie until update
	hasTariff <- false
	// We initialize startButtons with an empty slice
	startButtons <- nil

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

	return rootBox, func(t *types.Tariff) {
		setHasTariff(t != nil)
		withStartButtons(func(bs []*ui.QPushButton) []*ui.QPushButton {
			for _, b := range bs {
				b.SetEnabled(t != nil)
			}
			return bs
		})
	}
}
