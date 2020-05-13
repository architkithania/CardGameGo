package buttons

import (
	"CardGameGo/src/managers/eventmanager/events"
)

type Button interface {
	events.ClickEvent
}
