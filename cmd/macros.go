package main

import (
	"macropad/internal/general"
	macros_keypress "macropad/internal/macros/keypress"
	macros_rotary "macropad/internal/macros/rotary"
	macros_rotary_btn "macropad/internal/macros/rotary_btn"
	macros_serial "macropad/internal/macros/serial"
)

func Cmd(cmd general.CmdStruct) {

	switch cmd.Type {
	case "serial":
		macros_serial.Cmd(serialPort, cmd)
	case "keypress":
		macros_keypress.Cmd(serialPort, cmd)
	case "rotary":
		macros_rotary.Cmd(serialPort, cmd)
	case "rotary_btn":
		macros_rotary_btn.Cmd(serialPort, cmd)
	}

}
