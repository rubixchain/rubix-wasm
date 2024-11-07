package wasmbridge

// HostFunctionRegistry manages the registration of host functions.
type HostFunctionRegistry struct {
	hostFunctions []HostFunction
}

// NewHostFunctionRegistry creates a new registry with predefined host functions.
func NewHostFunctionRegistry() *HostFunctionRegistry {
	registry := &HostFunctionRegistry{
		hostFunctions: []HostFunction{},
	}

	// Register predefined host functions
	registry.Register(NewDoApiCall())
	registry.Register(NewDoMintNFTApiCall())

	return registry
}

// Register adds a new host function to the registry.
func (r *HostFunctionRegistry) Register(hf HostFunction) {
	r.hostFunctions = append(r.hostFunctions, hf)
}

// GetHostFunctions returns all registered host function names.
func (r *HostFunctionRegistry) GetHostFunctions() []HostFunction {
	return r.hostFunctions
}
