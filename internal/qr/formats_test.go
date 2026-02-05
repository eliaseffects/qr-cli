package qr_test

import (
	"strings"
	"testing"

	"github.com/eliaseffects/qr-cli/internal/qr"
)

func TestWifiConfig(t *testing.T) {
	tests := []struct {
		config qr.WifiConfig
		want   string
	}{
		{
			qr.WifiConfig{SSID: "MyNetwork", Password: "secret123", Security: "WPA"},
			"WIFI:T:WPA;S:MyNetwork;P:secret123;;",
		},
		{
			qr.WifiConfig{SSID: "Open;Network", Password: "", Security: "nopass"},
			"WIFI:T:nopass;S:Open\\;Network;P:;;",
		},
		{
			qr.WifiConfig{SSID: "Hidden", Password: "pass", Security: "WPA", Hidden: true},
			"WIFI:T:WPA;S:Hidden;P:pass;H:true;",
		},
	}

	for _, tt := range tests {
		if got := tt.config.String(); got != tt.want {
			t.Errorf("WifiConfig.String() = %q, want %q", got, tt.want)
		}
	}
}

func TestVCard(t *testing.T) {
	v := qr.VCard{
		Name:  "John Doe",
		Phone: "+1234567890",
		Email: "john@example.com",
	}

	result := v.String()

	if !strings.Contains(result, "BEGIN:VCARD") {
		t.Error("missing BEGIN:VCARD")
	}
	if !strings.Contains(result, "FN:John Doe") {
		t.Error("missing FN field")
	}
	if !strings.Contains(result, "TEL:+1234567890") {
		t.Error("missing TEL field")
	}
}
