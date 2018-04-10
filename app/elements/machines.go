package elements

import (
	"log"
	"sync"
	"time"

	"github.com/andlabs/ui"
	"github.com/waelbendhia/tariffs-app/types"
)

type machineSignal int

const (
	KILL_SIG  = 0
	START_SIG = 1
	END_SIG   = 2
)

func newMachinesElement(app machineCRUDERTimer) *ui.Box {
	var (
		mElemsLookUp = make(map[int64]int)
		mElemsLock   sync.Mutex
		rootBox      = ui.NewVerticalBox()
		label        = ui.NewLabel("Machines:")
		// Add machine elements
		addMachineBox     = ui.NewHorizontalBox()
		addMachineLabel   = ui.NewLabel("Ajouter:")
		addMachineInput   = ui.NewEntry()
		addMachineConfirm = ui.NewButton("Confirmer")
		seperator         = ui.NewHorizontalSeparator()
		// Machines list
		machinesListHolder = ui.NewVerticalBox()
		// Helper functions
		toggleAddButton = func() {
			if addMachineInput.Text() != "" {
				addMachineConfirm.Enable()
			} else {
				addMachineConfirm.Disable()
			}
		}
		createMachineElement = func(m types.Machine) (
			*ui.Box,
			chan<- *types.Playtime,
			<-chan machineSignal,
		) {
			var (
				playTimeChan = make(chan *types.Playtime)
				timerChan    = make(chan *time.Time)
				signalChan   = make(chan machineSignal)
				mBox         = ui.NewHorizontalBox()
				mLabel       = ui.NewLabel(m.Name)
				mTimer       = ui.NewLabel("")
				mStartTimer  = ui.NewButton("Commencer")
				mEndTimer    = ui.NewButton("Arreter")
				mDelete      = ui.NewButton("Suprimer")
			)
			// This go routine constantly reads from playTimeChan
			// and updates this element accordingly
			go func() {
				for t := range playTimeChan {
					if t != nil {
						mTimer.SetText(
							time.
								Since(t.Start).
								Truncate(time.Second).
								String(),
						)
						mTimer.Show()
						mDelete.Show()
						timerChan <- &t.Start
						// This go routine updates the timer display
						go func() {
							for start := range timerChan {
								if start != nil {
									for {
										select {
										case <-time.Tick(time.Second):
											mTimer.SetText(
												time.
													Since(*start).
													Truncate(time.Second).
													String(),
											)
										case <-timerChan:
											return

										}
									}
								}
							}
						}()
					} else {
						mTimer.Hide()
						mDelete.Hide()
						mStartTimer.Show()
						timerChan <- nil
					}
				}
				close(timerChan)
			}()

			mBox.Append(mLabel, true)

			mEndTimer.OnClicked(func(_ *ui.Button) {
				signalChan <- START_SIG
			})
			mBox.Append(mStartTimer, false)
			mBox.Append(mTimer, true)

			mEndTimer.OnClicked(func(_ *ui.Button) {
				signalChan <- END_SIG
			})
			mBox.Append(mEndTimer, false)

			mDelete.OnClicked(func(_ *ui.Button) {
				signalChan <- KILL_SIG
				close(signalChan)
			})
			mBox.Append(mDelete, false)
			return mBox, playTimeChan, signalChan
		}
		machineWatcherRoutine = func(
			ptChan chan<- *types.Playtime,
			signalChan <-chan machineSignal,
			ID int64,
			mElem *ui.Box,
		) {
			for sig := range signalChan {
				switch sig {
				case KILL_SIG:
					close(ptChan)
					defer mElemsLock.Unlock()
					mElemsLock.Lock()
					ind := mElemsLookUp[ID]
					delete(mElemsLookUp, ID)
					for k, v := range mElemsLookUp {
						if v > ind {
							mElemsLookUp[k] = mElemsLookUp[k] - 1
						}
					}
					machinesListHolder.Delete(ind)
					mElem.Destroy()
					app.DeleteMachine(types.Machine{ID: ID})
					return
				case START_SIG:
					pt := app.Start(ID)
					ptChan <- &pt
				default:
					pt := app.GetOpenPlayTime(ID)
					if pt != nil {
						app.End(pt.ID)
						ptChan <- nil
					}
				}
			}
		}
		insertMachine = func(i int, m types.Machine) {
			mElem, ptChan, signalChan := createMachineElement(m)
			machinesListHolder.Append(mElem, true)
			go machineWatcherRoutine(ptChan, signalChan, m.ID, mElem)
			mElemsLock.Lock()
			mElemsLookUp[m.ID] = i
			mElemsLock.Unlock()
		}
		submitMachine = func() {
			m := app.AddMachine(types.Machine{Name: addMachineInput.Text()})
			addMachineInput.SetText("")
			toggleAddButton()
			mElemsLock.Lock()
			highest := 0
			for _, v := range mElemsLookUp {
				if v > highest {
					highest = v
				}
			}
			mElemsLock.Unlock()
			insertMachine(highest+1, m)
		}
	)

	// Set up add machine elements
	addMachineBox.Append(addMachineLabel, false)
	addMachineBox.Append(addMachineInput, true)
	addMachineBox.Append(addMachineConfirm, false)

	addMachineBox.SetPadded(true)

	addMachineInput.OnChanged(func(i *ui.Entry) { toggleAddButton() })

	toggleAddButton()

	addMachineConfirm.OnClicked(func(_ *ui.Button) { submitMachine() })

	for i, m := range app.GetMachines() {
		insertMachine(i, m)
	}

	rootBox.Append(label, false)
	rootBox.Append(addMachineBox, false)
	rootBox.Append(seperator, false)
	rootBox.Append(machinesListHolder, true)

	rootBox.SetPadded(true)

	return rootBox
}

