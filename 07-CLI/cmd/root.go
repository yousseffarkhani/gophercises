package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "task",                       // Name of the command
	Short: "Task is a CLI task manager", // Description
}
