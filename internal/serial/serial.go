package serial

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

func getAllSerial() []string {
	ports, err := serial.GetPortsList()
	if err != nil {
		return nil
	}
	if len(ports) == 0 {
		return nil
	}
	return ports
}

func ReadSerial(port serial.Port) (string, error) {
	buff := make([]byte, 100)
	var sb strings.Builder
	for {
		n, err := port.Read(buff)
		if err != nil {
			return "", err
		}
		if n == 0 {
			break
		}

		cap := false

		for _, item := range buff {
			if string(item) == "\n" {
				cap = true
				break
			}
			if item != 0 {
				sb.WriteByte(item)
			}
		}

		if cap { // If we have captured the data, break
			break
		}

	}

	port.ResetOutputBuffer()

	return sb.String(), nil

}

func SendSerial(port serial.Port, msg string) error {
	data := msg + "\n"
	log.Debugln("Sending message over serial: %v", data)
	_, err := port.Write([]byte(data))
	if err != nil {
		port.Close()
		return err
	}
	time.Sleep(150 * time.Millisecond)
	return nil

}

func OpenSerial() (serial.Port, error) {

	allPorts := getAllSerial()

	var port serial.Port
	var err error

	for _, serialport := range allPorts {
		mode := &serial.Mode{
			BaudRate: 115200,
			DataBits: 8,
		}
		port, err = serial.Open(serialport, mode)

		if err != nil {
			continue // Skip it if we can't open it
		}

		port.SetReadTimeout(100 * time.Millisecond)
		port.ResetInputBuffer()
		port.ResetOutputBuffer()

		for i := 0; i <= 2; i++ {
			time.Sleep(500 * time.Millisecond)
			outp, err := ReadSerial(port)
			if err != nil {
				port.Close()
				return nil, err
			}
			if strings.Contains(outp, ".") {
				data := strings.Split(outp, ".")
				if len(data) == 5 {
					if data[4] == "heartbeat" {
						return port, nil
					}
				}
			}
			port.ResetOutputBuffer()
		}
	}

	return nil, fmt.Errorf("device not found")

}
