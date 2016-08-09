package main

import (
	"net/http"

	"github.infra.hana.ondemand.com/cloudfoundry/aker/plugin"
	"github.infra.hana.ondemand.com/cloudfoundry/gologger"

	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugins/handler"
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
