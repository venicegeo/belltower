package dsls

import "encoding/json"

type Graph struct {
	Id       Id        `json:"id"`
	Name     string    `json:"name"`
	Metadata *Metadata `json:"metadata"`

	Components  []*Component  `json:"components"`
	Connections []*Connection `json:"connections"`
}

type Id string

type Metadata struct {
	Contact     string `json:"contact,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
}

// Ticker, Adder, etc
type ComponentType struct {
	Id       Id        `json:"id"`
	Name     string    `json:"name"`
	Metadata *Metadata `json:"metadata,omitempty"`

	Input  *Port `json:"input,omitempty"`
	Output *Port `json:"output,omitempty"`
	Config *Port `json:"config,omitempty"`
}

// MyTicker, AdderFortyTwo, etc
type Component struct {
	Id       Id        `json:"id"`
	Name     string    `json:"name"`
	Metadata *Metadata `json:"metadata"`

	Type string `json:"type"`

	Precondition  string `json:"precondition,omitempty"`
	Postcondition string `json:"postcondition,omitempty"`
	Config        Map
}

type PortTypeEnum string

const (
	InvalidPort PortTypeEnum = "INVALID"
	InputPort                = "input"
	OutputPort               = "output"
	ConfigPort               = "config"
)

type Port struct {
	Id       Id           `json:"id"`
	Name     string       `json:"name"`
	DataType *DataType    `json:"datatype"`
	PortType PortTypeEnum `json:"porttype"`
}

type Connection struct {
	Id          Id     `json:"id"`
	Name        string `json:"name"`
	Source      string `json:"source"`      // component.port
	Destination string `json:"destination"` // component.port
}

//---------------------------------------------------------------------

type Validater interface {
	Validate() error
}

func (m *Metadata) Validate() error      { return nil }
func (p *Port) Validate() error          { return nil }
func (p *Component) Validate() error     { return nil }
func (p *ComponentType) Validate() error { return nil }
func (p *Graph) Validate() error         { return nil }
func (p *Connection) Validate() error    { return nil }

func NewObjectFromJSON(jsn string, obj Validater) (interface{}, error) {
	err := json.Unmarshal([]byte(jsn), obj)
	if err != nil {
		return nil, err
	}

	err = obj.Validate()
	if err != nil {
		return nil, err
	}

	return obj, nil
}
