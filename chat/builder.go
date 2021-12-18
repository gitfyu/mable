package chat

// Builder is a utility for constructing more complex Msg objects. An instance should be obtained using NewBuilder.
type Builder struct {
	parts []Msg
}

// TODO it might be useful to use a sync.Pool for Builder instances, similar to how network.PacketBuilder works

// NewBuilder constructs a new Builder instance with a message
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

func (b *Builder) Bold() *Builder {
	b.current().Bold = StyleOn
	return b
}

func (b *Builder) NotBold() *Builder {
	b.current().Bold = StyleOff
	return b
}

func (b *Builder) Italic() *Builder {
	b.current().Italic = StyleOn
	return b
}

func (b *Builder) NotItalic() *Builder {
	b.current().Italic = StyleOff
	return b
}

func (b *Builder) NotUnderlined() *Builder {
	b.current().Underlined = StyleOff
	return b
}

func (b *Builder) Strikethrough() *Builder {
	b.current().Strikethrough = StyleOn
	return b
}

func (b *Builder) NotStrikethrough() *Builder {
	b.current().Strikethrough = StyleOff
	return b
}

func (b *Builder) Obfuscated() *Builder {
	b.current().Obfuscated = StyleOn
	return b
}

func (b *Builder) NotObfuscated() *Builder {
	b.current().Obfuscated = StyleOff
	return b
}

func (b *Builder) Color(c Color) *Builder {
	b.current().Color = c
	return b
}

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
