package elements

import (
	"fmt"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/waelbendhia/tariffs-app/types"
)

func newMachinesElement(app machineCRUDERTimer) (*widgets.QGroupBox, func(t *types.Tariff)) {
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
		addMachineInput                    = widgets.NewQPlainTextEdit(nil)
		addMachineConfirm                  = newButton("Confirmer")
		// Machines list
		machineScroller                              = widgets.NewQScrollArea(nil)
		machinesListHolder, machinesListHolderLayout = newVBox()
		// toggleAddButton sets the enabled state of a the add button depending on
		// the machine name input's value
		toggleAddButton = func() {
			addMachineConfirm.SetEnabled(addMachineInput.ToPlainText() != "")
		}
		// startButtons holds a slice of references to the startTimer buttons
		// to all machines
		startButtons = make(chan []*widgets.QPushButton, 1)
		// withStartButtons provides a close to access startButtons' value
		withStartButtons = func(f func([]*widgets.QPushButton) []*widgets.QPushButton) {
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
				mSpacer        = widgets.NewQSpacerItem(
					0, 0,
					widgets.QSizePolicy__Expanding,
					widgets.QSizePolicy__Expanding,
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
			mBoxLayout.QLayout.AddWidget(mLabel)
			mBoxLayout.QLayout.AddWidget(mTimer)
			mBoxLayout.AddSpacerItem(mSpacer)
			mBoxLayout.QLayout.AddWidget(mStartTimer)
			mBoxLayout.QLayout.AddWidget(mEndTimer)
			mBoxLayout.QLayout.AddWidget(mDelete)

			mStartTimer.SetEnabled(getHasTariff() && getPlayTime() == nil)
			toggleButtons()
			mBox.SetFixedHeight(inputHeight)

			// set up button actions
			mStartTimer.ConnectClicked(func(_ bool) {
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
			mEndTimer.ConnectClicked(func(_ bool) {
				withPlayTime(func(pt *types.Playtime) *types.Playtime {
					app.End(pt.ID)
					return nil
				})
				toggleButtons()
			})
			mDelete.ConnectClicked(func(_ bool) {
				app.DeleteMachine(m)
				mBox.DeleteLater()
			})

			mBox.ConnectDestroyQObject(func() {
				close(playtimeChan)
				withStartButtons(func(bs []*widgets.QPushButton) []*widgets.QPushButton {
					for i, b := range bs {
						if b == mStartTimer {
							bs = append(bs[:i], bs[i+1:]...)
							break
						}
					}
					return bs
				})
			})
			withStartButtons(func(bs []*widgets.QPushButton) []*widgets.QPushButton {
				bs = append(bs, mStartTimer)
				return bs
			})
			machinesListHolderLayout.QLayout.AddWidget(mBox)
		}
		addMachine = func(b bool) {
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
	addMachineBoxLayout.QLayout.AddWidget(addMachineLabel)
	addMachineBoxLayout.QLayout.AddWidget(addMachineInput)
	addMachineBoxLayout.QLayout.AddWidget(addMachineConfirm)

	addMachineInput.ConnectTextChanged(toggleAddButton)

	toggleAddButton()
	addOnEnterHandler(addMachineInput, func() { addMachine(false) })
	addMachineInput.SetTabChangesFocus(true)

	addMachineConfirm.ConnectClicked(addMachine)

	for _, m := range app.GetMachines() {
		insertMachine(m)
	}

	machinesListHolderLayout.SetAlign(core.Qt__AlignTop)
	machineScroller.SetWidget(machinesListHolder)
	machineScroller.SetWidgetResizable(true)
	addMachineBox.SetMaximumHeight(inputHeight)

	rootBoxLayout.QLayout.AddWidget(addMachineBox)
	rootBoxLayout.QLayout.AddWidget(machineScroller)

	return rootBox, func(t *types.Tariff) {
		setHasTariff(t != nil)
		withStartButtons(func(bs []*widgets.QPushButton) []*widgets.QPushButton {
			for _, b := range bs {
				b.SetEnabled(t != nil)
			}
			return bs
		})
	}
}
