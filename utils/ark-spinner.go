package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type Spinner struct {
	spinner *spinner.Spinner
}

func NewArkSpinner() *Spinner {
	s := spinner.New(spinner.CharSets[18], 100*time.Millisecond)
	return &Spinner{
		spinner: s,
	}
}

func (a *Spinner) Start(message string) {
	a.spinner.Color("blue")
	a.spinner.Start()
	a.spinner.Suffix = " " + message
}
func (a *Spinner) Stop(err error, message string) {

	if err == nil {
		a.spinner.FinalMSG = "✅  " + message + "\n"
	} else {
		a.spinner.FinalMSG = "❌  " + message + "\n"
	}

	a.spinner.Stop()
}

func (a *Spinner) InfoStatusEvent(message string) {
	fmt.Println("ℹ️  " + message)
}

func (a *Spinner) SuccessStatusEvent(message string) {
	fmt.Println("✅  " + message)
}

func (a *Spinner) ErrorStatusEvent(message string) {
	fmt.Println("❌  " + message)
}
