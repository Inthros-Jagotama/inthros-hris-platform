package mask

import "testing"

func TestPartialMask(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  string
	}{
		{"17 chars: last 4 visible", "3201010101985678", "************5678"},
		{"exactly 10 chars: last 4 visible", "1234567890", "******7890"},
		{"9 chars: last 3 visible", "123456789", "******789"},
		{"6 chars: last 3 visible", "123456", "***456"},
		{"5 chars: fully masked", "12345", "*****"},
		{"3 chars: fully masked", "123", "***"},
		{"1 char: fully masked", "1", "*"},
		{"empty stays empty", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := PartialMask(c.input)
			if got != c.want {
				t.Errorf("PartialMask(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}
