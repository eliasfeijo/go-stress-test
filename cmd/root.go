package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-stress-test",
	Short: "A stress tester CLI like Apache ab",
	Long:  `A stress tester CLI like Apache ab`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-stress-test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("url", "u", "", "The URL to test")
	rootCmd.MarkFlagRequired("url")

	rootCmd.Flags().StringP("method", "m", "GET", "The HTTP method to use")
	rootCmd.Flags().IntP("concurrency", "c", 1, "The number of concurrent requests to make")
	rootCmd.Flags().IntP("requests", "r", 1, "The number of requests to make")
	rootCmd.Flags().IntP("timeout", "t", 10, "The timeout in seconds")
	rootCmd.Flags().BoolP("verbose", "v", false, "Verbose output")
}
