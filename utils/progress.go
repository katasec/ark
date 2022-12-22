package utils

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

type ArkSpinner struct {
	spinner *spinner.Spinner
}

func NewArkSpinner() *ArkSpinner {
	s := spinner.New(spinner.CharSets[18], 100*time.Millisecond)
	return &ArkSpinner{
		spinner: s,
	}
}

func (a *ArkSpinner) Start(message string) {
	a.spinner.Color("blue")
	a.spinner.Start()
	a.spinner.Suffix = " " + message
}
func (a *ArkSpinner) Stop(err error, message string) {

	if err == nil {
		a.spinner.FinalMSG = "✅  " + message + "\n"
	} else {
		a.spinner.FinalMSG = "❌  " + message + "\n"
	}

	a.spinner.Stop()
}

func InfoMessage(message string) {
	fmt.Println("ℹ️ " + message)
}