func redrawMachines(
	app machineCRUDERTimer,
	machinesListHolder *ui.Box,
	hasChildrenLock *sync.Mutex,
	hasChildren *bool,
) {
	machinesListBox := ui.NewVerticalBox()
	for _, m := range app.GetMachines() {
		func(m types.Machine) {
			var (
				playtime    = app.GetOpenPlayTime(m.ID)
				mBox        = ui.NewHorizontalBox()
				mLabel      = ui.NewLabel(m.Name)
				mTimer      = ui.NewLabel("")
				mStartTimer = ui.NewButton("Commencer")
				mEndTimer   = ui.NewButton("Arreter")
				mDelete     = ui.NewButton("Suprimer")
			)
			mBox.Append(mLabel, true)
			mEndTimer.OnClicked(func(_ *ui.Button) {
				app.End(playtime.ID)
				redrawMachines(app, machinesListHolder, hasChildrenLock, hasChildren)
			})
			mStartTimer.OnClicked(func(_ *ui.Button) {
				app.Start(m.ID)
				redrawMachines(app, machinesListHolder, hasChildrenLock, hasChildren)
			})
			if playtime == nil {
				mBox.Append(mStartTimer, false)
			} else {
				go func() {
					for _ = range time.Tick(time.Second) {
						playtime := app.GetOpenPlayTime(m.ID)
						if playtime != nil {
							mTimer.SetText(
								((time.Since(playtime.Start) / time.Second) * time.Second).
									String(),
							)
							log.Printf("Timing %d", playtime.ID)
						} else {
							log.Println("Exiting")
							return
						}
					}
				}()
				mBox.Append(mTimer, true)
				mBox.Append(mEndTimer, false)
			}
			mBox.Append(mDelete, false)
			mDelete.OnClicked(func(_ *ui.Button) {
				app.DeleteMachine(m)
				redrawMachines(app, machinesListHolder, hasChildrenLock, hasChildren)
			})
			machinesListBox.Append(mBox, false)
		}(m)
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
