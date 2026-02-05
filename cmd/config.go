package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	configErr error
)

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("qr-cli")
		viper.AddConfigPath(".")
		if home, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(filepath.Join(home, ".config", "qr-cli"))
			viper.AddConfigPath(home)
		}
	}

	viper.SetEnvPrefix("QR")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	setConfigDefaults()

	if err := viper.ReadInConfig(); err != nil {
		if cfgFile != "" {
			configErr = fmt.Errorf("failed to read config file: %w", err)
		}
	}
}

func setConfigDefaults() {
	viper.SetDefault("output", "")
	viper.SetDefault("size", 256)
	viper.SetDefault("format", "png")
	viper.SetDefault("level", "M")
	viper.SetDefault("fg", "#000000")
	viper.SetDefault("bg", "#ffffff")
	viper.SetDefault("border", 4)
	viper.SetDefault("logo", "")
	viper.SetDefault("logo-scale", 0.2)
	viper.SetDefault("invert", false)
	viper.SetDefault("terminal-color", false)
	viper.SetDefault("terminal", false)
	viper.SetDefault("open", false)
	viper.SetDefault("copy", false)
	viper.SetDefault("quiet", false)

	viper.SetDefault("wifi.ssid", "")
	viper.SetDefault("wifi.pass", "")
	viper.SetDefault("wifi.security", "WPA")
	viper.SetDefault("wifi.hidden", false)

	viper.SetDefault("vcard.name", "")
	viper.SetDefault("vcard.phone", "")
	viper.SetDefault("vcard.email", "")
	viper.SetDefault("vcard.org", "")
	viper.SetDefault("vcard.title", "")
	viper.SetDefault("vcard.url", "")
	viper.SetDefault("vcard.address", "")

	viper.SetDefault("batch.file", "")
	viper.SetDefault("batch.dir", "./qr-output")
	viper.SetDefault("batch.size", 256)
	viper.SetDefault("batch.format", "png")
	viper.SetDefault("batch.prefix", "qr-")
	viper.SetDefault("batch.quiet", false)
}

func bindOutputFlags(cmd *cobra.Command) {
	bindFlag(cmd, "output", "output")
	bindFlag(cmd, "size", "size")
	bindFlag(cmd, "format", "format")
	bindFlag(cmd, "level", "level")
	bindFlag(cmd, "fg", "fg")
	bindFlag(cmd, "bg", "bg")
	bindFlag(cmd, "border", "border")
	bindFlag(cmd, "logo", "logo")
	bindFlag(cmd, "logo-scale", "logo-scale")
	bindFlag(cmd, "invert", "invert")
	bindFlag(cmd, "terminal-color", "terminal-color")
	bindFlag(cmd, "terminal", "terminal")
	bindFlag(cmd, "open", "open")
	bindFlag(cmd, "copy", "copy")
	bindFlag(cmd, "quiet", "quiet")
}

func bindFlag(cmd *cobra.Command, key, flag string) {
	if f := cmd.Flags().Lookup(flag); f != nil {
		_ = viper.BindPFlag(key, f)
	}
}

func applyOutputConfig(cmd *cobra.Command, flags *OutputFlags) {
	if !cmd.Flags().Changed("output") && viper.IsSet("output") {
		flags.OutputPath = viper.GetString("output")
	}
	if !cmd.Flags().Changed("size") && viper.IsSet("size") {
		flags.Size = viper.GetInt("size")
	}
	if !cmd.Flags().Changed("format") && viper.IsSet("format") {
		flags.Format = viper.GetString("format")
	}
	if !cmd.Flags().Changed("level") && viper.IsSet("level") {
		flags.Level = viper.GetString("level")
	}
	if !cmd.Flags().Changed("fg") && viper.IsSet("fg") {
		flags.FgColor = viper.GetString("fg")
	}
	if !cmd.Flags().Changed("bg") && viper.IsSet("bg") {
		flags.BgColor = viper.GetString("bg")
	}
	if !cmd.Flags().Changed("border") && viper.IsSet("border") {
		flags.Border = viper.GetInt("border")
	}
	if !cmd.Flags().Changed("logo") && viper.IsSet("logo") {
		flags.LogoPath = viper.GetString("logo")
	}
	if !cmd.Flags().Changed("logo-scale") && viper.IsSet("logo-scale") {
		flags.LogoScale = viper.GetFloat64("logo-scale")
	}
	if !cmd.Flags().Changed("invert") && viper.IsSet("invert") {
		flags.Invert = viper.GetBool("invert")
	}
	if !cmd.Flags().Changed("terminal-color") && viper.IsSet("terminal-color") {
		flags.TermColor = viper.GetBool("terminal-color")
	}
	if !cmd.Flags().Changed("terminal") && viper.IsSet("terminal") {
		flags.Terminal = viper.GetBool("terminal")
	}
	if !cmd.Flags().Changed("open") && viper.IsSet("open") {
		flags.OpenViewer = viper.GetBool("open")
	}
	if !cmd.Flags().Changed("copy") && viper.IsSet("copy") {
		flags.CopyClip = viper.GetBool("copy")
	}
	if !cmd.Flags().Changed("quiet") && viper.IsSet("quiet") {
		flags.Quiet = viper.GetBool("quiet")
	}
}

