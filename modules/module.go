package modules

// Module represent a module used to index a block chain.
type Module interface {
	// GetName gets the name that identifies the module
	GetName() string
}
