package chat

import "strings"

type Style uint8

const (
	// StyleInherit indicates that a style option should be inherited from the parent. This is the default behavior.
	StyleInherit Style = iota
	// StyleOn explicitly enables a style option.
	StyleOn
	// StyleOff explicitly disables a style option.
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

// String returns the text from this Msg (and optionally extra ones from Msg.Extra) without any additional formatting.
func (m Msg) String() string {
	var builder strings.Builder
	builder.WriteString(m.Text)

	for _, extra := range m.Extra {
		builder.WriteString(extra.Text)
	}

	return builder.String()
}

// MarshalText implements encoding.TextMarshaler.
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
