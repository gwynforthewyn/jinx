package cmd

import (
	"github.com/spf13/cobra"
	"jinx/src/jinkiesengine"
	jinxtypes "jinx/types"
)

type ServeRuntime struct {
	GlobalRuntime jinxtypes.JinxData

	ContainerConfigPath string
	HostConfigPath      string
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Subcommands to allow you to start or stop a jinkies.",
	Long: `Why would you want an unconfigured instance of jinkies? Any time you want a jenkins instance
quickly for reasons unrelated to a specific job. Maybe you want to prototype some jcasc settings or something.

Maybe you want two instances of jinkies running at once? Use the -o flag to supply a yaml file overriding the hostconfig
(https://pkg.go.dev/github.com/docker/docker@v20.10.13+incompatible/api/types/container#HostConfig),
or use -c to supply a yaml file overriding the container config (https://pkg.go.dev/github.com/docker/docker@v20.10.13+incompatible/api/types/container#Config).
`,
}

func (server *ServeRuntime) startSubCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start jinkies!",
		Long:  `Starts the unconfigured jinkies container`,
		Run: func(cmd *cobra.Command, args []string) {
			jinkiesengine.RunRunRun(server.GlobalRuntime.ContainerName, server.GlobalRuntime.PullImages, server.ContainerConfigPath, server.HostConfigPath)
		},
	}
}

func (server *ServeRuntime) stopSubCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stops your jinkies container_info.",
		Long:  `No configuration is retained after a stop, so this gets you back to a clean slate.`,
		Run: func(cmd *cobra.Command, args []string) {
			jinkiesengine.StopGirl(server.GlobalRuntime.ContainerName)
		},
	}
}

func RegisterServe(jinxRunTime jinxtypes.JinxData) {
	config := ServeRuntime{GlobalRuntime: jinxRunTime}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(serveCmd)
	serveCmd.AddCommand(config.startSubCommand())
	serveCmd.AddCommand(config.stopSubCommand())

	serveCmd.PersistentFlags().StringVarP(&config.ContainerConfigPath, "containerconfig", "c", "", "Path to config file describing your container")
	serveCmd.PersistentFlags().StringVarP(&config.HostConfigPath, "hostconfig", "o", "", "Path to config file describing your container host ")
}
