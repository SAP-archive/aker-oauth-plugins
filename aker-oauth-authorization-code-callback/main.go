package main

import (
	"net/http"

	"github.infra.hana.ondemand.com/I061150/aker/logging"
	"github.infra.hana.ondemand.com/I061150/aker/plugin"

	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/handler"
	"github.infra.hana.ondemand.com/cloudfoundry/aker-oauth-plugin/handler/callback"
)

func main() {
	factory := func(data []byte) (http.Handler, error) {
		cfg, err := handler.ParseConfig(data)
		if err != nil {
			return nil, err
		}
		return callback.HandlerFromConfig(cfg)
	}
	if err := plugin.ListenAndServeHTTP(factory); err != nil {
		logging.Fatalf("Error creating plugin: %v", err)
	}
}
