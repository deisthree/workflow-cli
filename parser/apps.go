package parser

import (
	"strconv"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"github.com/teamhephy/workflow-cli/cmd"
	"github.com/teamhephy/workflow-cli/executable"
)

// Apps routes app commands to their specific function.
func Apps(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Valid commands for apps:

apps:create        create a new application
apps:list          list accessible applications
apps:info          view info about an application
apps:open          open the application in a browser
apps:logs          view aggregated application logs
apps:run           run a command in an ephemeral app container
apps:destroy       destroy an application
apps:transfer      transfer app ownership to another user

Use '{{.Name}} help [command]' to learn more.
`)

	switch argv[0] {
	case "apps:create":
		return appCreate(argv, cmdr)
	case "apps:list":
		return appsList(argv, cmdr)
	case "apps:info":
		return appInfo(argv, cmdr)
	case "apps:open":
		return appOpen(argv, cmdr)
	case "apps:logs":
		return appLogs(argv, cmdr)
	case "apps:run":
		return appRun(argv, cmdr)
	case "apps:destroy":
		return appDestroy(argv, cmdr)
	case "apps:transfer":
		return appTransfer(argv, cmdr)
	default:
		if printHelp(argv, usage) {
			return nil
		}

		if argv[0] == "apps" {
			argv[0] = "apps:list"
			return appsList(argv, cmdr)
		}

		PrintUsage(cmdr)
		return nil
	}
}

func appCreate(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Creates a new application.

- if no <id> is provided, one will be generated automatically.

Usage: {{.Name}} apps:create [<id>] [options]

Arguments:
  <id>
    a uniquely identifiable name for the application. No other app can already
    exist with this name.

Options:
  --no-remote
    do not create a '{{.Remote}}' git remote.
  -b --buildpack BUILDPACK
    a buildpack url to use for this app
  -r --remote REMOTE
    name of remote to create. [default: {{.Remote}}]
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	id := safeGetValue(args, "<id>")
	buildpack := safeGetValue(args, "--buildpack")
	remote := safeGetValue(args, "--remote")
	noRemote := args["--no-remote"].(bool)

	return cmdr.AppCreate(id, buildpack, remote, noRemote)
}

func appsList(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Lists applications visible to the current user.

Usage: {{.Name}} apps:list [options]

Options:
  -l --limit=<num>
    the maximum number of results to display, defaults to config setting
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	results, err := responseLimit(safeGetValue(args, "--limit"))

	if err != nil {
		return err
	}

	return cmdr.AppsList(results)
}

func appInfo(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Prints info about the current application.

Usage: {{.Name}} apps:info [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	return cmdr.AppInfo(app)
}

func appOpen(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Opens a URL to the application in the default browser.

Usage: {{.Name}} apps:open [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	return cmdr.AppOpen(app)
}

func appLogs(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Retrieves the most recent log events.

Usage: {{.Name}} apps:logs [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  -n --lines=<lines>
    the number of lines to display
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")

	linesStr := safeGetValue(args, "--lines")
	var lines int

	if linesStr == "" {
		lines = -1
	} else {
		lines, err = strconv.Atoi(linesStr)

		if err != nil {
			return err
		}
	}

	return cmdr.AppLogs(app, lines)
}

func appRun(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Runs a command inside an ephemeral app container. Default environment is
/bin/bash.

Usage: {{.Name}} apps:run [options] [--] <command>...

Arguments:
  <command>
    the shell command to run inside the container.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	command := strings.Join(args["<command>"].([]string), " ")

	return cmdr.AppRun(app, command)
}

func appDestroy(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Destroys an application.

Usage: {{.Name}} apps:destroy [options]

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
  --confirm=<app>
    skips the prompt for the application name. <app> is the uniquely identifiable
    name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	confirm := safeGetValue(args, "--confirm")

	return cmdr.AppDestroy(app, confirm)
}

func appTransfer(argv []string, cmdr cmd.Commander) error {
	usage := executable.Render(`
Transfer app ownership to another user.

Usage: {{.Name}} apps:transfer <username> [options]

Arguments:
  <username>
    the user that the app will be transferred to.

Options:
  -a --app=<app>
    the uniquely identifiable name for the application.
`)

	args, err := docopt.Parse(usage, argv, true, "", false, true)

	if err != nil {
		return err
	}

	app := safeGetValue(args, "--app")
	user := safeGetValue(args, "<username>")

	return cmdr.AppTransfer(app, user)
}
