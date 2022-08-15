package main

import (
	_ "embed"
	"fmt"
	"macropad/internal/general"
	macros_keypress "macropad/internal/macros/keypress"
	serialint "macropad/internal/serial"
	"os"
	"strconv"
	"strings"

	"github.com/getlantern/systray"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

var serialPort serial.Port

//go:embed winres/icon.ico
var iconImage []byte

func readSerial() bool {
	lastMsg := ""
	errCnt := 0

	modes, err := general.GetModes()
	if err != nil {
		return fmt.Errorf("%v", err.Error()) != nil
	}

	for errCnt != 10 {
		data, err := serialint.ReadSerial(serialPort)

		if err != nil {
			log.Infoln("Unable to connect to macropad.")
			serialPort.Close()
			return false
		}
		if data == "The handle is invalid." {
			serialPort.Close()
			return false
		}
		if data != lastMsg && data != "" {
			x := strings.Split(data, ".")
			if len(x) != 5 {
				errCnt += 1
				continue
			}
			inNum := x[0]
			inType := x[1]
			inMode := x[2]
			inSub := x[3]
			inData := x[4]

			cmdId, err := strconv.Atoi(inNum)
			if err != nil {
				log.Error(err)
				log.Infoln("Invalid message, invalid id")
				continue
			}

			var cmdType string
			var cmdMode string
			var cmdSub string

			switch inType {
			case "0":
				cmdType = "serial"
			case "1":
				cmdType = "keypress"
			case "2":
				cmdType = "rotary"
			case "3":
				cmdType = "rotary_btn"
			}

			for _, item := range modes {
				done := false
				if item.ID != inMode {
					continue
				}
				for _, x := range item.Submodes {
					if x.ID == inSub {
						cmdMode = item.Name
						cmdSub = x.Name
						done = true
					}
				}
				if done {
					break
				}
			}

			cmdStruct := general.CmdStruct{
				Id:        cmdId,
				Type:      cmdType,
				Mode:      cmdMode,
				ModeID:    inMode,
				SubMode:   cmdSub,
				SubModeID: inSub,
				Data:      inData,
			}

			Cmd(cmdStruct)

			lastMsg = data
			errCnt = 0
		}
	}

	return false
}

func bootstrap() serial.Port {
	for {
		serialPort, err := serialint.OpenSerial()
		if err != nil {
			continue
		}
		log.Infoln("Connected to macropad!")
		return serialPort
	}
}

func onQuit(mQuit *systray.MenuItem) {
	for {
		select {
		case <-mQuit.ClickedCh:
			systray.Quit()
			return
		}
	}
}

func onReady() {
	systray.SetIcon(iconImage)
	systray.SetTitle("Macropad")
	systray.SetTooltip("Macropad v0.0.1")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	serialPort = bootstrap()
	go onQuit(mQuit)
	for {
		readSerial()
		serialPort.Close()
		serialPort = nil
		serialPort = bootstrap()
	}

}

func onExit() {
	// Cleaning stuff here.
}

func setup() error {
	basePath := general.GetAppStorageFolder()
	macrosPath := fmt.Sprintf("%v/macros/", basePath)
	macrosGamePath := fmt.Sprintf("%v/game/", macrosPath)
	macrosCustomPath := fmt.Sprintf("%v/custom/", macrosPath)

	settingsPath := fmt.Sprintf("%v/settings.json", basePath)
	defaultKeybindsPath := fmt.Sprintf("%v/default.json", macrosPath)

	blankKeybinds := []macros_keypress.KeypressStruct{
		{KeyNum: "0", Keypress: ""}, {KeyNum: "1", Keypress: ""}, {KeyNum: "2", Keypress: ""},
		{KeyNum: "3", Keypress: ""}, {KeyNum: "4", Keypress: ""}, {KeyNum: "5", Keypress: ""},
		{KeyNum: "6", Keypress: ""}, {KeyNum: "7", Keypress: ""}, {KeyNum: "8", Keypress: ""},
		{KeyNum: "9", Keypress: ""}, {KeyNum: "10", Keypress: ""}, {KeyNum: "11", Keypress: ""},
	}

	// Settings
	if !general.PathExists(settingsPath) {
		os.MkdirAll(basePath, os.ModePerm)
		defSettings := general.SettingsStruct{
			Color: []int{
				160, 160, 160,
				180, 180, 180,
				200, 200, 202,
				220, 220, 220,
			},
			Tone: []int{
				196, 220, 246,
				262, 294, 330,
				349, 392, 440,
				494, 523, 587,
			},
		}
		err := general.WriteStructFile(settingsPath, defSettings)
		if err != nil {
			return err
		}
	}

	// Default keybinds
	if !general.PathExists(defaultKeybindsPath) {
		os.MkdirAll(macrosPath, os.ModePerm)
		keyBinds := []macros_keypress.KeypressStruct{
			{KeyNum: "0", Keypress: "num-"}, {KeyNum: "1", Keypress: "num+"}, {KeyNum: "2", Keypress: "backspace"},
			{KeyNum: "3", Keypress: "num7"}, {KeyNum: "4", Keypress: "num8"}, {KeyNum: "5", Keypress: "num9"},
			{KeyNum: "6", Keypress: "num4"}, {KeyNum: "7", Keypress: "num5"}, {KeyNum: "8", Keypress: "num6"},
			{KeyNum: "9", Keypress: "num1"}, {KeyNum: "10", Keypress: "num2"}, {KeyNum: "11", Keypress: "num3"},
		}
		defPath := fmt.Sprintf("%v/default.json", macrosPath)
		err := general.WriteStructFile(defPath, keyBinds)
		if err != nil {
			return err
		}
	}

	// Game macros

	if !general.PathExists(macrosGamePath) {
		os.MkdirAll(macrosGamePath, os.ModePerm)
	}
	empt, err := general.CheckFolderEmpty(macrosGamePath)
	if err != nil {
		return err
	}
	if empt {
		defPath := fmt.Sprintf("%v/default.json", macrosGamePath)
		err := general.WriteStructFile(defPath, blankKeybinds)
		if err != nil {
			return err
		}
	}

	// Custom macros

	if !general.PathExists(macrosCustomPath) {
		os.MkdirAll(macrosCustomPath, os.ModePerm)
	}

	empt, err = general.CheckFolderEmpty(macrosCustomPath)
	if err != nil {
		return err
	}
	if empt {
		defPath := fmt.Sprintf("%v/default.json", macrosCustomPath)
		err := general.WriteStructFile(defPath, blankKeybinds)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	general.LoggingInit()
	setup()

	log.Infoln("Waiting for macropad connection...")
	systray.Run(onReady, onExit)
}
