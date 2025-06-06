package types

import "github.com/spf13/cobra"

type BeforeStartHook func(*cobra.Command, *CliContext) error

type RawConfigLoadedHook func(ctx *CliContext, rawConfig []byte) error
