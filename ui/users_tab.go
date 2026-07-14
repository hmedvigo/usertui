package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type UsersTabPanel struct {
	MainFlex    *tview.Flex
	List        *tview.List
	Buttons     *tview.Flex
	ConfirmedId int
}

func NewUsersTabPanel(onProceed func(actionIdx int), onCancel func()) *UsersTabPanel {
	panel := &UsersTabPanel{
		MainFlex:    tview.NewFlex().SetDirection(tview.FlexRow),
		List:        tview.NewList(),
		Buttons:     tview.NewFlex().SetDirection(tview.FlexColumn),
		ConfirmedId: 0,
	}

	options := []string{
		"Create New System User Account",
		"Modify / Edit Existing User Account",
		"Delete User Account From System",
		"View / Search Comprehensive User Directory",
	}

	// Helper function to render a line mimicking a retro radio button state
	getOptionLabel := func(idx int, text string) string {
		prefix := "( )"
		if idx == panel.ConfirmedId {
			prefix = "(*)"
		}
		// Fix 2: Add left indentation padding ("  ") so it isn't flush with the edge
		return fmt.Sprintf("  %s %s", prefix, text)
	}
	// Turn off secondary text so each option takes exactly 1 row of screen height
	panel.List.ShowSecondaryText(false)
	panel.List.SetCurrentItem(0)

	// Fix 2: Build spacing rows into the list layout to expand vertical item height
	for i, opt := range options {
		panel.List.AddItem(getOptionLabel(i, opt), "", 0, nil)
		if i < len(options)-1 {
			// Add an unselectable blank row as a layout padding spacer
			panel.List.AddItem("", "", 0, nil)
		}
	}

	panel.List.SetChangedFunc(nil)
	// Update brackets indicator whenever the user shifts focus with Up/Down arrows
	panel.List.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index%2 == 0 { // Ensure it's a real item row
			panel.ConfirmedId = index / 2
			// Redraw all items to move the (*)
			for i, opt := range options {
				panel.List.SetItemText(i*2, getOptionLabel(i, opt), "")
			}
		}
	})

	// Capture spacebar cleanly inside the list box
	panel.List.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == ' ' {
			curr := panel.List.GetCurrentItem()
			if curr%2 == 0 {
				panel.ConfirmedId = curr / 2
				for i, opt := range options {
					panel.List.SetItemText(i*2, getOptionLabel(i, opt), "")
				}
			}
			return nil
		}
		return event
	})

	// Setup retro button controls at the bottom
	proceedBtn := tview.NewButton("< Proceed to Action >")
	cancelBtn := tview.NewButton("< Cancel >")

	proceedBtn.SetSelectedFunc(func() {
		if onProceed != nil {
			onProceed(panel.ConfirmedId)
		}
	})
	cancelBtn.SetSelectedFunc(func() {
		if onCancel != nil {
			onCancel()
		}
	})

	// Configure Retro Color Themes
	panel.List.SetBackgroundColor(tcell.ColorBlack)
	panel.List.SetMainTextColor(tcell.ColorDarkGreen)

	// Active selection color config (Turning from dim grey/green to vibrant phosphor green background)
	panel.List.SetSelectedTextColor(tcell.ColorBlack)
	panel.List.SetSelectedBackgroundColor(tcell.ColorGreen)

	// Style Bottom Buttons
	btnActiveStyle := tcell.StyleDefault.Background(tcell.ColorGreen).Foreground(tcell.ColorBlack)
	btnInactiveStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorDarkGreen)

	proceedBtn.SetStyle(btnInactiveStyle)
	proceedBtn.SetLabelColorActivated(tcell.ColorGreen)
	proceedBtn.SetBackgroundColorActivated(tcell.ColorBlack)

	cancelBtn.SetStyle(btnInactiveStyle)
	cancelBtn.SetLabelColorActivated(tcell.ColorGreen)
	cancelBtn.SetBackgroundColorActivated(tcell.ColorBlack)

	// Tie button focus shifts to visual highlight swaps
	proceedBtn.SetFocusFunc(func() { proceedBtn.SetStyle(btnActiveStyle) })
	proceedBtn.SetBlurFunc(func() { proceedBtn.SetStyle(btnInactiveStyle) })
	cancelBtn.SetFocusFunc(func() { cancelBtn.SetStyle(btnActiveStyle) })
	cancelBtn.SetBlurFunc(func() { cancelBtn.SetStyle(btnInactiveStyle) })

	// Lay out the bottom button strip horizontally
	panel.Buttons.SetBackgroundColor(tcell.ColorBlack)
	panel.Buttons.
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 0, 1, false). // Flexible left margin spacer
		AddItem(proceedBtn, 23, 0, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 4, 0, false). // Gap between buttons
		AddItem(cancelBtn, 12, 0, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 0, 1, false) // Flexible right margin spacer
	// Combine components into the container block
	titleBox := tview.NewTextView().
		SetTextColor(tcell.ColorDarkGreen).
		SetText(" User Administration Actions:").
		SetBackgroundColor(tcell.ColorBlack)

	panel.MainFlex.SetBackgroundColor(tcell.ColorBlack)
	panel.MainFlex.
		AddItem(titleBox, 1, 0, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 1, 0, false).
		AddItem(panel.List, 7, 0, true).                                           // Expanded window height to hold spacer lines
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 2, 0, false). // Fix 2: Substantial spacing split before buttons
		AddItem(panel.Buttons, 1, 0, false)

	return panel
}
