package builder

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/database/manager"
	"github.com/milkyway-labs/flux/indexer"
	"github.com/milkyway-labs/flux/modules"
	modulesmanager "github.com/milkyway-labs/flux/modules/manager"
	"github.com/milkyway-labs/flux/node"
	nodemanager "github.com/milkyway-labs/flux/node/manager"
	"github.com/milkyway-labs/flux/types"
	"github.com/milkyway-labs/flux/utils"
)

// IndexersBuilder its an object to create the various indexers with they required
// databases, nodes and modules.
type IndexersBuilder struct {
	databasesManager *manager.DatabasesManager
	nodesManager     *nodemanager.NodesManager
	modulesManager   *modulesmanager.ModulesManager
	globalObjects    map[string]any
}

func NewIndexersBuilder(
	databasesManager *manager.DatabasesManager,
	nodesManager *nodemanager.NodesManager,
	modulesManager *modulesmanager.ModulesManager,
) *IndexersBuilder {
	return &IndexersBuilder{
		databasesManager: databasesManager,
		nodesManager:     nodesManager,
		modulesManager:   modulesManager,
	}
}

func (b *IndexersBuilder) BuildAll(ctx context.Context, cfg *types.Config) ([]indexer.Indexer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config can't be nil")
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	logger, err := utils.NewLoggerFromConfig(&cfg.Logging)
	if err != nil {
		return nil, fmt.Errorf("create logger instance: %w", err)
	}

	indexers := make([]indexer.Indexer, len(cfg.Indexers))
	for i, indexerCfg := range cfg.Indexers {
		indexerCtx := types.NewIndexerContext(cfg, &indexerCfg, b.globalObjects, logger)
		ctx = types.InjectIndexerContext(ctx, indexerCtx)

		// Build the indexer's database instance
		indexerDB, err := b.buildDatabase(ctx, cfg, indexerCfg.DatabaseID)
		if err != nil {
			return nil, fmt.Errorf("build database for indexer %s: %w", indexerCfg.Name, err)
		}

		// Build the indexer's node
		indexerNode, err := b.buildNode(ctx, cfg, indexerCfg.NodeID)
		if err != nil {
			return nil, fmt.Errorf("build node for indexer %s: %w", indexerCfg.Name, err)
		}

		// Build the indexer's modules
		indexerModules, err := b.buildModules(ctx, cfg, indexerDB, indexerNode, &indexerCfg)
		if err != nil {
			return nil, fmt.Errorf("build modules for indexer %s: %w", indexerCfg.Name, err)
		}

		// Build the indexer
		indexers[i] = indexer.NewIndexer(&indexerCfg, logger, indexerDB, indexerNode, indexerModules)
	}

	return indexers, nil
}

func (b *IndexersBuilder) BuildByName(ctx context.Context, cfg *types.Config, name string) (indexer.Indexer, error) {
	if cfg == nil {
		return indexer.Indexer{}, fmt.Errorf("config can't be nil")
	}

	if err := cfg.Validate(); err != nil {
		return indexer.Indexer{}, fmt.Errorf("invalid config: %w", err)
	}

	logger, err := utils.NewLoggerFromConfig(&cfg.Logging)
	if err != nil {
		return indexer.Indexer{}, fmt.Errorf("create logger instance: %w", err)
	}

	indexerCfg, err := cfg.GetIndexerConfig(name)
	if err != nil {
		return indexer.Indexer{}, err
	}

	indexerCtx := types.NewIndexerContext(cfg, indexerCfg, b.globalObjects, logger)
	ctx = types.InjectIndexerContext(ctx, indexerCtx)

	// Build the indexer's database instance
	indexerDB, err := b.buildDatabase(ctx, cfg, indexerCfg.DatabaseID)
	if err != nil {
		return indexer.Indexer{}, fmt.Errorf("build database for indexer %s: %w", indexerCfg.Name, err)
	}

	// Build the indexer's node
	indexerNode, err := b.buildNode(ctx, cfg, indexerCfg.NodeID)
	if err != nil {
		return indexer.Indexer{}, fmt.Errorf("build node for indexer %s: %w", indexerCfg.Name, err)
	}

	// Build the indexer's modules
	indexerModules, err := b.buildModules(ctx, cfg, indexerDB, indexerNode, indexerCfg)
	if err != nil {
		return indexer.Indexer{}, fmt.Errorf("build modules for indexer %s: %w", indexerCfg.Name, err)
	}

	// Build the indexer
	return indexer.NewIndexer(indexerCfg, logger, indexerDB, indexerNode, indexerModules), nil
}

