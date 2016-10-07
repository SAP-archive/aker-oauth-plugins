package main

import (
	"net/http"

	"github.com/SAP/aker/plugin"
	"github.com/SAP/gologger"

	"github.com/SAP/aker-oauth-plugins/handler"
)

func main() {
	factory := func(data []byte) (http.Handler, error) {
		cfg, err := handler.ParseConfig(data)
		if err != nil {
			return nil, err
		}
		return handler.CallbackHandlerFromConfig(cfg)
	}
	if err := plugin.ListenAndServeHTTP(factory); err != nil {
		gologger.Fatalf("Error creating plugin: %v", err)
	}
}
