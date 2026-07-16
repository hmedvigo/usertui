package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TabPanel struct {
	MainFlex    *tview.Flex
	List        *tview.List
	Buttons     *tview.Flex
	ConfirmedId int
	PanelType   string
}

type PanelConfig struct {
	Title   string
	Options []string
}

func NewTabPanel(panelType string, onProceed func(actionIdx int), onCancel func()) *TabPanel {
	config := getPanelConfig(panelType)
	panel := &TabPanel{
		MainFlex:    tview.NewFlex().SetDirection(tview.FlexRow),
		List:        tview.NewList(),
		Buttons:     tview.NewFlex().SetDirection(tview.FlexColumn),
		ConfirmedId: 0,
		PanelType:   panelType,
	}

	// Helper function to render a line mimicking a retro radio button state
	// MOVE THIS HELPER FUNCTION LATER!!!!
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

	// Build spacing rows into the list layout to expand vertical item height
	for i, opt := range config.Options {
		panel.List.AddItem(getOptionLabel(i, opt), "", 0, nil)
		if i < len(config.Options)-1 {
			// Add an unselectable blank row as a layout padding spacer
			panel.List.AddItem("", "", 0, nil)
		}
	}
	// Nothing happens when navigating up and down
	panel.List.SetChangedFunc(nil)
	// Update brackets indicator whenever the user shifts focus with Up/Down arrows
	panel.List.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index%2 == 0 { // Ensure it's a real item row
			panel.ConfirmedId = index / 2
			// Redraw all items to move the (*)
			for i, opt := range config.Options {
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
				for i, opt := range config.Options {
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
	panel.List.SetSelectedBackgroundColor(tcell.ColorLawnGreen)

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
	titleSpacer := tview.NewBox().SetBackgroundColor(tcell.ColorBlack)
	// Combine components into the container block
	titleBox := tview.NewTextView().
		SetTextColor(tcell.ColorLawnGreen).
		SetText(config.Title)
	panel.MainFlex.SetBackgroundColor(tcell.ColorBlack)
	panel.MainFlex.
		AddItem(titleSpacer, 1, 0, false). // Top spacer
		AddItem(titleBox, 1, 0, false).
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 1, 0, false).
		AddItem(panel.List, 7, 0, true).                                           // Expanded window height to hold spacer lines
		AddItem(tview.NewBox().SetBackgroundColor(tcell.ColorBlack), 2, 0, false). // Fix 2: Substantial spacing split before buttons
		AddItem(panel.Buttons, 1, 0, false)

	return panel
}
func getPanelConfig(panelType string) PanelConfig {
	switch panelType {
	case "users":
		return PanelConfig{
			Title: " User Management Options",
			Options: []string{
				"Create New System User Account",
				"Modify / Edit Existing User Account",
				"Delete User Account From System",
				"View / Search Comprehensive User Directory",
				"View / Search  User Directory",
				"View  Comprehensive User Directory",
			},
		}
	case "groups":
		return PanelConfig{
			Title: " Group Administration Actions:",
			Options: []string{
				"Create New User Group",
				"Modify / Edit Existing Group",
				"Delete Group From System",
				"View / Search Group Directory",
				"Manage Group Memberships",
			},
		}
	default:
		return PanelConfig{
			Title: " Administration Actions:",
			Options: []string{
				"Option 1",
				"Option 2",
				"Option 3",
			},
		}
	}
}
