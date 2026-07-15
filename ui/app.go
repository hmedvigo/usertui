package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppContext struct {
	App         *tview.Application
	MainGrid    *tview.Grid
	InnerFlex   *tview.Flex
	HeaderTabs  *tview.Flex
	UsersBtn    *tview.Button
	GroupsBtn   *tview.Button
	BodyPages   *tview.Pages
	CurrentTab  int
	UsersPanel  *TabPanel
	GroupsPanel *TabPanel
}

func NewApp() *AppContext {
	ctx := &AppContext{
		App:        tview.NewApplication(),
		MainGrid:   tview.NewGrid(),
		InnerFlex:  tview.NewFlex().SetDirection(tview.FlexRow),
		HeaderTabs: tview.NewFlex().SetDirection(tview.FlexColumn),
		UsersBtn:   tview.NewButton(" USERS "),
		GroupsBtn:  tview.NewButton(" GROUPS "),
		BodyPages:  tview.NewPages(),
		CurrentTab: 0,
	}

	ctx.setupLayout()
	ctx.setupKeybindings()
	ctx.TabVisuals()

	return ctx
}

func (ctx *AppContext) setupLayout() {
	// RETRO FIX: True matrix black backgrounds with classic phosphor green borders
	ctx.UsersPanel = NewTabPanel("users",
		func(actionIdx int) {
			// Handle user action
		},
		func() {
			ctx.App.Stop()
		},
	)
	ctx.GroupsPanel = NewTabPanel("groups",
		func(actionIdx int) {
			// Handle group action
		},
		func() {
			ctx.App.Stop()
		},
	)
	ctx.BodyPages.AddPage("users", ctx.UsersPanel.MainFlex, true, true)
	ctx.BodyPages.AddPage("groups", ctx.GroupsPanel.MainFlex, true, false)

	ctx.UsersBtn.SetSelectedFunc(func() {
		ctx.CurrentTab = 0
		ctx.BodyPages.SwitchToPage("users")
		ctx.TabVisuals()
	})

	ctx.GroupsBtn.SetSelectedFunc(func() {
		ctx.CurrentTab = 1
		ctx.BodyPages.SwitchToPage("groups")
		ctx.TabVisuals()
	})

	// RETRO FIX: Keep the entire header backdrop flat black
	ctx.HeaderTabs.
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 0, 0, false).
		AddItem(ctx.UsersBtn, 9, 0, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorRed), 0, 0, false).
		AddItem(ctx.GroupsBtn, 10, 0, false)

	ctx.InnerFlex.AddItem(ctx.HeaderTabs, 1, 1, true)
	ctx.InnerFlex.AddItem(ctx.BodyPages, 0, 1, false)

	ctx.MainGrid.
		SetColumns(0, 70, 0).
		SetRows(0, 22, 0).
		AddItem(ctx.InnerFlex, 1, 1, 1, 1, 0, 0, true)
}
func (ctx *AppContext) TabVisuals() {

	activeStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	inactiveStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkGreen)

	// Single helper function to cleanly apply the retro style parameters to any given button
	applyRetroStyle := func(btn *tview.Button, isFieldActive bool) {
		if isFieldActive {
			btn.SetStyle(activeStyle)
			btn.SetLabelColor(tcell.ColorGreen)
			btn.SetLabelColorActivated(tcell.ColorGreen)
		} else {
			btn.SetStyle(inactiveStyle)
			btn.SetLabelColor(tcell.ColorDarkGreen)
			btn.SetLabelColorActivated(tcell.ColorDarkGreen)
		}
		// Both states share a solid, glare-free black background block
		btn.SetBackgroundColorActivated(tcell.ColorBlack)
	}

	// Apply styles using our optimized helper
	applyRetroStyle(ctx.UsersBtn, ctx.CurrentTab == 0)
	applyRetroStyle(ctx.GroupsBtn, ctx.CurrentTab == 1)

	if ctx.CurrentTab == 0 {
		if ctx.UsersPanel != nil && ctx.UsersPanel.List != nil {
			ctx.App.SetFocus(ctx.UsersPanel.List)
			ctx.UsersPanel.List.SetCurrentItem(0) // Reset to first item
		}
	}
	if ctx.CurrentTab == 1 {
		if ctx.GroupsPanel != nil && ctx.GroupsPanel.List != nil {
			ctx.App.SetFocus(ctx.GroupsPanel.List)
			ctx.GroupsPanel.List.SetCurrentItem(0) // Reset to first item
		}
	}
}
func (ctx *AppContext) setupKeybindings() {
	ctx.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentFocus := ctx.App.GetFocus()

		activePanel := ctx.getActivePanel()
		if activePanel != nil {
			// handlePanelNavigation returns true if it consumed the event
			if ctx.handlePanelNavigation(event, currentFocus, activePanel) {
				return nil // Event was consumed, stop processing
			}
		}

		// Tab Navigation Processing
		if event.Key() == tcell.KeyTab {
			ctx.CurrentTab = (ctx.CurrentTab + 1) % 2
			ctx.switchTab(ctx.CurrentTab)
			return nil
		}

		// CRUCIAL FIX: Let UP/DOWN keys pass through so tview can scroll list items
		return event
	})
}
func (ctx *AppContext) getActivePanel() *TabPanel {
	switch ctx.CurrentTab {
	case 0:
		return ctx.UsersPanel
	case 1:
		return ctx.GroupsPanel
	default:
		return nil
	}
}
func (ctx *AppContext) switchTab(tabIndex int) {
	switch tabIndex {
	case 0:
		ctx.BodyPages.SwitchToPage("users")
	case 1:
		ctx.BodyPages.SwitchToPage("groups")
	}
	ctx.TabVisuals()
}
func (ctx *AppContext) handlePanelNavigation(event *tcell.EventKey, currentFocus tview.Primitive, panel *TabPanel) bool {
	// Guard clause: make sure event isn't nil before doing anything
	if event == nil {
		return false
	}

	// Get button references
	proceedBtn := panel.Buttons.GetItem(1)
	cancelBtn := panel.Buttons.GetItem(3)
	currIdx := panel.List.GetCurrentItem()
	maxIdx := panel.List.GetItemCount() - 1

	// Find the last selectable item (even index)
	lastSelectable := maxIdx
	if lastSelectable%2 != 0 {
		lastSelectable--
	}
	// Using a switch statement prevents multiple blocks from running sequentially
	switch event.Key() {
	case tcell.KeyUp:
		if currentFocus == panel.List {
			if currIdx == 0 {
				return true
			}

			// Move to previous selectable item
			newIdx := currIdx - 1
			if newIdx%2 != 0 {
				newIdx--
			}
			if newIdx >= 0 {
				panel.List.SetCurrentItem(newIdx)
			}
			return true
		} else if currentFocus == proceedBtn || currentFocus == cancelBtn {
			// Move focus to list and set to last selectable item
			panel.List.SetCurrentItem(lastSelectable)
			ctx.App.SetFocus(panel.List)
			// Force redraw of the list to show selection
			panel.List.SetSelectedBackgroundColor(tcell.ColorGreen)
			return true
		}

	case tcell.KeyDown:

		if currentFocus == panel.List {
			if currIdx == lastSelectable {
				// Move focus to proceed button, but keep the list selection visible
				ctx.App.SetFocus(proceedBtn)
				// Keep the current item highlighted by not changing it
				return true
			}

			// Move to next selectable item
			newIdx := currIdx + 1
			if newIdx%2 != 0 {
				newIdx++
			}
			if newIdx <= maxIdx {
				panel.List.SetCurrentItem(newIdx)
			}
			return true
		}

	case tcell.KeyRight:
		if currentFocus == proceedBtn {
			ctx.App.SetFocus(cancelBtn)
			return true
		}

	case tcell.KeyLeft:
		if currentFocus == cancelBtn {
			ctx.App.SetFocus(proceedBtn)
			return true
		}
	default:
		return false
	}
	return false
}

func (ctx *AppContext) Run() error {
	// Set focus to the list right before the application loop kicks off
	go func() {
		ctx.App.QueueUpdate(func() {
			if ctx.UsersPanel != nil && ctx.UsersPanel.List != nil {
				ctx.App.SetFocus(ctx.UsersPanel.List)
				ctx.UsersPanel.List.SetCurrentItem(0)
			}
		})
	}()
	return ctx.App.SetRoot(ctx.MainGrid, true).EnableMouse(false).Run()
}
