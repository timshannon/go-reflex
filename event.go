// Copyright 2020 Tim Shannon. All rights reserved.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package reflex

import "reflect"

// Event is a DOM event optionally available as a parameter to EventFuncs
type Event struct {
	Type string
	MouseEvent
}

var eventType = reflect.TypeOf(Event{})

// MouseEvent is the base structure from a DOM mouse event
type MouseEvent struct {
	AltKey    bool `json:"altKey"`
	Button    uint `json:"button"`
	Buttons   uint `json:"buttons"`
	ClientX   uint `json:"clientX"`
	ClientY   uint `json:"clientY"`
	CtrlKey   bool `json:"ctrlKey"`
	MetaKey   bool `json:"metaKey"`
	MovementX uint `json:"movementX"`
	MovementY uint `json:"movementY"`
	ScreenX   uint `json:"screenX"`
	ScreenY   uint `json:"screenY"`
	ShiftKey  bool `json:"shiftKey"`
}
