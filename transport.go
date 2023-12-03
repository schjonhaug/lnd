package tapprotocol

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/fxamacker/cbor/v2"
)

type Transport struct {
	connection net.Conn
}

func (transport *Transport) reader(r io.Reader, command any, channel chan any) {
	buf := make([]byte, 1024)
	_, err := r.Read(buf[:])

	if err != nil {
		print(err)
		return
	}

	decMode, _ := cbor.DecOptions{ExtraReturnErrors: cbor.ExtraDecErrorUnknownField}.DecMode()

	switch command.(type) {
	case StatusCommand:

		var v StatusData

		if err := decMode.Unmarshal(buf, &v); err != nil {

			var e ErrorData

			if err := decMode.Unmarshal(buf, &e); err != nil {
				channel <- err
			}

			channel <- e

		}

		channel <- v

	case unsealCommand:

		var v unsealData

		if err := decMode.Unmarshal(buf, &v); err != nil {

			var e ErrorData

			if err := decMode.Unmarshal(buf, &e); err != nil {
				channel <- err
			}

			channel <- e

		}

		channel <- v
	case newCommand:

		var v newData

		if err := decMode.Unmarshal(buf, &v); err != nil {

			var e ErrorData

			if err := decMode.Unmarshal(buf, &e); err != nil {
				channel <- err
			}

			channel <- e

		}

		channel <- v
	case certsCommand:

		var v certsData

		if err := decMode.Unmarshal(buf, &v); err != nil {

			var e ErrorData

			if err := decMode.Unmarshal(buf, &e); err != nil {
				channel <- err

			}

			channel <- e

		}

		channel <- v
	case checkCommand:

		var v checkData

		if err := decMode.Unmarshal(buf, &v); err != nil {

			var e ErrorData

			if err := decMode.Unmarshal(buf, &e); err != nil {
				channel <- err

			}

			channel <- e

		}

		channel <- v
	case readCommand:

		var v readData

		if err := decMode.Unmarshal(buf, &v); err != nil {

			var e ErrorData

			if err := decMode.Unmarshal(buf, &e); err != nil {
				channel <- err
			}

			channel <- e

		}

		channel <- v
	case waitCommand:

		var v waitData

		if err := decMode.Unmarshal(buf, &v); err != nil {

			var e ErrorData

			if err := decMode.Unmarshal(buf, &e); err != nil {
				channel <- err
			}

			channel <- e

		}

		channel <- v

	default:

		var v ErrorData

		if err := decMode.Unmarshal(buf, &v); err != nil {
			channel <- err
		}

		channel <- v

	}

}

func (transport *Transport) Connect() {
	connection, err := net.Dial("unix", "/tmp/ecard-pipe")
	if err != nil {
		log.Fatal(err)
	}
	transport.connection = connection
}

func (transport Transport) Disconnect() {
	transport.connection.Close()
}

func (transport Transport) Send(command any, channel chan any) {

	cbor_serialized, err := cbor.Marshal(command)
	if err != nil {
		fmt.Println("error:", err)
	}

	go transport.reader(transport.connection, command, channel)
	_, err = transport.connection.Write(cbor_serialized)

	if err != nil {
		log.Fatal("write error:", err)
	}

	time.Sleep(100 * time.Millisecond)

}
