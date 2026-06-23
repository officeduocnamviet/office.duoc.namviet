package app

import "testing"

func TestFormatOrderCode(t *testing.T) {
	cases := []struct {
		seq  int64
		want string
	}{
		{1, "DH00000001"},
		{123, "DH00000123"},
		{99999999, "DH99999999"},
		{100000000, "DH100000000"}, // vượt pad width → không cắt, vẫn đúng/duy nhất
	}
	for _, c := range cases {
		if got := formatOrderCode(c.seq); got != c.want {
			t.Errorf("formatOrderCode(%d) = %q, want %q", c.seq, got, c.want)
		}
	}
}
