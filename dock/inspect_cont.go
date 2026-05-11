package dock

// InspectCont inspects the state of a container.
func InspectCont(c *Client, name string) (*ContInfo, error) {
	cont := NewCont(c, name)
	return cont.Inspect()
}

// HasCont checks if a container exists.
func HasCont(c *Client, name string) (bool, error) {
	cont := NewCont(c, name)
	return cont.Exists()
}
