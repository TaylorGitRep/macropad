package general

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func GetModes() ([]ModeStruct, error) {

	modes := []ModeStruct{
		{
			ID:   "0",
			Name: "Bootstrap",
			Submodes: []SubModeStruct{
				{ID: "0", Name: "Default"},
			},
		},
		{
			ID:   "1",
			Name: "Default",
			Submodes: []SubModeStruct{
				{ID: "0", Name: "Default"},
			},
		},
	}

	gameMacros := fmt.Sprintf("%v/macros/game", GetAppStorageFolder())
	customMacros := fmt.Sprintf("%v/macros/custom", GetAppStorageFolder())

	os.MkdirAll(gameMacros, os.ModePerm)
	os.MkdirAll(customMacros, os.ModePerm)

	modeGame := ModeStruct{
		ID:   "2",
		Name: "Game",
	}

	files, err := FilePathWalkDir(gameMacros)
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		modeGame.Submodes = append(modeGame.Submodes, SubModeStruct{ID: fmt.Sprintf("%v", i), Name: strings.ReplaceAll(filepath.Base(file), filepath.Ext(file), "")})
	}

	modes = append(modes, modeGame)

	modeCustom := ModeStruct{
		ID:   "3",
		Name: "Custom",
	}

	files, err = FilePathWalkDir(customMacros)
	if err != nil {
		return nil, err
	}

	for i, file := range files {
		modeCustom.Submodes = append(modeGame.Submodes, SubModeStruct{ID: fmt.Sprintf("%v", i), Name: strings.ReplaceAll(filepath.Base(file), filepath.Ext(file), "")})
	}

	modes = append(modes, modeCustom)

	return modes, nil
}

func PathExists(checkPath string) bool {
	if _, err := os.Stat(checkPath); !os.IsNotExist(err) {
		return true
	}
	return false
}

func GetAppStorageFolder() string {
	var appstorage string
	if runtime.GOOS == "windows" {
		appstorage = fmt.Sprintf("%v/TaylorGitRep/macropad/", os.Getenv("localappdata"))
	} else {
		dirname, _ := os.UserHomeDir()
		dirPath := "/.TaylorGitRep/macropad"
		if runtime.GOOS == "darwin" {
			dirPath = "/Library/Application Support/TaylorGitRep/macropad/"
		}
		appstorage = fmt.Sprintf("%v/%v", dirname, dirPath)
	}

	return appstorage
}

func WriteStructFile(filePath string, data interface{}) error {
	datastr, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filePath, []byte(datastr), 0644)
	if err != nil {
		return err
	}
	return nil
}

func CheckFolderEmpty(folderPath string) (bool, error) {
	files, err := FilePathWalkDir(folderPath)
	if err != nil {
		return false, err
	}
	if len(files) == 0 {
		return true, nil
	}
	return false, nil
}

func GetSettings() (SettingsStruct, error) {
	ret := SettingsStruct{}
	settingsPath := fmt.Sprintf("%v/settings.json", GetAppStorageFolder())
	data, err := ioutil.ReadFile(settingsPath)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return ret, err
	}
	return ret, nil
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
