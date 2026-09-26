package quickshare

import (
	"net/url"
	"strings"
	"testing"
)

func TestTelegramShareURL(t *testing.T) {
	target := "https://demo.trycloudflare.com/a?x=1"
	got := TelegramShareURL(target, "cftunnel 分享地址: "+target)
	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("parse share URL: %v", err)
	}
	q := parsed.Query()
	if q.Get("url") != target {
		t.Fatalf("url query = %q, want %q", q.Get("url"), target)
	}
	if q.Get("text") == "" {
		t.Fatal("text query is empty")
	}
}

func TestQRText(t *testing.T) {
	got, err := QRText("https://demo.trycloudflare.com")
	if err != nil {
		t.Fatalf("QRText: %v", err)
	}
	if !strings.Contains(got, "█") || !strings.Contains(got, "\n") {
		t.Fatalf("QR output does not look like a terminal QR code: %q", got)
	}
}
