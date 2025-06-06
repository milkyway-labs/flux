package types

import "github.com/spf13/cobra"

type BeforeStartHook func(*cobra.Command, *CliContext) error
