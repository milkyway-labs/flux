package types

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	database "github.com/milkyway-labs/chain-indexer/database/manager"
	modulesmanager "github.com/milkyway-labs/chain-indexer/modules/manager"
	nodemanager "github.com/milkyway-labs/chain-indexer/node/manager"
	"github.com/milkyway-labs/chain-indexer/types"
)

type CliContextKey string

const ContextKey = CliContextKey("cli.context")

// CliContext represents the context that is passed to all the root's sob commands
type CliContext struct {
	Name             string
	CfgDir           string
	DatabasesManager *database.DatabasesManager
	NodesManager     *nodemanager.NodesManager
	ModulesManager   *modulesmanager.ModulesManager
}

func NewCliContext(
	name string,
) *CliContext {
	return &CliContext{
		Name:             name,
		DatabasesManager: database.NewDatabasesManager(),
		NodesManager:     nodemanager.NewNodesManager(),
		ModulesManager:   modulesmanager.NewModuleManager(),
	}
}

func (c *CliContext) GetName() string {
	return c.Name
}

func (c *CliContext) SetCfgPath(path string) {
	c.CfgDir = path
}

func (c *CliContext) GetConfigFilePath() string {
	if c.CfgDir == "" {
		panic("Can't get config file path, config path is not set")
	}

	return path.Join(c.CfgDir, "config.yaml")
}

func (c *CliContext) LoadConfig() (*types.Config, error) {
	configFilePath := c.GetConfigFilePath()

	// Make sure the path exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist (%s)", configFilePath)
	}

	// Read the config
	config, err := types.ParseConfig(configFilePath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func InjectCmdContext(ctx context.Context, cmdCtx *CliContext) context.Context {
	return context.WithValue(ctx, ContextKey, cmdCtx)
}

func GetCmdContext(cmd *cobra.Command) *CliContext {
	var ctx *CliContext
	currCmd := cmd
	for {
		ctxValue, ok := currCmd.Context().Value(ContextKey).(*CliContext)
		if !ok {
			currCmd = currCmd.Parent()
			// No more parents
			if currCmd == nil {
				break
			}
		} else {
			ctx = ctxValue
			break
		}
	}
	if ctx == nil {
		panic("no cmd context found, please inject it with the InjectCmdContext function")
	}

	// Set the context home path from the cmd flag
	homePath, err := cmd.Flags().GetString(FlagHome)
	if err != nil {
		panic(fmt.Sprintf("can't get context from cmd, cmd don't have the %s flag", FlagHome))
	}
	ctx.SetCfgPath(homePath)

	return ctx
}
