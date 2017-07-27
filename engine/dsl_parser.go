package engine

import (
	"encoding/json"
	"fmt"

	"github.com/venicegeo/belltower/mpg/mlog"
)

//---------------------------------------------------------------------

type Parser struct {
	tokenizer *Tokenizer

	graph       *GraphModel
	components  []*ComponentModel
	connections []*ConnectionModel
}

func ParseDSL(lines string) (*GraphModel, error) {
	tokenizer := &Tokenizer{}

	err := tokenizer.Scan(lines)
	if err != nil {
		return nil, err
	}

	parser := &Parser{
		tokenizer: tokenizer,
	}

	err = parser.parse()
	if err != nil {
		return nil, err
	}

	return parser.graph, nil
}

func (p *Parser) parse() error {
	var err error

	done := false

	for {
		err = nil

		token := p.tokenizer.Pop()

		switch {
		case token.isEOL():

			p.skipEOLs()
		case token.isIdentN("graph"):
			err = p.toGraphState()
			done = true
		default:
			err = p.toErrorState(token, "expected EOF, EOL, or graph")
		}

		if err != nil {
			return err
		}

		if done {
			break
		}
	}

	err = p.matchEOF()
	if err != nil {
		return err
	}

	if p.graph == nil {
		return fmt.Errorf("did not find a graph!")
	}

	return nil
}

// Skip over any EOLs. If there aren't any, that's okay
func (p *Parser) skipEOLs() {
	for {
		t := p.tokenizer.Peek()
		if !t.isEOL() {
			return
		}
		_ = p.tokenizer.Pop()
	}
}

// Skip over any EOLs, but there must be at least one
func (p *Parser) matchEOLs() error {
	t := p.tokenizer.Pop()
	if !t.isEOL() {
		return p.toErrorState(t, "expected EOL")
	}

	p.skipEOLs()

	return nil
}

// Skip over any EOLs, then match EOF
func (p *Parser) matchEOF() error {
	p.skipEOLs()

	t := p.tokenizer.Pop()
	if !t.isEOF() {
		return p.toErrorState(t, "expected EOF")
	}
	return nil
}

func (p *Parser) toErrorState(tok Token, msg string) error {
	return fmt.Errorf("ERROR: %s\n  got: %s\n", msg, tok)
}

func (p *Parser) toGraphState() error {
	mlog.Debugf("at actionGraph")
	var err error

	p.graph = &GraphModel{}

	token := p.tokenizer.Pop()
	if !token.isIdent() {
		err = p.toErrorState(token, "expected identifier (6)")
		return err
	}
	p.graph.Name = token.str

	err = p.matchEOLs()
	if err != nil {
		return err
	}

	done := false

	for {
		err = nil

		token = p.tokenizer.Pop()

		switch {
		case token.isIdentN("component"):
			err = p.toComponentState()
		case token.isIdentN("metadata"):
			err = p.toMetadataState()
		case token.isIdentN("end"):
			err = p.matchEOLs()
			done = true
		case token.isIdent():
			p.tokenizer.PutBack(token)
			mlog.Debugf("PUSHING %s", token)
			err = p.toConnectionState()
		default:
			err = p.toErrorState(token, "expected 'component', 'metadata', 'end', or identifier")
		}

		if err != nil {
			return err
		}

		if done {
			break
		}
	}

	return nil
}

func (p *Parser) toComponentState() error {
	mlog.Debugf("STATE: component")
	var err error

	component := &ComponentModel{
		Config: ArgMap{},
	}

	token := p.tokenizer.Pop()
	if !token.isIdent() {
		err = p.toErrorState(token, "expected identifier (1)")
		return err
	}
	component.Name = token.str

	err = p.matchEOLs()
	if err != nil {
		return err
	}

	done := false

	for {
		err = nil

		token = p.tokenizer.Pop()

		switch {
		case token.isIdentN("type"):
			err = p.toFieldValueState(&component.Type)
		case token.isIdentN("precondition"):
			err = p.toFieldValueState(&component.Precondition)
		case token.isIdentN("postcondition"):
			err = p.toFieldValueState(&component.Postcondition)
		case token.isIdentN("config"):
			err = p.toConfigState(&component.Config)
		case token.isIdentN("end"):
			err = p.matchEOLs()
			done = true
		default:
			err = p.toErrorState(token, "expected 'type', 'name', 'precondition', etc.")
		}

		if err != nil {
			return err
		}

		if done {
			break
		}
	}

	p.graph.Components = append(p.graph.Components, component)

	return nil
}

