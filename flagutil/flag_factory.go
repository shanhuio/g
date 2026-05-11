package flagutil

import (
	"flag"
)

// FlagFactory is a factory that can be used to create flag sets.
type FlagFactory struct {
	name string
}

// NewFactory creates a new flag factory with the given name.
func NewFactory(name string) *FlagFactory {
	return &FlagFactory{name: name}
}

// ParseArgs parses the args using the given FlagSet, and returns
// the parsed args that does not include flags.
func ParseArgs(s *flag.FlagSet, args []string) []string {
	s.Parse(args)
	return s.Args()
}

// New creates a new flag set.
func (f *FlagFactory) New() *FlagSet {
	set := flag.NewFlagSet(f.name, flag.ExitOnError)
	return &FlagSet{
		FlagSet: set,
	}
}

// PlainArgs parse the args with no flags.
func (f *FlagFactory) PlainArgs(args []string) []string {
	flags := f.New()
	return flags.ParseArgs(args)
}

// FlagSet extends the *flag.FlagSet with some more common helper functions.
type FlagSet struct {
	*flag.FlagSet
}

// ParseArgs parses the args and returns the parsed args that does not include
// flags.
func (s *FlagSet) ParseArgs(args []string) []string {
	return ParseArgs(s.FlagSet, args)
}
