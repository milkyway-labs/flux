package start

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/milkyway-labs/chain-indexer/cli/types"
	"github.com/milkyway-labs/chain-indexer/prometheus"
)

func NewStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start parsing the blockchain data",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := types.GetCliContext(cmd)
			return startParsing(cmd.Context(), ctx)
		},
	}
}

func startParsing(ctx context.Context, cliCtx *types.CliContext) error {
	cfg, err := cliCtx.LoadConfig()
	if err != nil {
		return err
	}

	// Start the monitoring server
	prometheusServer := prometheus.NewServer(&cfg.Monitoring)
	prometheusServer.Start()

	indexers, err := cliCtx.IndexersBuilder.BuildAll(ctx, cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(ctx)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	waitGroup := sync.WaitGroup{}
	for _, indexer := range indexers {
		err := indexer.Start(ctx, &waitGroup)
		if err != nil {
			return fmt.Errorf("staring indexer %s, %w", indexer.GetName(), err)
		}
	}

	waitGroup.Wait()
	prometheusServer.Stop()

	return nil
}
