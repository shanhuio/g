package caco3

type loadTracer struct {
	trace []string
	m     map[string]bool
}

func newLoadTracer() *loadTracer {
	return &loadTracer{
		m: make(map[string]bool),
	}
}

func (t *loadTracer) push(name string) bool {
	if t.m[name] {
		return false
	}
	t.trace = append(t.trace, name)
	t.m[name] = true
	return true
}

func (t *loadTracer) pop() {
	n := len(t.trace)
	if n == 0 {
		return
	}
	last := t.trace[n-1]
	delete(t.m, last)
	t.trace = t.trace[:n-1]
}

func (t *loadTracer) stack() []string {
	return t.trace
}
