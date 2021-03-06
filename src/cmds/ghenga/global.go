package main

import (
	"ghenga/db"

	"github.com/jessevdk/go-flags"
)

type globalOptions struct {
	DB    string `short:"d" long:"database" default:"" env:"GHENGA_DB"   description:"Connection string for postgresql database" default:"host=/var/run/postgresql"`
	Debug bool   `short:"D" long:"debug"                            env:"GHENGA_DEBUG" description:"Enable debug messages for development"`
}

// OpenDB opens the database, which is initialized if necessary. Before exit,
// cleanup() should be called to properly close the database connection.
func OpenDB() (*db.DB, error) {
	return db.Init(globalOpts.DB)
}

var globalOpts = globalOptions{}
var parser = flags.NewParser(&globalOpts, flags.HelpFlag|flags.PassDoubleDash)
