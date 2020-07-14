package parser

import (
	"fmt"

	docopt "github.com/docopt/docopt-go"
	"github.com/teamhephy/workflow-cli/cmd"
	"github.com/teamhephy/workflow-cli/executable"
)

// Git routes git commands to their specific function.
func Git(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Valid commands for git:

git:remote          Adds git remote of application to repository
git:remove          Removes git remote of application from repository

Use '{{.Name}} help [command]' to learn more.
`)

	switch argv[0] {
	case "git:remote":
		return gitRemote(argv, cmdr)
	case "git:remove":
		return gitRemove(argv, cmdr)
	case "git":
		fmt.Print(usage)
		return nil
	default:
		PrintUsage(cmdr)
		return nil
	}
}

func gitRemote(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Adds git remote of application to repository

Usage: {{.Name}} git:remote [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -r --remote=REMOTE
    name of remote to create. [default: {{.Remote}}]
  -f --force
    overwrite remote of the given name if it already exists.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	remote := safeGetValue(args, "--remote")
	force := args["--force"].(bool)

	return cmdr.GitRemote(app, remote, force)
}

func gitRemove(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Removes git remotes of application from repository.

Usage: {{.Name}} git:remove [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	return cmdr.GitRemove(safeGetValue(args, "--app"))
}
