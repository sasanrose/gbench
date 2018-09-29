package cmd

var (
	driver, address, port, input string
)

func initRenderFlags() {
	renderCmd.Flags().StringVarP(&input, "input", "i", "./report.json", "Path to the report file.")
	renderCmd.Flags().StringVarP(&driver, "driver", "d", "cli", "Driver to use for rendering the report. Accepted values are 'cli'and 'html'.")
	renderCmd.Flags().StringVarP(&address, "address", "a", "localhost", "Address to access the html report.")
	renderCmd.Flags().StringVarP(&address, "port", "p", "8080", "Port to access the html report.")
}
