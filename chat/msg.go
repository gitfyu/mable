package chat

type Msg struct {
	Text          string `json:"text"`
	Bold          bool   `json:"bold,omitempty"`
	Italic        bool   `json:"italic,omitempty"`
	Underlined    bool   `json:"underlined,omitempty"`
	Strikethrough bool   `json:"strikethrough,omitempty"`
	Obfuscated    bool   `json:"obfuscated,omitempty"`
	Color         Color  `json:"color,omitempty"`
	Extra         []Msg  `json:"extra,omitempty"`
}
