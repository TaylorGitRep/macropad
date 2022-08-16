package macros_rotary

import (
	"fmt"
	"strconv"

	"github.com/TaylorGitRep/macropad/internal/general"
	serialint "github.com/TaylorGitRep/macropad/internal/serial"

	log "github.com/sirupsen/logrus"

	"go.bug.st/serial"
)

func setmode(port serial.Port, mode string, submode string) bool {
	err := serialint.SendSerial(port, fmt.Sprintf("mode.%v.%v", mode, submode))
	return err == nil
}

func Cmd(port serial.Port, message general.CmdStruct) bool {

	if message.Data != "1" { // Ignore up press
		return true
	}

	modes, err := general.GetModes()
	if err != nil {
		return false
	}

	curMode, err := strconv.Atoi(message.ModeID)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	nextMode := "1"

	for _, item := range modes {
		itemid, err := strconv.Atoi(item.ID)
		if err != nil {
			log.Error(err.Error())
			continue
		}
		if itemid > curMode {
			nextMode = item.ID
			break
		}
	}

	return setmode(port, nextMode, "0")
}
