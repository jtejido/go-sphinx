package util

import (
	"encoding/xml"
	"io"
)

type Handler interface {
	//called when XML document start
	StartDocument()
	//called when XML document end
	EndDocument()
	//called when XML tag start
	StartElement(xml.StartElement)
	//called when XML tag end
	EndElement(xml.EndElement)
	//called when the parser encount chardata
	CharData(xml.CharData)
	//called when the parser encount comment
	Comment(xml.Comment)
	//called when the parser encount procInst
	//<!procinst >
	ProcInst(xml.ProcInst)
	//called when the parser encount directive
	//
	Directive(xml.Directive)
	//called when the parser hit an error
	Error(error)
}

//BaseHandler is a implemented Handler
//All methods do nothing
//You need not implement all Handler methods.
/*
type PartialHandler struct{BaseHandler}
func (h PartialHandler) StartElement(xml.StartElement){
  //do something
}
*/
type BaseHandler struct{}

func (h BaseHandler) StartDocument()                {}
func (h BaseHandler) EndDocument()                  {}
func (h BaseHandler) StartElement(xml.StartElement) {}
func (h BaseHandler) EndElement(xml.EndElement)     {}
func (h BaseHandler) CharData(xml.CharData)         {}
func (h BaseHandler) Comment(xml.Comment)           {}
func (h BaseHandler) ProcInst(xml.ProcInst)         {}
func (h BaseHandler) Directive(xml.Directive)       {}
func (h BaseHandler) Error(error)                   {}

// SAX-like XML Parser
type Parser struct {
	*xml.Decoder
	handler Handler
}

// Create a New Parser
func NewParser(reader io.Reader, handler Handler) *Parser {
	decoder := xml.NewDecoder(reader)
	return &Parser{decoder, handler}
}

// SetHTMLMode make Parser can parse invalid HTML
func (p *Parser) SetHTMLMode() {
	p.Strict = false
	p.AutoClose = xml.HTMLAutoClose
	p.Entity = xml.HTMLEntity
}

// Parse calls handler's methods
// when the parser encount a start-element,a end-element, a comment and so on.
func (p *Parser) Parse() (err error) {
	p.handler.StartDocument()

	for {
		token, err := p.Token()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			p.handler.Error(err)
		}

		switch token.(type) {
		case xml.StartElement:
			s := token.(xml.StartElement)
			p.handler.StartElement(s)
		case xml.EndElement:
			e := token.(xml.EndElement)
			p.handler.EndElement(e)
		case xml.CharData:
			c := token.(xml.CharData)
			p.handler.CharData(c)
		case xml.Comment:
			com := token.(xml.Comment)
			p.handler.Comment(com)
		case xml.ProcInst:
			pro := token.(xml.ProcInst)
			p.handler.ProcInst(pro)
		case xml.Directive:
			dir := token.(xml.Directive)
			p.handler.Directive(dir)
		default:
			panic("unknown xml token.")
		}
	}

	p.handler.EndDocument()
	return
}
