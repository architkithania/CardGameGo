package events

type ClickEvent interface {
	GetX() int32
	GetY() int32
	GetWidth() int32
	GetHeight() int32
	RunCallback(...interface{}) error
}
