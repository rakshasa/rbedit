package rbeditCmd

import (
	"context"

	"github.com/spf13/cobra"
)

type addCommandFn func(context.Context) (*cobra.Command, context.Context)

const defaultUsageTemplate = `Usage:  {{if .Runnable}}{{.UseLine}}{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

func setupDefaultCommand(cmd *cobra.Command, stateKeyPrefix string) {
	if cmd.Annotations == nil {
		cmd.Annotations = make(map[string]string)
	}

	cmd.Long = "\n" + cmd.Short
	cmd.Annotations[stateKeyPrefixKey] = stateKeyPrefix
	cmd.DisableFlagsInUseLine = true
	cmd.TraverseChildren = true

	cmd.SetUsageTemplate(defaultUsageTemplate)
}

func addCommand(ctx context.Context, rootCmd *cobra.Command, newCmdFn addCommandFn) context.Context {
	cmd, ctx := newCmdFn(ctx)
	rootCmd.AddCommand(cmd)
	return ctx
}

func commandPathAsList(cmd *cobra.Command) []string {
	path := []string{}
	for c := cmd; c.HasParent(); c = c.Parent() {
		path = append([]string{c.Name()}, path...)
	}

	return path
}
