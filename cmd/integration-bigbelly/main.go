package main

import (
	"context"
	"fmt"
	"os"

	"github.com/diwise/integration-bigbelly/internal/pkg/application"
	"github.com/diwise/service-chassis/pkg/infrastructure/buildinfo"
	"github.com/diwise/service-chassis/pkg/infrastructure/env"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y"
)

const (
	serviceName string = "integration-bigbelly"
)

func main() {
	serviceVersion := buildinfo.SourceVersion()

	ctx, log, cleanup := o11y.Init(context.Background(), serviceName, serviceVersion)
	defer cleanup()

	bigBellyApiUrl := env.GetVariableOrDie(ctx, "BIGBELLY_API", "bigbelly url")
	xToken := env.GetVariableOrDie(ctx, "XTOKEN", "API key")

	app := application.New(bigBellyApiUrl, xToken)

	assets, err := app.GetAssets(ctx)
	if err != nil {
		// felhantera...
		log.Error("failed to get assets", "err", err.Error())
		os.Exit(1)
	}
	fmt.Println(assets)
}
