package main

import (
	"context"

	"github.com/diwise/service-chassis/pkg/infrastructure/buildinfo"
	"github.com/diwise/service-chassis/pkg/infrastructure/env"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y"
)

const (
	serviceName string = "integration-bigbelly"
)

func main() {
	serviceVersion := buildinfo.SourceVersion()

	ctx, _, cleanup := o11y.Init(context.Background(), serviceName, serviceVersion)
	defer cleanup()

	env.GetVariableOrDie(ctx, "BIGBELLY_API", "bigbelly url")

}
