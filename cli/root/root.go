package root

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/milkyway-labs/chain-indexer/cli/types"
	"github.com/spf13/cobra"
)

func NewRootCommad(ctx context.Context, cmdContext *types.CliContext) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   cmdContext.GetName(),
		Short: fmt.Sprintf("%s is a chain data aggregator and exporter", cmdContext.GetName()),
		Long: fmt.Sprintf(`%s is a chain data aggregator. It improves the chain's data accessibility
by providing an indexed database exposing aggregated resources.`, cmdContext.GetName()),
	}

	rootCmd.SetContext(types.InjectCmdContext(ctx, cmdContext))

	// Set the default home path
	home, _ := os.UserHomeDir()
	defaultConfigPath := path.Join(home, fmt.Sprintf(".%s", cmdContext.GetName()))
	rootCmd.PersistentFlags().String(
		types.FlagHome,
		defaultConfigPath,
		"Set the home folder of the application, where all files will be stored",
	)

	return rootCmd
}