func applyWifiConfig(cmd *cobra.Command) {
	if !cmd.Flags().Changed("ssid") && viper.IsSet("wifi.ssid") {
		wifiSSID = viper.GetString("wifi.ssid")
	}
	if !cmd.Flags().Changed("pass") && viper.IsSet("wifi.pass") {
		wifiPassword = viper.GetString("wifi.pass")
	}
	if !cmd.Flags().Changed("security") && viper.IsSet("wifi.security") {
		wifiSecurity = viper.GetString("wifi.security")
	}
	if !cmd.Flags().Changed("hidden") && viper.IsSet("wifi.hidden") {
		wifiHidden = viper.GetBool("wifi.hidden")
	}
}

func applyVCardConfig(cmd *cobra.Command) {
	if !cmd.Flags().Changed("name") && viper.IsSet("vcard.name") {
		vcardName = viper.GetString("vcard.name")
	}
	if !cmd.Flags().Changed("phone") && viper.IsSet("vcard.phone") {
		vcardPhone = viper.GetString("vcard.phone")
	}
	if !cmd.Flags().Changed("email") && viper.IsSet("vcard.email") {
		vcardEmail = viper.GetString("vcard.email")
	}
	if !cmd.Flags().Changed("org") && viper.IsSet("vcard.org") {
		vcardOrg = viper.GetString("vcard.org")
	}
	if !cmd.Flags().Changed("title") && viper.IsSet("vcard.title") {
		vcardTitle = viper.GetString("vcard.title")
	}
	if !cmd.Flags().Changed("url") && viper.IsSet("vcard.url") {
		vcardURL = viper.GetString("vcard.url")
	}
	if !cmd.Flags().Changed("address") && viper.IsSet("vcard.address") {
		vcardAddress = viper.GetString("vcard.address")
	}
}

func applyBatchConfig(cmd *cobra.Command) {
	if !cmd.Flags().Changed("file") && viper.IsSet("batch.file") {
		batchCfg.File = viper.GetString("batch.file")
	}
	if !cmd.Flags().Changed("dir") && viper.IsSet("batch.dir") {
		batchCfg.Dir = viper.GetString("batch.dir")
	}
	if !cmd.Flags().Changed("size") && viper.IsSet("batch.size") {
		batchCfg.Size = viper.GetInt("batch.size")
	}
	if !cmd.Flags().Changed("format") && viper.IsSet("batch.format") {
		batchCfg.Format = viper.GetString("batch.format")
	}
	if !cmd.Flags().Changed("prefix") && viper.IsSet("batch.prefix") {
		batchCfg.Prefix = viper.GetString("batch.prefix")
	}
	if !cmd.Flags().Changed("quiet") && viper.IsSet("batch.quiet") {
		batchCfg.Quiet = viper.GetBool("batch.quiet")
	}
}

func bindWifiFlags(cmd *cobra.Command) {
	bindFlag(cmd, "wifi.ssid", "ssid")
	bindFlag(cmd, "wifi.pass", "pass")
	bindFlag(cmd, "wifi.security", "security")
	bindFlag(cmd, "wifi.hidden", "hidden")
}

func bindVCardFlags(cmd *cobra.Command) {
	bindFlag(cmd, "vcard.name", "name")
	bindFlag(cmd, "vcard.phone", "phone")
	bindFlag(cmd, "vcard.email", "email")
	bindFlag(cmd, "vcard.org", "org")
	bindFlag(cmd, "vcard.title", "title")
	bindFlag(cmd, "vcard.url", "url")
	bindFlag(cmd, "vcard.address", "address")
}

func bindBatchFlags(cmd *cobra.Command) {
	bindFlag(cmd, "batch.file", "file")
	bindFlag(cmd, "batch.dir", "dir")
	bindFlag(cmd, "batch.size", "size")
	bindFlag(cmd, "batch.format", "format")
	bindFlag(cmd, "batch.prefix", "prefix")
	bindFlag(cmd, "batch.quiet", "quiet")
}
