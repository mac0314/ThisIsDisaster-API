package controllers

import (
	"encoding/xml"
	"fmt"

	"github.com/revel/revel"
)

type TEST struct {
	*revel.Controller
}

type (
	stringMap map[string]interface{}

	ErrorResponse struct {
		ID   string    `xml:"id"`
		Meta stringMap `xml:"meta,omitempty"`
	}
)

func (s stringMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	tokens := []xml.Token{start}
	for k, v := range s {
		v := fmt.Sprintf("%v", v)
		t := xml.StartElement{Name: xml.Name{"", k}}
		tokens = append(tokens, t, xml.CharData(v), xml.EndElement{t.Name})
	}
	tokens = append(tokens, xml.EndElement{start.Name})
	for _, t := range tokens {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}
	err := e.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (c TEST) Index() revel.Result {
	m := make(stringMap, 2)
	m["key1"] = "value"
	m["key2"] = 1024

	return c.RenderXML(m)
}