func (p *Parser) toConnectionState() error {
	mlog.Debugf("STATE: connection %s", p.tokenizer.Peek())
	var err error

	// name.name -> name.name

	token := p.tokenizer.Pop()
	if !token.isIdent() {
		return p.toErrorState(token, "expected identifier (2)")
	}
	sourceComponent := token.str

	token = p.tokenizer.Pop()
	if token.typ != PeriodMarker {
		return p.toErrorState(token, "expected '.'")
	}

	token = p.tokenizer.Pop()
	if !token.isIdent() {
		return p.toErrorState(token, "expected identifier (3)")
	}
	sourcePort := token.str

	token = p.tokenizer.Pop()
	if token.typ != HyphenMarker {
		return p.toErrorState(token, "expected '-'")
	}

	token = p.tokenizer.Pop()
	if token.typ != GreaterMarker {
		return p.toErrorState(token, "expected '>'")
	}

	token = p.tokenizer.Pop()
	if !token.isIdent() {
		return p.toErrorState(token, "expected identifier (4)")
	}
	destinationComponent := token.str

	token = p.tokenizer.Pop()
	if token.typ != PeriodMarker {
		return p.toErrorState(token, "expected '.'")
	}

	token = p.tokenizer.Pop()
	if !token.isIdent() {
		return p.toErrorState(token, "expected identifier (5)")
	}
	destinationPort := token.str

	err = p.matchEOLs()
	if err != nil {
		return err
	}

	mlog.Debugf("GOT: %s.%s -> %s.%s", sourceComponent, sourcePort, destinationComponent, destinationPort)

	connection := &ConnectionModel{
		Source:      sourceComponent + "." + sourcePort,
		Destination: destinationComponent + "." + destinationPort,
	}

	p.graph.Connections = append(p.graph.Connections, connection)

	return nil
}

func (p *Parser) toConfigState(argMap *ArgMap) error {
	mlog.Debugf("STATE: config")
	var err error

	err = p.matchEOLs()
	if err != nil {
		return err
	}

	done := false

	for {
		err = nil

		token := p.tokenizer.Pop()

		switch {
		case token.isIdentN("end"):
			err = p.matchEOLs()
			done = true
		case token.isIdent():
			key := token.str
			var value interface{}
			value, err = p.toConfigValueState()
			if err == nil {
				(*argMap)[key] = value
			}
		default:
			err = p.toErrorState(token, "expected identifier or 'end'")
		}

		if err != nil {
			return err
		}

		if done {
			break
		}
	}

	return nil
}

func (p *Parser) toConfigValueState() (interface{}, error) {

	var err error

	token := p.tokenizer.Pop()

	if token.typ == ColonMarker {
		token = p.tokenizer.Pop()
		value := token.str
		err = p.matchEOLs()

		// TODO: handle ints and floats here

		// if quoted string, remove quotes
		if len(value) > 2 && value[0:1] == `"` && value[len(value)-1:len(value)] == `"` {
			value = value[1 : len(value)-1]
		}
		return value, err
	}

	if token.typ == WigglyMarker {

		s := ""
		for {
			err = nil

			token = p.tokenizer.Pop()
			if token.isEOL() {
				p.skipEOLs()
				break
			}
			s += " " + token.str
		}

		if err != nil {
			return nil, err
		}

		mlog.Debugf("MAPa %s\n", s)

		value := &map[string]interface{}{}
		err = json.Unmarshal([]byte(s), value)
		if err != nil {
			return nil, err
		}
		mlog.Debugf("MAPb %#v\n", value)
		return value, err
	}

	err = p.toErrorState(token, "expected ':' or '~'")

	return nil, err
}

func (p *Parser) toMetadataState() error {
	mlog.Debugf("STATE: metadata")
	var err error

	meta := map[string]interface{}{}

	err = p.matchEOLs()
	if err != nil {
		return err
	}

	done := false

	for {
		err = nil

		token := p.tokenizer.Pop()

		switch {
		case token.isIdent():
			key := token.str
			token = p.tokenizer.Pop()
			if token.typ != ColonMarker {
				err = p.toErrorState(token, "expected ':'")
			} else {
				err = p.toErrorState(token, "???")
				token = p.tokenizer.Pop()
				value := token.str
				meta[key] = value
				err = p.matchEOLs()
			}
		case token.isIdentN("end"):
			err = p.matchEOLs()
			done = true
		default:
			err = p.toErrorState(token, "expected identifier or 'end'")
		}

		if err != nil {
			return err
		}

		if done {
			break
		}
	}

	return nil
}

func (p *Parser) toFieldValueState(field *string) error {

	token := p.tokenizer.Pop()
	if token.typ != ColonMarker {
		return p.toErrorState(token, "expected ':'")
	}

	token = p.tokenizer.Pop()
	*field = token.str

	err := p.matchEOLs()
	return err
}
