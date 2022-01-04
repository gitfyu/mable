package chat

// Builder is a utility for constructing more complex Msg objects. An instance should be obtained using NewBuilder.
type Builder struct {
	parts []Msg
}

// TODO it might be useful to use a sync.Pool for Builder instances, similar to how network.PacketBuilder works

// NewBuilder constructs a new Builder instance with initial text.
func NewBuilder(text string) *Builder {
	return &Builder{
		parts: []Msg{
			{Text: text},
		},
	}
}

func (b *Builder) current() *Msg {
	return &b.parts[len(b.parts)-1]
}

// Bold enables the bold style for the current text.
func (b *Builder) Bold() *Builder {
	b.current().Bold = StyleOn
	return b
}

// NotBold disables the bold style for the current text.
func (b *Builder) NotBold() *Builder {
	b.current().Bold = StyleOff
	return b
}

// Italic enables the italic style for the current text.
func (b *Builder) Italic() *Builder {
	b.current().Italic = StyleOn
	return b
}

// NotItalic disables the italic style for the current text.
func (b *Builder) NotItalic() *Builder {
	b.current().Italic = StyleOff
	return b
}

// Underlined enables the underlined style for the current text.
func (b *Builder) Underlined() *Builder {
	b.current().Underlined = StyleOn
	return b
}

// NotUnderlined disables the underlined style for the current text.
func (b *Builder) NotUnderlined() *Builder {
	b.current().Underlined = StyleOff
	return b
}

// Strikethrough enables the strikethrough style for the current text.
func (b *Builder) Strikethrough() *Builder {
	b.current().Strikethrough = StyleOn
	return b
}

// NotStrikethrough disables the strikethrough style for the current text.
func (b *Builder) NotStrikethrough() *Builder {
	b.current().Strikethrough = StyleOff
	return b
}

// Obfuscated enables the obfuscated style for the current text.
func (b *Builder) Obfuscated() *Builder {
	b.current().Obfuscated = StyleOn
	return b
}

// NotObfuscated disables the obfuscated style for the current text.
func (b *Builder) NotObfuscated() *Builder {
	b.current().Obfuscated = StyleOff
	return b
}

// Color changes the color for the current text.
func (b *Builder) Color(c Color) *Builder {
	b.current().Color = c
	return b
}

// Append adds a new child message with the specified text.
func (b *Builder) Append(text string) *Builder {
	b.parts = append(b.parts, Msg{Text: text})
	return b
}

// Build constructs a Msg using all the previously applied options. You should not use and/or keep a reference to the
// Builder anymore after calling this function.
func (b *Builder) Build() *Msg {
	msg := b.parts[0]
	if len(b.parts) > 1 {
		msg.Extra = b.parts[1:]
	}

	return &msg
}
