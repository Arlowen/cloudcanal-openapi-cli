package main

import (
	"os"

	"cloudcanal-openapi-cli/internal/app"
	"cloudcanal-openapi-cli/internal/config"
	"cloudcanal-openapi-cli/internal/console"
	"cloudcanal-openapi-cli/internal/i18n"
	"cloudcanal-openapi-cli/internal/repl"
	"cloudcanal-openapi-cli/internal/util"
)

func main() {
	io := console.NewStdIO(os.Stdin, os.Stdout)
	runtime := app.NewRuntime(config.NewService(""))
	ok, err := runtime.InitializeIfNeeded(io)
	if err != nil {
		io.Println(i18n.T("common.fatalErrorPrefix", err.Error()))
		os.Exit(1)
	}
	if !ok {
		return
	}

	shell := repl.NewShell(io, runtime)
	if len(os.Args) > 1 {
		if err := shell.ExecuteArgs(os.Args[1:]); err != nil {
			io.Println(i18n.T("common.fatalErrorPrefix", util.SummarizeError(err)))
			os.Exit(1)
		}
		return
	}

	if err := shell.Run(); err != nil {
		io.Println(i18n.T("common.fatalErrorPrefix", util.SummarizeError(err)))
		os.Exit(1)
	}
}
