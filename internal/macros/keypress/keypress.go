package macros_keypress

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"macropad/internal/general"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

func macroKeypress(msg general.CmdStruct, keyarr []KeypressStruct) bool {
	for _, item := range keyarr {
		if msg.Data == item.KeyNum {
			robotgo.KeyTap(item.Keypress)
		}
	}
	return false
}

func getMacro(filepath string) ([]KeypressStruct, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	data := []KeypressStruct{}

	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getMacroFolder(filepath string) (map[string][]KeypressStruct, error) {
	files, err := general.FilePathWalkDir(filepath)
	if err != nil {
		return nil, err
	}

	out := map[string][]KeypressStruct{}

	for i, file := range files {
		data, err := getMacro(file)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		out[fmt.Sprintf("%v", i)] = data

	}

	return out, nil

}

func Cmd(port serial.Port, msg general.CmdStruct) bool {

	var keyarr []KeypressStruct
	// https://github.com/go-vgo/robotgo/blob/master/docs/keys.md

	macroFolder := fmt.Sprintf("%v/macros/", general.GetAppStorageFolder())
	macroGameFolder := fmt.Sprintf("%v/game/", macroFolder)
	macroCustomFolder := fmt.Sprintf("%v/custom/", macroFolder)

	defaultCombo, err := getMacro(fmt.Sprintf("%v/default.json", macroFolder))
	if err != nil {
		log.Error(err.Error())
		return false
	}

	keyCombos := map[string]map[string][]KeypressStruct{
		"1": { // Default mode
			"0": defaultCombo, // Default submode
		},
	}

	keyCombos["2"], err = getMacroFolder(macroGameFolder)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	keyCombos["3"], err = getMacroFolder(macroCustomFolder)
	if err != nil {
		log.Error(err.Error())
		return false
	}

	keyarr = keyCombos[msg.ModeID][msg.SubModeID]

	macroKeypress(msg, keyarr)
	return false
}
