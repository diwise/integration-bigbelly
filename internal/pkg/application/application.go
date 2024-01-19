package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/diwise/integration-bigbelly/internal/pkg/domain"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("integration-bigbelly/app")

type App interface {
	GetAssets(ctx context.Context) ([]domain.Asset, error)
}

type app struct {
	bigBellyApiUrl string
	xToken         string
}

type BigBellyResponse struct {
	Assets    []domain.Asset `json:"assets"`
	ErrorCode string         `json:"errorCode"`
}

func New(bigBellyApiUrl string /*, xToken string*/) App {
	return &app{
		bigBellyApiUrl: bigBellyApiUrl,
		// xToken: xToken,
	}
}

func (a *app) GetAssets(ctx context.Context) ([]domain.Asset, error) {
	var err error

	ctx, span := tracer.Start(ctx, "get-assets")
	defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()

	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	params := url.Values{}
	params.Add("objectType", "assets")
	params.Add("action", "load")

	apiUrl := fmt.Sprintf("%s?%s", a.bigBellyApiUrl, params.Encode())

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %s", err.Error())
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Token", a.xToken)

	var resp *http.Response
	resp, err = httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("request failed: %s", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	var b []byte
	b, err = io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body: %s", err.Error())
		return nil, err
	}

	var bigBellyResponse BigBellyResponse

	err = json.Unmarshal(b, &bigBellyResponse)
	if err != nil {
		err = fmt.Errorf("failed to unmarshal response body: %s", err.Error())
		return nil, err
	}

	// skapa en ny struct för att kunna hantera svaret från servern
	// plocka ut alla assets ur svaret från servern och returnera

	return bigBellyResponse.Assets, nil
}
