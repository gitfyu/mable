/*
Package chat is used to construct messages that can be converted to Minecraft's JSON chat format.
For more information about this format, visit https:wiki.vg/index.php?title=Chat&oldid=5640. Please note that
this package does not implement new functionality from versions newer than 1.8.X, such as web colors.

A simple message can be constructed like this:
	msg := Msg{
		Text: "Hello",
		Bold: StyleOn,
		Color: ColorGreen,
	}
You can use json.Encoder to convert msg to JSON.

For more complex messages, you might want to use a Builder:
	msg := NewBuilder("Hello, ").Bold().Color(ColorRed).
				Append("world").NotBold().Color(ColorYellow).
				Build()

If you change a style or color setting in a Builder, that change applies to the current text but also to all
text that is appended later. If you want to prevent this, you must explicitly disable the style or change the color
back after calling Append.
*/
package chat
