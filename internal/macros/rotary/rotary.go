package macros_rotary

import (
	"fmt"
	"strconv"

	"github.com/TaylorGitRep/macropad/internal/general"
	serialint "github.com/TaylorGitRep/macropad/internal/serial"

	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var lastRot int

func setmode(port serial.Port, mode string, submode string) bool {
	err := serialint.SendSerial(port, fmt.Sprintf("mode.%v.%v", mode, submode))
	return err == nil
}

func Cmd(port serial.Port, message general.CmdStruct) bool {

	modes, err := general.GetModes()
	if err != nil {
		return false
	}

	curSub, err := strconv.Atoi(message.SubModeID)
	if err != nil {
		log.Error(err)
		return false
	}
	curRot, err := strconv.Atoi(message.Data)
	if err != nil {
		log.Error(err)
		return false
	}
	nextSub := "0"
	reverse := false

	if lastRot > curRot {
		reverse = true
	}
	lastRot = curRot

	for _, item := range modes {
		if item.ID != message.ModeID {
			continue // Skip if it's the wrong mode
		}
		for _, data := range item.Submodes {
			check := fmt.Sprintf("%v", curSub+1)
			if reverse {
				check = fmt.Sprintf("%v", curSub-1)
			}
			if data.ID == check {
				nextSub = data.ID
				break
			}
		}
	}

	return setmode(port, message.ModeID, nextSub)
}
