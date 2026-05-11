package dock

import (
	"bytes"
	"encoding/binary"
	"testing"
)

func dockerFrame(stream byte, payload []byte) []byte {
	frame := make([]byte, 8+len(payload))
	frame[0] = stream
	binary.BigEndian.PutUint32(frame[4:8], uint32(len(payload)))
	copy(frame[8:], payload)
	return frame
}

func TestLogDispatcherRoutes(t *testing.T) {
	for _, test := range []struct {
		name           string
		stream         byte
		payload        string
		wantStdout     string
		wantStderr     string
	}{
		{
			name:    "stdout (stream 1)",
			stream:  1,
			payload: "hello\n",
			wantStdout: "hello\n",
		},
		{
			name:    "stream 0 also stdout",
			stream:  0,
			payload: "out\n",
			wantStdout: "out\n",
		},
		{
			name:    "stderr (stream 2)",
			stream:  2,
			payload: "boom\n",
			wantStderr: "boom\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			d := &logDispatcher{Stdout: &stdout, Stderr: &stderr}
			frame := dockerFrame(test.stream, []byte(test.payload))
			if err := d.pipe(bytes.NewReader(frame)); err != nil {
				t.Fatalf("pipe: %v", err)
			}
			if stdout.String() != test.wantStdout {
				t.Errorf(
					"stdout = %q, want %q",
					stdout.String(), test.wantStdout,
				)
			}
			if stderr.String() != test.wantStderr {
				t.Errorf(
					"stderr = %q, want %q",
					stderr.String(), test.wantStderr,
				)
			}
		})
	}
}

func TestLogDispatcherInterleaved(t *testing.T) {
	var stdout, stderr bytes.Buffer
	d := &logDispatcher{Stdout: &stdout, Stderr: &stderr}

	var buf bytes.Buffer
	buf.Write(dockerFrame(1, []byte("o1\n")))
	buf.Write(dockerFrame(2, []byte("e1\n")))
	buf.Write(dockerFrame(1, []byte("o2\n")))
	buf.Write(dockerFrame(2, []byte("e2\n")))

	if err := d.pipe(&buf); err != nil {
		t.Fatalf("pipe: %v", err)
	}
	if stdout.String() != "o1\no2\n" {
		t.Errorf("stdout = %q, want %q", stdout.String(), "o1\no2\n")
	}
	if stderr.String() != "e1\ne2\n" {
		t.Errorf("stderr = %q, want %q", stderr.String(), "e1\ne2\n")
	}
}

func TestLogDispatcherEmptyInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	d := &logDispatcher{Stdout: &stdout, Stderr: &stderr}
	if err := d.pipe(bytes.NewReader(nil)); err != nil {
		t.Errorf("pipe: %v", err)
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Errorf("unexpected output: stdout=%q stderr=%q",
			stdout.String(), stderr.String())
	}
}

func TestLogDispatcherEmptyPayload(t *testing.T) {
	var stdout, stderr bytes.Buffer
	d := &logDispatcher{Stdout: &stdout, Stderr: &stderr}
	frame := dockerFrame(1, nil)
	if err := d.pipe(bytes.NewReader(frame)); err != nil {
		t.Errorf("pipe: %v", err)
	}
	if stdout.Len() != 0 || stderr.Len() != 0 {
		t.Errorf("unexpected output: stdout=%q stderr=%q",
			stdout.String(), stderr.String())
	}
}

func TestLogDispatcherInvalidStream(t *testing.T) {
	var stdout, stderr bytes.Buffer
	d := &logDispatcher{Stdout: &stdout, Stderr: &stderr}
	frame := dockerFrame(99, []byte("garbage"))
	err := d.pipe(bytes.NewReader(frame))
	if err == nil {
		t.Fatal("got nil, want error for invalid stream byte")
	}
}

func TestLogDispatcherTruncatedHeader(t *testing.T) {
	var stdout, stderr bytes.Buffer
	d := &logDispatcher{Stdout: &stdout, Stderr: &stderr}
	// Only 5 of 8 header bytes — ReadFull fails with ErrUnexpectedEOF.
	err := d.pipe(bytes.NewReader([]byte{1, 0, 0, 0, 0}))
	if err == nil {
		t.Fatal("got nil, want error for truncated header")
	}
}

func TestLogDispatcherTruncatedPayload(t *testing.T) {
	var stdout, stderr bytes.Buffer
	d := &logDispatcher{Stdout: &stdout, Stderr: &stderr}
	// Header claims 10 bytes payload, only 3 supplied.
	buf := dockerFrame(1, []byte("xxxxxxxxxx"))
	buf = buf[:8+3]
	err := d.pipe(bytes.NewReader(buf))
	if err == nil {
		t.Fatal("got nil, want error for truncated payload")
	}
}

func TestLogSinkRoundTrip(t *testing.T) {
	var stdout, stderr bytes.Buffer
	s := newLogSink(&stdout, &stderr)

	if _, err := s.Write(dockerFrame(1, []byte("hello\n"))); err != nil {
		t.Fatalf("Write stdout frame: %v", err)
	}
	if _, err := s.Write(dockerFrame(2, []byte("oops\n"))); err != nil {
		t.Fatalf("Write stderr frame: %v", err)
	}

	if err := s.waitDone(); err != nil {
		t.Fatalf("waitDone: %v", err)
	}
	if stdout.String() != "hello\n" {
		t.Errorf("stdout = %q, want %q", stdout.String(), "hello\n")
	}
	if stderr.String() != "oops\n" {
		t.Errorf("stderr = %q, want %q", stderr.String(), "oops\n")
	}
}

func TestLogSinkEmpty(t *testing.T) {
	var stdout, stderr bytes.Buffer
	s := newLogSink(&stdout, &stderr)
	if err := s.waitDone(); err != nil {
		t.Errorf("waitDone with no writes: %v", err)
	}
}

func TestLogSinkNilDefaults(t *testing.T) {
	// nil stdout/stderr should default to os.Stdout/os.Stderr without panic.
	// With no writes, waitDone closes the writer and the goroutine returns
	// on EOF — nothing is emitted, so this is safe to run.
	s := newLogSink(nil, nil)
	if s.d.Stdout == nil {
		t.Error("Stdout default not set")
	}
	if s.d.Stderr == nil {
		t.Error("Stderr default not set")
	}
	if err := s.waitDone(); err != nil {
		t.Errorf("waitDone: %v", err)
	}
}
