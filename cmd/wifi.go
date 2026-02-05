package cmd

import (
	"fmt"
	"strings"

	"github.com/eliaseffects/qr-cli/internal/qr"
	"github.com/spf13/cobra"
)

var (
	wifiFlags    OutputFlags
	wifiSSID     string
	wifiPassword string
	wifiSecurity string
	wifiHidden   bool

	wifiCmd = &cobra.Command{
		Use:   "wifi",
		Short: "Generate QR code for WiFi network connection",
		RunE:  runWifi,
	}
)

func init() {
	wifiCmd.Flags().StringVar(&wifiSSID, "ssid", "", "Network name (required)")
	wifiCmd.Flags().StringVar(&wifiPassword, "pass", "", "Network password")
	wifiCmd.Flags().StringVar(&wifiSecurity, "security", "WPA", "Security type: WPA, WEP, nopass")
	wifiCmd.Flags().BoolVar(&wifiHidden, "hidden", false, "Network is hidden")
	_ = wifiCmd.MarkFlagRequired("ssid")

	addOutputFlags(wifiCmd, &wifiFlags, true)
	bindOutputFlags(wifiCmd)
	bindWifiFlags(wifiCmd)
}

func runWifi(cmd *cobra.Command, args []string) error {
	applyOutputConfig(cmd, &wifiFlags)
	applyWifiConfig(cmd)

	security := strings.ToUpper(strings.TrimSpace(wifiSecurity))
	switch security {
	case "NOPASS":
		security = "nopass"
		wifiPassword = ""
	case "WPA", "WEP":
	default:
		return fmt.Errorf("invalid security type: %s", wifiSecurity)
	}

	data := qr.WifiConfig{
		SSID:     wifiSSID,
		Password: wifiPassword,
		Security: security,
		Hidden:   wifiHidden,
	}.String()

	return runGenerate(data, wifiFlags, cmd.Flags().Changed("format"))
}
