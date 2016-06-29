package main

import (
	"ghenga/db"

	"github.com/jessevdk/go-flags"
	"github.com/jmoiron/modl"
)

type globalOptions struct {
	DB    string `short:"d" long:"database" default:"" env:"GHENGA_DB"   description:"Connection string for postgresql database" default:"host=/var/run/postgresql"`
	Debug bool   `short:"D" long:"debug"                            env:"GHENGA_DEBUG" description:"Enable debug messages for development"`
}

// OpenDB opens the database, which is initialized if necessary. Before exit,
// cleanup() should be called to properly close the database connection.
func OpenDB() (dbm *modl.DbMap, cleanup func() error, err error) {
	dbm, err = db.Init(globalOpts.DB)
	if err != nil {
		return nil, nil, err
	}

	err = db.Migrate(dbm)
	if err != nil {
		return nil, nil, err
	}

	cleanup = func() error { return dbm.Db.Close() }
	return dbm, cleanup, nil
}

var globalOpts = globalOptions{}
var parser = flags.NewParser(&globalOpts, flags.HelpFlag|flags.PassDoubleDash)
