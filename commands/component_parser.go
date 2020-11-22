package commands

// TODO: try to extract a component string from a free text line

type component string

func (c component) toString() string {
	return string(c)
}

func parse(comp string) component {
	// TODO: implement component detection
	return component(comp)
}
