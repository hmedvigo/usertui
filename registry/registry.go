package registry

import "github.com/rivo/tview"

type PageCreator func() tview.Primitive
type PageRegistry interface {
	RegisterPageCreator(name string, creator PageCreator)
	GoBackToMainMenu()
}