// WithGlobalObject adds an object that can be accessed from all the modules
// during their initialization. This can be useful in case you want to share some
// global objects across all the modules.
func (b *IndexersBuilder) WithGlobalObject(key string, value any) *IndexersBuilder {
	b.globalObjects[key] = value
	return b
}

// GetGlobalObject gets one of the global objects given its key.
func (b *IndexersBuilder) GetGlobalObject(key string) any {
	return b.globalObjects[key]
}

func (b *IndexersBuilder) buildDatabase(
	ctx context.Context,
	cfg *types.Config,
	dbID string,
) (database.Database, error) {
	databaseCfg, databaseFound := cfg.Databases[dbID]
	if !databaseFound {
		return nil, fmt.Errorf("database %s not found", dbID)
	}

	dbType, foundDBType := databaseCfg["type"].(string)
	if !foundDBType {
		return nil, fmt.Errorf("can't find 'type' field in database %s", dbID)
	}

	rawConfig, err := yaml.Marshal(databaseCfg)
	if err != nil {
		return nil, fmt.Errorf("marshal database %s config", dbID)
	}

	return b.databasesManager.GetDatabase(ctx, dbType, dbID, rawConfig)
}

func (b *IndexersBuilder) buildNode(ctx context.Context, cfg *types.Config, nodeID string) (node.Node, error) {
	nodeCfg, nodeCfgFound := cfg.Nodes[nodeID]
	if !nodeCfgFound {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	nodeType, foundNodeType := nodeCfg["type"].(string)
	if !foundNodeType {
		return nil, fmt.Errorf("can't find 'type' field in node %s", nodeID)
	}

	rawConfig, err := yaml.Marshal(nodeCfg)
	if err != nil {
		return nil, fmt.Errorf("marshal node %s config", nodeID)
	}

	return b.nodesManager.GetNode(ctx, nodeType, nodeID, rawConfig)
}

func (b *IndexersBuilder) buildModules(
	ctx context.Context,
	cfg *types.Config,
	db database.Database,
	node node.Node,
	indexerCfg *types.IndexerConfig,
) ([]modules.Module, error) {
	modules := make([]modules.Module, len(indexerCfg.Modules))

	for i, moduleName := range indexerCfg.Modules {
		moduleCfg, foundModuleCfg := cfg.Modules[moduleName]
		if !foundModuleCfg {
			moduleCfg = types.RawConfig{}
		}
		// If OverrideModuleConfig is set for this module, override the module config
		// with it
		overrideModuleCfg, foundOverrideModuleCfg := indexerCfg.OverrideModuleConfig[moduleName]
		if foundOverrideModuleCfg {
			utils.CopyMap(moduleCfg, overrideModuleCfg)
		}

		// Convert the module config back to its binary representation
		var rawConfig []byte
		if foundModuleCfg {
			byteConfig, err := yaml.Marshal(moduleCfg)
			if err != nil {
				return nil, fmt.Errorf("marshal module %s config", moduleName)
			}
			rawConfig = byteConfig
		}

		// Build the module
		module, err := b.modulesManager.GetModule(ctx, moduleName, db, node, rawConfig)
		if err != nil {
			return nil, fmt.Errorf("build module `%s` for indexer `%s`: %w", moduleName, indexerCfg.Name, err)
		}
		modules[i] = module
	}

	return modules, nil
}
