package chat

type Style uint8

const (
	// StyleInherit indicates that a style option should be inherited from the parent
	StyleInherit Style = iota
	// StyleOn enables a style option
	StyleOn
	// StyleOff disables a style option
	StyleOff
)

// Msg represents a chat message. For complex messages, it might be easier to use Builder to construct them.
type Msg struct {
	Text          string `json:"text"`
	Bold          Style  `json:"bold,omitempty"`
	Italic        Style  `json:"italic,omitempty"`
	Underlined    Style  `json:"underlined,omitempty"`
	Strikethrough Style  `json:"strikethrough,omitempty"`
	Obfuscated    Style  `json:"obfuscated,omitempty"`
	Color         Color  `json:"color,omitempty"`
	Extra         []Msg  `json:"extra,omitempty"`
}

func (s Style) MarshalText() ([]byte, error) {
	switch s {
	case StyleOn:
		return []byte("true"), nil
	case StyleOff:
		return []byte("false"), nil
	default:
		return nil, nil
	}
}
