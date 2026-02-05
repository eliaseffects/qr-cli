package qr

import (
	"fmt"
	"strings"
)

type WifiConfig struct {
	SSID     string
	Password string
	Security string // WPA, WEP, nopass
	Hidden   bool
}

func (w WifiConfig) String() string {
	hidden := ""
	if w.Hidden {
		hidden = "H:true;"
	}

	ssid := escapeWifi(w.SSID)
	pass := escapeWifi(w.Password)

	suffix := ";"
	if hidden != "" {
		suffix = ""
	}

	return fmt.Sprintf("WIFI:T:%s;S:%s;P:%s;%s%s", w.Security, ssid, pass, hidden, suffix)
}

type VCard struct {
	Name    string
	Phone   string
	Email   string
	Org     string
	Title   string
	URL     string
	Address string
}

func (v VCard) String() string {
	var b strings.Builder
	b.WriteString("BEGIN:VCARD\n")
	b.WriteString("VERSION:3.0\n")

	if v.Name != "" {
		b.WriteString(fmt.Sprintf("FN:%s\n", v.Name))
		parts := strings.SplitN(v.Name, " ", 2)
		if len(parts) == 2 {
			b.WriteString(fmt.Sprintf("N:%s;%s;;;\n", parts[1], parts[0]))
		} else {
			b.WriteString(fmt.Sprintf("N:%s;;;;\n", v.Name))
		}
	}
	if v.Phone != "" {
		b.WriteString(fmt.Sprintf("TEL:%s\n", v.Phone))
	}
	if v.Email != "" {
		b.WriteString(fmt.Sprintf("EMAIL:%s\n", v.Email))
	}
	if v.Org != "" {
		b.WriteString(fmt.Sprintf("ORG:%s\n", v.Org))
	}
	if v.Title != "" {
		b.WriteString(fmt.Sprintf("TITLE:%s\n", v.Title))
	}
	if v.URL != "" {
		b.WriteString(fmt.Sprintf("URL:%s\n", v.URL))
	}
	if v.Address != "" {
		b.WriteString(fmt.Sprintf("ADR:;;%s;;;;\n", v.Address))
	}

	b.WriteString("END:VCARD")
	return b.String()
}

func escapeWifi(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, ";", "\\;")
	s = strings.ReplaceAll(s, ",", "\\,")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
