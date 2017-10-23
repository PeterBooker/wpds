package main

import (
	"github.com/go-chi/valve"
	"github.com/peterbooker/wpds/internal/http"
	"github.com/peterbooker/wpds/internal/worker"
)

func main() {

	// Setup Context
	valv := valve.New()
	baseCtx := valv.Context()

	// Start background Workers
	worker.Init(baseCtx)

	// Start the HTTP Server
	http.Init(baseCtx, valv)

}
