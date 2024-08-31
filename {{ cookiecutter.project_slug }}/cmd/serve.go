package main

import (
	"slices"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/config"
	"github.com/{{ cookiecutter.github_org }}/{{ cookiecutter.project_slug }}/server"
)

const (
	categoryServer = "server"
)

func CommandServe(cfg *config.Config) *cli.Command {
	serverFlags := []cli.Flag{
		&cli.StringFlag{
			Category:    strings.ToUpper(categoryServer),
			Destination: &cfg.Server.ListenAddress,
			EnvVars:     []string{envPrefix + strings.ToUpper(categoryServer) + "_LISTEN_ADDRESS"},
			Name:        categoryServer + "-listen-address",
			Usage:       "`host:port` for the server to listen on",
			Value:       "0.0.0.0:8080",
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
