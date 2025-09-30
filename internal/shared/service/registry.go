package service

import (
	"context"
	"fmt"
	"log/slog"
	"slices"

	"github.com/aby-med/medical-platform/internal/shared/config"
	"github.com/go-chi/chi/v5"
)

const (
	// AllModulesWildcard is used to indicate all modules should be enabled
	AllModulesWildcard = "*"
)

// Module defines the interface that all service modules must implement
type Module interface {
	// Name returns the unique identifier of the module
	Name() string

	// Initialize sets up the module's dependencies (database, handlers, etc.)
	// Must be called before MountRoutes
	Initialize(ctx context.Context) error

	// MountRoutes registers the module's HTTP routes on the provided router
	MountRoutes(r chi.Router)

	// Start initializes and starts the module's background processes
	// It should block until the context is canceled or an error occurs
	Start(ctx context.Context) error
}

// Registry manages the collection of available modules and their lifecycle
type Registry struct {
	modules []Module
	logger  *slog.Logger
	config  *config.Config
}

// NewRegistry creates a new module registry with the provided configuration
func NewRegistry(cfg *config.Config, logger *slog.Logger) *Registry {
	return &Registry{
		modules: []Module{},
		logger:  logger,
		config:  cfg,
	}
}

// Register adds a module to the registry
func (r *Registry) Register(module Module) {
	r.modules = append(r.modules, module)
	r.logger.Info("Module registered", slog.String("module", module.Name()))
}

// GetModules returns all modules that match the enabled list
// If the enabled list contains the wildcard "*", all modules are returned
// If a module name in the enabled list doesn't exist, an error is returned
func (r *Registry) GetModules(enabled []string) ([]Module, error) {
	// If wildcard is specified, return all modules
	if len(enabled) == 1 && enabled[0] == AllModulesWildcard {
		r.logger.Info("All modules enabled via wildcard")
		return r.modules, nil
	}

	// Filter modules based on the enabled list
	var result []Module
	var unknownModules []string

	for _, name := range enabled {
		found := false
		for _, m := range r.modules {
			if m.Name() == name {
				result = append(result, m)
				found = true
				break
			}
		}
		if !found {
			unknownModules = append(unknownModules, name)
		}
	}

	// Report any unknown modules
	if len(unknownModules) > 0 {
		return nil, fmt.Errorf("unknown modules requested: %v", unknownModules)
	}

	return result, nil
}

// GetModuleByName returns a specific module by name
func (r *Registry) GetModuleByName(name string) (Module, error) {
	for _, m := range r.modules {
		if m.Name() == name {
			return m, nil
		}
	}
	return nil, fmt.Errorf("module not found: %s", name)
}

// AllModuleNames returns the names of all registered modules
func (r *Registry) AllModuleNames() []string {
	names := make([]string, 0, len(r.modules))
	for _, m := range r.modules {
		names = append(names, m.Name())
	}
	return names
}

// HasModule checks if a module with the given name is registered
func (r *Registry) HasModule(name string) bool {
	return slices.ContainsFunc(r.modules, func(m Module) bool {
		return m.Name() == name
	})
}

// Count returns the number of registered modules
func (r *Registry) Count() int {
	return len(r.modules)
}

// MountAllRoutes mounts routes for all modules in the registry
func (r *Registry) MountAllRoutes(router chi.Router) {
	for _, m := range r.modules {
		m.MountRoutes(router)
	}
}

// StartAll starts all modules in the registry
func (r *Registry) StartAll(ctx context.Context) []error {
	var errs []error
	for _, m := range r.modules {
		if err := m.Start(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to start module %s: %w", m.Name(), err))
		}
	}
	return errs
}
