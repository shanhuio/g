package dock

import (
	"reflect"
	"testing"
)

func TestUnmapEnv(t *testing.T) {
	for _, test := range []struct {
		name string
		in   map[string]string
		want []string
	}{
		{name: "nil", in: nil, want: nil},
		{name: "empty", in: map[string]string{}, want: nil},
		{
			name: "single",
			in:   map[string]string{"FOO": "bar"},
			want: []string{"FOO=bar"},
		},
		{
			name: "sorted",
			in:   map[string]string{"B": "2", "A": "1", "C": "3"},
			want: []string{"A=1", "B=2", "C=3"},
		},
		{
			name: "empty value",
			in:   map[string]string{"EMPTY": ""},
			want: []string{"EMPTY="},
		},
		{
			name: "value with equals",
			in:   map[string]string{"PATH": "/bin:/usr/bin", "EQ": "a=b=c"},
			want: []string{"EQ=a=b=c", "PATH=/bin:/usr/bin"},
		},
		{
			name: "empty key",
			in:   map[string]string{"": "value"},
			want: []string{"=value"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := unmapEnv(test.in)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("unmapEnv(%v), got %v, want %v", test.in, got, test.want)
			}
		})
	}
}
