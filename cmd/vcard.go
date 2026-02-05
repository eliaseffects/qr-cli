package cmd

import (
	"github.com/eliaseffects/qr-cli/internal/qr"
	"github.com/spf13/cobra"
)

var (
	vcardFlags   OutputFlags
	vcardName    string
	vcardPhone   string
	vcardEmail   string
	vcardOrg     string
	vcardTitle   string
	vcardURL     string
	vcardAddress string

	vcardCmd = &cobra.Command{
		Use:   "vcard",
		Short: "Generate QR code for contact card",
		RunE:  runVCard,
	}
)

func init() {
	vcardCmd.Flags().StringVar(&vcardName, "name", "", "Full name (required)")
	vcardCmd.Flags().StringVar(&vcardPhone, "phone", "", "Phone number")
	vcardCmd.Flags().StringVar(&vcardEmail, "email", "", "Email address")
	vcardCmd.Flags().StringVar(&vcardOrg, "org", "", "Organization")
	vcardCmd.Flags().StringVar(&vcardTitle, "title", "", "Job title")
	vcardCmd.Flags().StringVar(&vcardURL, "url", "", "Website URL")
	vcardCmd.Flags().StringVar(&vcardAddress, "address", "", "Street address")
	_ = vcardCmd.MarkFlagRequired("name")

	addOutputFlags(vcardCmd, &vcardFlags, true)
	bindOutputFlags(vcardCmd)
	bindVCardFlags(vcardCmd)
}

func runVCard(cmd *cobra.Command, args []string) error {
	applyOutputConfig(cmd, &vcardFlags)
	applyVCardConfig(cmd)

	data := qr.VCard{
		Name:    vcardName,
		Phone:   vcardPhone,
		Email:   vcardEmail,
		Org:     vcardOrg,
		Title:   vcardTitle,
		URL:     vcardURL,
		Address: vcardAddress,
	}.String()

	return runGenerate(data, vcardFlags, cmd.Flags().Changed("format"))
}
