package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [shell]",
	Short: "Generate autocompletion script",
	Long: `Bash:

  $ source <(envManager completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ envManager completion bash > /etc/bash_completion.d/envManager
  # macOS:
  $ envManager completion bash > /usr/local/etc/bash_completion.d/envManager

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, store the autocompletion script
  # somewhere in your $fpath
  envManager completion zsh > /file/in/fpath

  # You will need to start a new shell for this setup to take effect.`,
	ValidArgs: []string{"bash", "zsh"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		}
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
