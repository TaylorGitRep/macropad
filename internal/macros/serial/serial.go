package macros_serial

import (
	"errors"
	"fmt"
	"macropad/internal/general"
	serialint "macropad/internal/serial"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

func arrIntToStr(inp []int) string {
	var data string
	for _, x := range inp {
		data += fmt.Sprintf(".%v", x)
	}
	return data
}

func bootstrap(port serial.Port) error {
	settings, err := general.GetSettings()
	if err != nil {
		return err
	}
	err = serialint.SendSerial(port, "bootstrap")
	if err != nil {
		return err
	}
	err = serialint.SendSerial(port, "mode.1.0") // Set mode to default submode to default
	if err != nil {
		return err
	}
	err = serialint.SendSerial(port, fmt.Sprintf("color%v", arrIntToStr(settings.Color))) // Set mode to default submode to default
	if err != nil {
		return err
	}
	err = serialint.SendSerial(port, fmt.Sprintf("tone%v", arrIntToStr(settings.Tone))) // Set mode to default submode to default
	if err != nil {
		return err
	}
	err = serialint.SendSerial(port, "brightness.1-0") // Set mode to default submode to default
	return err
}

func heartbeat(port serial.Port) error {
	return serialint.SendSerial(port, "heartbeat")
}

func Cmd(port serial.Port, message general.CmdStruct) error {
	log.Debugln(message)
	switch message.Data {
	case "bootstrap":
		return bootstrap(port)
	case "heartbeat":
		return heartbeat(port)
	default:
		return errors.New("command not found")
	}
}
