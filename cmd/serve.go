package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yellowhama/musu-crawl-ai/internal/web"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web dashboard to browse knowledge",
	Run: func(cmd *cobra.Command, args []string) {
		out, _ := cmd.Flags().GetString("out")
		port, _ := cmd.Flags().GetInt("port")

		server := web.NewServer(out, port)
		server.Start()
	},
}

func init() {
	serveCmd.Flags().String("out", "./wiki", "Wiki directory to serve")
	serveCmd.Flags().Int("port", 8080, "Port to run the web server on")
	rootCmd.AddCommand(serveCmd)
}
