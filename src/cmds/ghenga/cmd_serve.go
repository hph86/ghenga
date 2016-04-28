package main

import (
	"fmt"
	"ghenga/server"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type cmdServe struct {
	Port   uint   `short:"p" long:"port"   default:"8080"   description:"set the port for the HTTP server"`
	Addr   string `short:"b" long:"bind"   default:""       description:"bind to this address"`
	Public string `          long:"public" default:"public" description:"directory for serving static files"`
}

func init() {
	_, err := parser.AddCommand("serve",
		"start server",
		"The server command starts the HTTP server",
		&cmdServe{})
	if err != nil {
		panic(err)
	}
}

func (opts *cmdServe) Execute(args []string) (err error) {
	dbmap, cleanup, e := OpenDB()
	if e != nil {
		return e
	}
	defer CleanupErr(&err, cleanup)

	log.Printf("starting server at %v:%d", opts.Addr, opts.Port)

	env := &server.Env{
		ListenAddr: fmt.Sprintf("%s:%d", opts.Addr, opts.Port),
		DbMap:      dbmap,
		Public:     opts.Public,
		Debug:      globalOpts.Debug,
	}

	router := server.PeopleHandler(env, mux.NewRouter())

	// server static files on the root path
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(env.Public)))

	// activate logging to stdout
	http.Handle("/", handlers.CombinedLoggingHandler(os.Stdout, router))

	return http.ListenAndServe(env.ListenAddr, nil)
}
