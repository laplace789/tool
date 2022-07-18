package model

const (
	Level_Info    = "Information"
	Level_Warning = "Warning"
	Level_System  = "System"
	Level_Error   = "Error"
)

type Message interface {
	SetTime()
	String() string
}
