package lndmobile

import (
	"github.com/schjonhaug/tapcards"
)

var satscard tapcards.Satscard

func TapcardsISOAppletSelectRequest() ([]byte, error) {
	return satscard.ISOAppletSelectRequest()
}

func TapcardsStatus() ([]byte, error) {
	return satscard.StatusRequest()
}
