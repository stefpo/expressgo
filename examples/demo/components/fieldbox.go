package components

import (
	h "github.com/stefpo/html5"
)

func FieldBox(title string, field *h.HTMLElement) *h.HTMLElement {
	box := h.DIV(nil,
		h.DIV(nil, title),
		field)
	return box
}
