package dock

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"shanhu.io/g/errcode"
)

type trackingReader struct {
	io.Reader
	closed bool
}

func (r *trackingReader) Close() error {
	r.closed = true
	return nil
}

func TestPrintStreamMessageEmpty(t *testing.T) {
	var out bytes.Buffer
	err := printStreamMessage(io.NopCloser(strings.NewReader("")), &out)
	if err != nil {
		t.Errorf("got %v, want nil", err)
	}
}

func TestPrintStreamMessageStream(t *testing.T) {
	input := `{"stream":"Step 1/3\n"}{"stream":"Step 2/3\n"}`
	var out bytes.Buffer
	if err := printStreamMessage(
		io.NopCloser(strings.NewReader(input)), &out,
	); err != nil {
		t.Fatalf("printStreamMessage: %v", err)
	}
	s := out.String()
	if !strings.Contains(s, "Step 1/3") {
		t.Errorf("missing 'Step 1/3' in output: %q", s)
	}
	if !strings.Contains(s, "Step 2/3") {
		t.Errorf("missing 'Step 2/3' in output: %q", s)
	}
}

func TestPrintStreamMessageError(t *testing.T) {
	input := `{"error":"build failed"}`
	var out bytes.Buffer
	err := printStreamMessage(
		io.NopCloser(strings.NewReader(input)), &out,
	)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if !errcode.IsInternal(err) {
		t.Errorf("got %v (code %q), want internal", err, errcode.Of(err))
	}
	if !strings.Contains(err.Error(), "build failed") {
		t.Errorf("error %q does not contain 'build failed'", err)
	}
}

func TestPrintStreamMessageLastErrorWins(t *testing.T) {
	input := `{"error":"first"}{"error":"second"}`
	var out bytes.Buffer
	err := printStreamMessage(
		io.NopCloser(strings.NewReader(input)), &out,
	)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if !strings.Contains(err.Error(), "second") {
		t.Errorf("got %q, want last error 'second'", err)
	}
}

func TestPrintStreamMessageStreamThenError(t *testing.T) {
	input := `{"stream":"working...\n"}{"error":"boom"}`
	var out bytes.Buffer
	err := printStreamMessage(
		io.NopCloser(strings.NewReader(input)), &out,
	)
	if err == nil {
		t.Fatal("got nil, want error")
	}
	if !errcode.IsInternal(err) {
		t.Errorf("got %v, want internal", err)
	}
	if !strings.Contains(out.String(), "working...") {
		t.Errorf("expected stream content in output, got %q", out.String())
	}
}

func TestPrintStreamMessageMalformedJSON(t *testing.T) {
	input := `{"stream":"ok"} not-json`
	var out bytes.Buffer
	err := printStreamMessage(
		io.NopCloser(strings.NewReader(input)), &out,
	)
	if err == nil {
		t.Fatal("got nil, want decode error")
	}
}

func TestPrintStreamMessageClosesReader(t *testing.T) {
	r := &trackingReader{Reader: strings.NewReader("")}
	var out bytes.Buffer
	if err := printStreamMessage(r, &out); err != nil {
		t.Fatalf("printStreamMessage: %v", err)
	}
	if !r.closed {
		t.Errorf("reader was not closed")
	}
}

func TestProgressPrinterStream(t *testing.T) {
	var out bytes.Buffer
	p := &progressPrinter{out: &out}
	p.Print(&streamMessage{Stream: "hello world"})
	if !strings.Contains(out.String(), "hello world") {
		t.Errorf("got %q, want to contain 'hello world'", out.String())
	}
	if p.ids != nil {
		t.Errorf("ids = %v, want nil for stream-only message", p.ids)
	}
}

func TestProgressPrinterIDTracking(t *testing.T) {
	var out bytes.Buffer
	p := &progressPrinter{out: &out}
	p.Print(&streamMessage{ID: "abc", Status: "downloading"})
	p.Print(&streamMessage{ID: "def", Status: "downloading"})
	p.Print(&streamMessage{ID: "abc", Status: "extracting"})

	if len(p.ids) != 2 {
		t.Fatalf("got %d ids, want 2: %v", len(p.ids), p.ids)
	}
	if p.ids[0] != "abc" || p.ids[1] != "def" {
		t.Errorf("ids = %v, want [abc def]", p.ids)
	}
}

func TestProgressPrinterStreamResetsIDs(t *testing.T) {
	var out bytes.Buffer
	p := &progressPrinter{out: &out}
	p.Print(&streamMessage{ID: "abc", Status: "downloading"})
	if len(p.ids) != 1 {
		t.Fatalf("ids = %v, want [abc]", p.ids)
	}
	p.Print(&streamMessage{Stream: "Step 1\n"})
	if p.ids != nil {
		t.Errorf("ids = %v, want nil after stream", p.ids)
	}
}

func TestProgressPrinterEmptyMessage(t *testing.T) {
	var out bytes.Buffer
	p := &progressPrinter{out: &out}
	p.Print(&streamMessage{})
	if len(p.ids) != 0 {
		t.Errorf("ids = %v, want empty", p.ids)
	}
	if out.Len() != 0 {
		t.Errorf("output = %q, want empty", out.String())
	}
}

func TestProgressPrinterProgressDetail(t *testing.T) {
	var out bytes.Buffer
	p := &progressPrinter{out: &out}
	p.Print(&streamMessage{
		ID:     "layer1",
		Status: "downloading",
		ProgressDetail: &progressDetail{
			Current: 50,
			Total:   100,
		},
	})
	s := out.String()
	if !strings.Contains(s, "layer1") {
		t.Errorf("missing id 'layer1' in output: %q", s)
	}
	if !strings.Contains(s, "downloading") {
		t.Errorf("missing status 'downloading' in output: %q", s)
	}
	if !strings.Contains(s, "50.00%") {
		t.Errorf("missing percentage '50.00%%' in output: %q", s)
	}
}
