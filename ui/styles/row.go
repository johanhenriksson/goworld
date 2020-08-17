package styles


type Row struct {

}

/*

rect {
	layout: column

	text {
	}

	rect {
		layout: "row"
		align: ${align}
		padding: 5px

		text {
			color: rgb(0.5, 0.3, 0.5)
			value: "hello " ${msg} " pls"
		}
	}
}

if oldState != newState
  render()

NewElement(
	"rect",
	Definition: map[string]Value {
		"layout": Const("row"),
		"align": State("align"),
	},
	Children: []Element {
		// children
		Text {
			Value: []Values {
				String("hello "),
				State("msg"),
				String(" pls"),
			},
			Definition: map[string]Value {
				"font-size": String("16px"),
				"background-color": RGB(0.5, 0.3, 0.5)
			},
		},

	}
}