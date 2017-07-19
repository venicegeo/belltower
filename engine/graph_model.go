package engine

type GraphModel struct {
	Name string `json:"name"`

	Components  []*ComponentModel  `json:"components"`
	Connections []*ConnectionModel `json:"connections"`
}

type ComponentModel struct {
	Name string `json:"name"`

	Type string `json:"type"`

	Precondition  string `json:"precondition,omitempty"`
	Postcondition string `json:"postcondition,omitempty"`
	Config        ArgMap
}

type ConnectionModel struct {
	Source      string `json:"source"`      // component.port
	Destination string `json:"destination"` // component.port
}

func (m *GraphModel) Validate() error               { return nil } // TODO
func (m *GraphModel) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *GraphModel) WriteToJSON() (string, error)  { return WriteToJSON(m) }

func (m *ComponentModel) Validate() error               { return nil } // TODO
func (m *ComponentModel) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *ComponentModel) WriteToJSON() (string, error)  { return WriteToJSON(m) }

func (m *ConnectionModel) Validate() error               { return nil } // TODO
func (m *ConnectionModel) ReadFromJSON(jsn string) error { return ReadFromJSON(jsn, m) }
func (m *ConnectionModel) WriteToJSON() (string, error)  { return WriteToJSON(m) }
