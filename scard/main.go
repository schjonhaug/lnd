package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/ebfe/scard"
	"github.com/fxamacker/cbor/v2"
	"github.com/skythen/apdu"

	tapprotocol "github.com/schjonhaug/coinkite-tap-proto-go"
)

func die(err error) {
	fmt.Println(err)
	os.Exit(1)
}

func waitUntilCardPresent(ctx *scard.Context, readers []string) (int, error) {
	rs := make([]scard.ReaderState, len(readers))
	for i := range rs {
		rs[i].Reader = readers[i]
		rs[i].CurrentState = scard.StateUnaware
	}

	for {
		for i := range rs {
			if rs[i].EventState&scard.StatePresent != 0 {
				return i, nil
			}
			rs[i].CurrentState = rs[i].EventState
		}
		err := ctx.GetStatusChange(rs, -1)
		if err != nil {
			return -1, err
		}
	}
}

func main() {

	var tapProtocol tapprotocol.TapProtocol

	// Establish a context
	ctx, err := scard.EstablishContext()
	if err != nil {
		die(err)
	}
	defer ctx.Release()

	// List available readers
	readers, err := ctx.ListReaders()
	if err != nil {
		die(err)
	}

	fmt.Printf("Found %d readers:\n", len(readers))
	for i, reader := range readers {
		fmt.Printf("[%d] %s\n", i, reader)
	}

	if len(readers) > 0 {

		fmt.Println("Waiting for a Card")
		index, err := waitUntilCardPresent(ctx, readers)
		if err != nil {
			die(err)
		}

		// Connect to card
		fmt.Println("Connecting to card in ", readers[index])
		card, err := ctx.Connect(readers[index], scard.ShareExclusive, scard.ProtocolAny)
		if err != nil {
			die(err)
		}
		defer card.Disconnect(scard.ResetCard)

		fmt.Println("Card status:")
		status, err := card.Status()
		if err != nil {
			die(err)
		}

		fmt.Printf("\treader: %s\n\tstate: %x\n\tactive protocol: %x\n\tatr: % x\n",
			status.Reader, status.State, status.ActiveProtocol, status.Atr)

		cmd := cmd()

		fmt.Println("Transmit:")
		fmt.Printf("\tc-apdu: % x\n", cmd)
		rsp, err := card.Transmit(cmd)
		if err != nil {
			die(err)
		}
		fmt.Printf("\tr-apdu: % x\n", rsp)

		rapdu, err := apdu.ParseRapdu(rsp)

		if err != nil {
			die(err)

		}

		fmt.Println("SW1:", hex.EncodeToString([]byte{rapdu.SW1}))
		fmt.Println("SW2:", hex.EncodeToString([]byte{rapdu.SW2}))

		decMode, _ := cbor.DecOptions{ExtraReturnErrors: cbor.ExtraDecErrorUnknownField}.DecMode()

		var v tapprotocol.StatusData

		if err := decMode.Unmarshal(rapdu.Data, &v); err != nil {

			var e tapprotocol.ErrorData

			if err := decMode.Unmarshal(rapdu.Data, &e); err != nil {
				fmt.Println("error:", err)
				//channel <- err
			}

			fmt.Println(e)

			//channel <- e

		}

		fmt.Println(v)

		//	channel <- v

		// API

		cmd, err = tapProtocol.StatusRequest()

		if err != nil {
			die(err)

		}

		fmt.Println("Transmit:")
		fmt.Printf("\tc-apdu: % x\n", cmd)
		rsp, err = card.Transmit(cmd)
		if err != nil {
			die(err)
		}
		fmt.Printf("\tr-apdu: % x\n", rsp)

		rapdu, err = apdu.ParseRapdu(rsp)

		if err != nil {
			die(err)

		}

		fmt.Println("SW1:", hex.EncodeToString([]byte{rapdu.SW1}))

		var w tapprotocol.StatusData

		if err := decMode.Unmarshal(rapdu.Data, &w); err != nil {

			var e tapprotocol.ErrorData

			if err := decMode.Unmarshal(rapdu.Data, &e); err != nil {
				fmt.Println("error:", err)
				//channel <- err
			}

			fmt.Println(e)

			//channel <- e

		}

		fmt.Println(w)

	}
}

func cmd() []byte {

	//var cmd = []byte{0x00, 0xa4, 0x04, 0x00, 0x0f}
	data := []byte{0xf0, 'C', 'o', 'i', 'n', 'k', 'i', 't', 'e', 'C', 'A', 'R', 'D', 'v', '1'}

	//cmd = append(cmd, data[:]...)
	capdu := apdu.Capdu{Cla: 0x00, Ins: 0xa4, P1: 0x04, Data: data}

	bytes, err := capdu.Bytes()

	if err != nil {

		fmt.Println("error:", err)
	}
	fmt.Println("bytes:", bytes)

	return bytes
}
