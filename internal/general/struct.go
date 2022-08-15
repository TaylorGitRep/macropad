package general

type CmdStruct struct {
	Id        int
	Type      string
	Mode      string
	ModeID    string
	SubMode   string
	SubModeID string
	Data      string
}

type SubModeStruct struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ModeStruct struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Submodes []SubModeStruct `json:"submodes"`
}

type SettingsStruct struct {
	Color []int `json:"color"`
	Tone  []int `json:"tone"`
}
