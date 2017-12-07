package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
)

var build = "0" // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "downstream plugin"
	app.Usage = "downstream plugin"
	app.Action = run
	app.Version = fmt.Sprintf("1.0.%s", build)
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "repositories",
			Usage:  "List of repositories to trigger",
			EnvVar: "PLUGIN_REPOSITORIES",
		},
		cli.StringFlag{
			Name:   "server",
			Usage:  "Trigger a drone build on a custom server",
			EnvVar: "DOWNSTREAM_SERVER,PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Drone API token from your user settings",
			EnvVar: "DRONE_TOKEN,DOWNSTREAM_TOKEN,PLUGIN_TOKEN",
		},
		cli.BoolFlag{
			Name:   "fork",
			Usage:  "Trigger a new build for a repository",
			EnvVar: "PLUGIN_FORK",
		},
		cli.BoolFlag{
			Name:   "wait",
			Usage:  "Wait for any currently running builds to finish",
			EnvVar: "PLUGIN_WAIT",
		},
		cli.DurationFlag{
			Name:   "timeout",
			Value:  time.Duration(60) * time.Second,
			Usage:  "How long to wait on any currently running builds",
			EnvVar: "PLUGIN_WAIT_TIMEOUT",
		},
		cli.BoolFlag{
			Name:   "track",
			Usage:  "Track triggered build status, error on non-success",
			EnvVar: "PLUGIN_TRACK",
		},
		cli.DurationFlag{
			Name:   "track-interval",
			Value:  time.Duration(60) * time.Second,
			Usage:  "How often to poll build status while tracking builds",
			EnvVar: "PLUGIN_TRACK_INTERVAL",
		},
		cli.DurationFlag{
			Name:   "track-timeout",
			Value:  time.Duration(600) * time.Second,
			Usage:  "How long to wait while tracking builds",
			EnvVar: "PLUGIN_TRACK_TIMEOUT",
		},
		cli.StringSliceFlag{
			Name:   "params",
			Usage:  "List of params (key=value or file paths of params) to pass to triggered builds",
			EnvVar: "PLUGIN_PARAMS",
		},
		cli.StringSliceFlag{
			Name:   "params-from-env",
			Usage:  "List of environment variables to pass to triggered builds",
			EnvVar: "PLUGIN_PARAMS_FROM_ENV",
		},
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repos:         c.StringSlice("repositories"),
		Server:        c.String("server"),
		Token:         c.String("token"),
		Fork:          c.Bool("fork"),
		Wait:          c.Bool("wait"),
		Timeout:       c.Duration("timeout"),
		Params:        c.StringSlice("params"),
		ParamsEnv:     c.StringSlice("params-from-env"),
		Track:         c.Bool("track"),
		TrackInterval: c.Duration("track-interval"),
		TrackTimeout:  c.Duration("track-timeout"),
	}

	return plugin.Exec()
}
