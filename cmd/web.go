package cmd

import (
	"fmt"
	"os"

	"subdomain-finder/internal/web"

	"github.com/spf13/cobra"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start web interface",
	Long:  `Start the web interface for subdomain enumeration`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")

		server := web.NewWebServer(port)
		if err := server.Start(); err != nil {
			fmt.Printf("Error starting web server: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().IntP("port", "p", 8080, "Port to run web interface on")
}
