package app

import "github.com/urfave/cli"

const Version = "1.0.0-beta.1"
const Name = "local-ci"

func CreateApp() *cli.App {
	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Usage = "Run your GitLab CI pipelines from your machine using Docker"

	return app
}
