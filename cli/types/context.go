package types

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	database "github.com/milkyway-labs/flux/database/manager"
	indexerbuilder "github.com/milkyway-labs/flux/indexer/builder"
	modulesmanager "github.com/milkyway-labs/flux/modules/manager"
	nodemanager "github.com/milkyway-labs/flux/node/manager"
	"github.com/milkyway-labs/flux/types"
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
	IndexersBuilder  *indexerbuilder.IndexersBuilder
	// Function that is called before the start cmd is executed.
	BeforeStartHook BeforeStartHook
	// Function that is called after the configurations have been loaded from the
	// disk.
	RawConfigLoadedHook RawConfigLoadedHook
}

func NewCliContext(
	name string,
) *CliContext {
	databaseManager := database.NewDatabasesManager()
	nodeManager := nodemanager.NewNodesManager()
	modulesManager := modulesmanager.NewModuleManager()

	return &CliContext{
		Name:             name,
		DatabasesManager: databaseManager,
		NodesManager:     nodeManager,
		ModulesManager:   modulesManager,
		IndexersBuilder: indexerbuilder.NewIndexersBuilder(
			databaseManager, nodeManager, modulesManager,
		),
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

func (c *CliContext) LoadConfigFileContent() ([]byte, error) {
	configFilePath := c.GetConfigFilePath()

	// Make sure the path exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist (%s)", configFilePath)
	}

	// Read the file content
	configData, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	if c.RawConfigLoadedHook != nil {
		err := c.RawConfigLoadedHook(c, configData)
		if err != nil {
			return nil, fmt.Errorf("raw configuration loaded hook: %w", err)
		}
	}

	return configData, nil
}

func (c *CliContext) LoadConfig() (*types.Config, error) {
	// Read the config file
	configFileContent, err := c.LoadConfigFileContent()
	if err != nil {
		return nil, err
	}

	// Parse the config
	config, err := types.ParseConfig(configFileContent)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *CliContext) WithBeforeStartHook(hook BeforeStartHook) *CliContext {
	c.BeforeStartHook = hook
	return c
}

func (c *CliContext) WithRawConfigLoadedHook(hook RawConfigLoadedHook) *CliContext {
	c.RawConfigLoadedHook = hook
	return c
}

func InjectCliContext(ctx context.Context, cliCtx *CliContext) context.Context {
	return context.WithValue(ctx, ContextKey, cliCtx)
}

func GetCliContext(cmd *cobra.Command) *CliContext {
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
		panic("no CliContext found, please inject it with the InjectCliContext function")
	}

	// Set the context home path from the cmd flag
	homePath, err := cmd.Flags().GetString(FlagHome)
	if err != nil {
		panic(fmt.Sprintf("can't get context from cmd, cmd don't have the %s flag", FlagHome))
	}
	ctx.SetCfgPath(homePath)

	return ctx
}
