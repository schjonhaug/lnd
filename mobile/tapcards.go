package lndmobile

import (
	"github.com/schjonhaug/tapcards"
)

var satscard tapcards.Satscard

func SatscardStatus() ([]byte, error) {
	return satscard.StatusRequest()
}
