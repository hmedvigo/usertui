package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type AppContext struct {
	App        *tview.Application
	MainGrid   *tview.Grid
	InnerFlex  *tview.Flex
	HeaderTabs *tview.Flex
	UsersBtn   *tview.Button
	GroupsBtn  *tview.Button
	BodyPages  *tview.Pages
	CurrentTab int
	UsersPanel *UsersTabPanel
	DebugView  *tview.TextView
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
	ctx.UsersPanel = NewUsersTabPanel(
		func(actionIdx int) {
			// Triggered when user hits <Proceed to Action>
			// Inside here, you will route to sub-panels (e.g., actual create user forms)
		},
		func() {
			// Triggered when user hits <Cancel>
			ctx.App.Stop()
		},
	)
	groupsPlaceholder := tview.NewBox().
		SetTitle(" [ GROUPS PANEL ] ").
		SetBorder(true).
		SetBorderColor(tcell.ColorGreen).
		SetBackgroundColor(tcell.ColorBlack)

	ctx.BodyPages.AddPage("users", ctx.UsersPanel.MainFlex, true, true)
	ctx.BodyPages.AddPage("groups", groupsPlaceholder, true, false)

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

	ctx.DebugView = tview.NewTextView().
		SetTextColor(tcell.ColorRed)
	ctx.DebugView.SetBackgroundColor(tcell.ColorBlack)
	ctx.DebugView.SetText(" [WAITING FOR KEY] ")

	// RETRO FIX: Keep the entire header backdrop flat black
	ctx.HeaderTabs.SetBackgroundColor(tcell.ColorBlack)
	ctx.HeaderTabs.
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 1, 0, false).
		AddItem(ctx.UsersBtn, 9, 0, true).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 2, 0, false).
		AddItem(ctx.GroupsBtn, 10, 0, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 0, 1, false).
		AddItem(ctx.DebugView, 45, 0, false)

	ctx.InnerFlex.AddItem(ctx.HeaderTabs, 1, 1, true)
	ctx.InnerFlex.AddItem(ctx.BodyPages, 0, 1, false)
	ctx.InnerFlex.SetBackgroundColor(tcell.ColorBlack)

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

	// Hand off operational focus to the active tab component
	if ctx.CurrentTab == 0 {
		ctx.App.SetFocus(ctx.UsersBtn)
	} else {
		ctx.App.SetFocus(ctx.GroupsBtn)
	}
}
func (ctx *AppContext) setupKeybindings() {
	ctx.App.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentFocus := ctx.App.GetFocus()

		// 🔍 LIVE DEBUG MONITORING
		focusStr := "NIL"
		if currentFocus != nil {
			focusStr = fmt.Sprintf("%T", currentFocus)
		}
		keyName := tcell.KeyNames[event.Key()]
		ctx.DebugView.SetText(fmt.Sprintf(" [F-PRE: %s | K: %s] ", focusStr, keyName))

		// 1. If the user is viewing the Users tab panel layout
		if ctx.CurrentTab == 0 && ctx.UsersPanel != nil {
			proceedBtn := ctx.UsersPanel.Buttons.GetItem(1)
			cancelBtn := ctx.UsersPanel.Buttons.GetItem(3)
			currIdx := ctx.UsersPanel.List.GetCurrentItem()
			maxIdx := ctx.UsersPanel.List.GetItemCount() - 1
			// Route Down arrow away from the List onto the Proceed button
			if event.Key() == tcell.KeyDown {
				if currentFocus == ctx.UsersPanel.List {
					if currIdx == maxIdx {
						ctx.App.SetFocus(proceedBtn)
						return nil
					} else if currIdx%2 == 0 {
						ctx.UsersPanel.List.SetCurrentItem(currIdx + 2)
						return nil
					}
				}
			}

			// Fix 3: Intercept Up arrow actions on empty spacer lines
			if event.Key() == tcell.KeyUp {
				if currentFocus == ctx.UsersPanel.List {
					if currIdx == 0 {
						return nil
					} else if currIdx%2 == 0 {
						ctx.UsersPanel.List.SetCurrentItem(currIdx - 2)
						return nil
					}
				} else if currentFocus == proceedBtn || currentFocus == cancelBtn {
					ctx.UsersPanel.List.SetCurrentItem(maxIdx)
					ctx.App.SetFocus(ctx.UsersPanel.List)
					return nil
				}
			}

			// Route LEFT and RIGHT arrow keys horizontally across the button strip
			if event.Key() == tcell.KeyRight && currentFocus == proceedBtn {
				ctx.App.SetFocus(cancelBtn)
				return nil
			}
			if event.Key() == tcell.KeyLeft && currentFocus == cancelBtn {
				ctx.App.SetFocus(proceedBtn)
				return nil
			}
		}

		// Tab Navigation Processing
		if event.Key() == tcell.KeyTab {
			ctx.CurrentTab = (ctx.CurrentTab + 1) % 2
			if ctx.CurrentTab == 0 {
				ctx.BodyPages.SwitchToPage("users")

				if ctx.UsersPanel != nil && ctx.UsersPanel.List != nil {
					ctx.UsersPanel.List.SetCurrentItem(0)
					ctx.App.SetFocus(ctx.UsersPanel.MainFlex)
				}
			} else {
				ctx.BodyPages.SwitchToPage("groups")
				ctx.App.SetFocus(ctx.GroupsBtn)
			}
			ctx.TabVisuals()
			// Post-switch focus validation: see if tview accepted our focus request
			postFocus := ctx.App.GetFocus()
			postFocusStr := "NIL"
			if postFocus != nil {
				postFocusStr = fmt.Sprintf("%T", postFocus)
			}
			ctx.DebugView.SetText(fmt.Sprintf(" [TAB SWAP -> F-POST: %s] ", postFocusStr))
			return nil
		}

		// CRUCIAL FIX: Let UP/DOWN keys pass through so tview can scroll list items
		return event
	})
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
