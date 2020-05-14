package eventmanager

import (
	"CardGameGo/src/managers/eventmanager/events"
	"github.com/veandco/go-sdl2/sdl"
)

type EventManager struct {
	screen           int
	RegisteredClicks []events.ClickEvent
}

func New(screen int) *EventManager {
	return &EventManager{
		screen:           screen,
		RegisteredClicks: make([]events.ClickEvent, 0, 5),
	}
}

func (em *EventManager) RegisterEvent(event events.ClickEvent) {
	em.RegisteredClicks = append(em.RegisteredClicks, event)
}

func (em *EventManager) ProcessClickEvents(mouseEv *sdl.MouseButtonEvent) error {
	for _, e := range em.RegisteredClicks {
		if mouseEv.X >= e.GetX() && mouseEv.X <= (e.GetX() + e.GetWidth()) &&
			mouseEv.Y >= e.GetY() && mouseEv.Y <= (e.GetY() + e.GetHeight()) {
			return e.RunCallback(e)
		}
	}
	return nil
}
