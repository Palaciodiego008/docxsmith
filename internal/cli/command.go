package cli

import (
	"flag"
)

// Command represents a CLI command
type Command interface {
	// Name returns the command name
	Name() string

	// Description returns a short description
	Description() string

	// Usage returns detailed usage information
	Usage() string

	// Execute runs the command with the given arguments
	Execute(args []string) error
}

// CommandRegistry manages all available commands
type CommandRegistry struct {
	commands map[string]Command
}

// NewCommandRegistry creates a new command registry
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]Command),
	}
}

// Register registers a new command
func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

// Get retrieves a command by name
func (r *CommandRegistry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// All returns all registered commands
func (r *CommandRegistry) All() map[string]Command {
	return r.commands
}

// BaseCommand provides common functionality for commands
type BaseCommand struct {
	name        string
	description string
	usage       string
	flagSet     *flag.FlagSet
}

// NewBaseCommand creates a new base command
func NewBaseCommand(name, description, usage string) *BaseCommand {
	return &BaseCommand{
		name:        name,
		description: description,
		usage:       usage,
		flagSet:     flag.NewFlagSet(name, flag.ExitOnError),
	}
}

// Name returns the command name
func (b *BaseCommand) Name() string {
	return b.name
}

// Description returns the command description
func (b *BaseCommand) Description() string {
	return b.description
}

// Usage returns the usage information
func (b *BaseCommand) Usage() string {
	return b.usage
}

// FlagSet returns the flag set
func (b *BaseCommand) FlagSet() *flag.FlagSet {
	return b.flagSet
}
