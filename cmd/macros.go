package main

import (
	"github.com/TaylorGitRep/macropad/internal/general"
	macros_keypress "github.com/TaylorGitRep/macropad/internal/macros/keypress"
	macros_rotary "github.com/TaylorGitRep/macropad/internal/macros/rotary"
	macros_rotary_btn "github.com/TaylorGitRep/macropad/internal/macros/rotary_btn"
	macros_serial "github.com/TaylorGitRep/macropad/internal/macros/serial"
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
