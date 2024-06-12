package main

import (
	"slices"

	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/config"
	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/server"
	"github.com/urfave/cli/v2"
)

const (
	categoryServer = "SERVER:"
)

func CommandServe(cfg *config.Config) *cli.Command {
	serverFlags := []cli.Flag{
		&cli.StringFlag{
			Category:    categoryServer,
			Destination: &cfg.Server.ListenAddress,
			EnvVars:     []string{envPrefix + "LISTEN_ADDRESS"},
			Name:        "listen-address",
			Usage:       "`host:port` for the server to listen on",
			Value:       "0.0.0.0:8080",
		},

		&cli.StringFlag{
			Category:    categoryServer,
			Destination: &cfg.Server.Name,
			EnvVars:     []string{envPrefix + "SERVER_NAME"},
			Name:        "server-name",
			Usage:       "service `name` to report in prometheus metrics",
			Value:       "{{ cookiecutter.project_slug }}",
		},
	}

	flags := slices.Concat(
		serverFlags,
	)

	return &cli.Command{
		Name:  "serve",
		Usage: "run {{ cookiecutter.project_slug }} server",
		Flags: flags,

		Before: func(ctx *cli.Context) error {
			// TODO: validate inputs
			return nil
		},

		Action: func(_ *cli.Context) error {
			s, err := server.New(cfg)
			if err != nil {
				return err
			}
			return s.Run()
		},
	}
}
