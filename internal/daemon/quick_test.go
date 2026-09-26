package daemon

import "testing"

func TestExtractURL(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{"INF | Your quick Tunnel has been created! Visit it at https://demo.trycloudflare.com", "https://demo.trycloudflare.com"},
		{"no public address yet", ""},
	}
	for _, tt := range tests {
		if got := extractURL(tt.line); got != tt.want {
			t.Errorf("extractURL(%q) = %q, want %q", tt.line, got, tt.want)
		}
	}
}
