package controllers

import (
	express "github.com/stefpo/expressgo"
	comp "github.com/stefpo/expressgo/examples/demo/components"
	h "github.com/stefpo/html5"
)

func Html5ViewPage(req *express.Request, resp *express.Response, next func(...express.Error)) {
	var x *h.HTMLElement
	var textField *h.HTMLElement
	title := "Document title"

	document := h.HTML(nil,
		h.HEAD(nil,
			h.TITLE(nil, title)),
		h.BODY(nil,
			h.H1(h.Attr{"id": "toto"}, "The initial title").AssignTo(&x),
			comp.FieldBox("Field label", h.INPUT(h.Attr{"type": "text", "class": "number integer"}).AssignTo(&textField)),
			h.DIV(nil, []interface{}{h.SPAN(nil, "Text"), " [Raw text] ", h.H2(nil, "more text")})))

	x.SetInnerText("The modified H1")

	textField.SetAttribute("value", "\"Stéphane\"")
	textField.SetAttribute("data-test", "A test value")
	textField.DelAttribute("data-test")
	textField.AddClass("red")
	textField.DelClass("integer")

	resp.End(document.ToHTML(true))
}
