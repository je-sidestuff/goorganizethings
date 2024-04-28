package focusblock

import "testing"

func TestFocusDay(t *testing.T) {
	got := PrintFocusDays(1900, 6, 6, 1, "Hello, World!")
	want := "1900-6-6 Hello, World!"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestFocusBlock(t *testing.T) {
	got := PrintFocusDays(1900, 6, 6, 3, "Hello, World!")
	want := `1900-6-6 Hello, World!
1900-6-7 Hello, World!
1900-6-8 Hello, World!`

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
