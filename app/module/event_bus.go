package module

import (
	"github.com/asaskevich/EventBus"
)

// https://github.com/asaskevich/EventBus

// NewEventBus 全局缓存
func NewEventBus() (EventBus.Bus, func(), error) {
	bus := EventBus.New()
	return bus, func() {}, nil
}
