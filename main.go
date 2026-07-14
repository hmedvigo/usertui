package main

import (
	"usertui/ui"
)

/*
	func geteuid() int {
		return syscall.Geteuid()
	}
*/
func main() {
	// Enforce root/sudo execution privilege check before allocating app memory
	/*	if geteuid() != 0 {
			_, err := fmt.Fprintln(os.Stderr, "\033[31m[ERROR] This application requires administrative privileges.\033[0m")
			if err != nil {
				return
			}
			os.Exit(1)
		}
	*/
	appUI := ui.NewApp()
	if err := appUI.Run(); err != nil {
		panic(err)
	}
}
