package manager

import (
	"context"
	"fmt"

	"github.com/milkyway-labs/flux/database"
	"github.com/milkyway-labs/flux/modules"
	"github.com/milkyway-labs/flux/node"
)

// ModulesManager represents a component that is capable of constructing
// indexing modules and can be used to register custom user's defined modules.
type ModulesManager struct {
	registered map[string]Builder
}

func NewModuleManager() *ModulesManager {
	return &ModulesManager{
		registered: make(map[string]Builder),
	}
}

// RegisterModule register a new indexing module.
func (mm *ModulesManager) RegisterModule(moduleName string, builder Builder) *ModulesManager {
	mm.registered[moduleName] = builder
	return mm
}

// GetModule builds an return the module with the requested name.
func (mm *ModulesManager) GetModule(
	ctx context.Context,
	moduleName string,
	db database.Database,
	node node.Node,
	cfg []byte,
) (modules.Module, error) {
	// Get the module builder config
	moduleBuilder, foundModuleBuilder := mm.registered[moduleName]
	if !foundModuleBuilder {
		return nil, fmt.Errorf("module `%s` not registered", moduleName)
	}

	// Build the module
	return moduleBuilder(ctx, db, node, cfg)
}
