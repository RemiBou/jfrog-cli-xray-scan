package commands

// TODO: try to extract a component string from a free text line

type component string

func (c component) toString() string {
	return string(c)
}

func parse(comp string) component {
	return "gav://org.apache.httpcomponents:httpclient:4.5.9"
}
