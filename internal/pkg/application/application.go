package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/diwise/integration-bigbelly/internal/pkg/domain"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("integration-bigbelly/app")

type App struct {
	bigBellyApiUrl string
	xToken         string
	diwiseApiUrl   string
}

type BigBellyResponse struct {
	Assets    []domain.Asset `json:"assets"`
	ErrorCode string         `json:"errorCode"`
}

const OutOfService string = "OUT_OF_SERVICE"

func New(bigBellyApiUrl string, xToken string, diwiseApiUrl string) App {
	return App{
		bigBellyApiUrl: bigBellyApiUrl,
		xToken:         xToken,
		diwiseApiUrl:   diwiseApiUrl,
	}
}

func (a *App) GetAssets(ctx context.Context) ([]domain.Asset, error) {
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %s", err.Error())
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("X-Token", a.xToken)
	
	resp, err := httpClient.Do(req)
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

	return bigBellyResponse.Assets, nil
}

func (a *App) MapToFillingLevels(ctx context.Context, assets []domain.Asset) ([]domain.FillingLevel, error) {

	var fillingLevels []domain.FillingLevel

	for _, asset := range assets {
		if asset.Status == OutOfService {
			continue
		} else {

			containerFull := false
			containerEmpty := true

			actualFillingPercentage := float64(asset.LatestFullness) * 10
			highThreshold := float64(asset.FullnessThreshold) * 10

			if actualFillingPercentage >= highThreshold {
				containerFull = true
			}
			if actualFillingPercentage > 0 {
				containerEmpty = false
			}

			fl := domain.NewFillingLevel(strconv.Itoa(int(asset.SerialNumber)), actualFillingPercentage, containerFull, containerEmpty, time.Now().UTC())
			fl.HighThreshold = highThreshold

			fillingLevels = append(fillingLevels, fl)
		}
	}

	return fillingLevels, nil
}

func (a *App) Send(ctx context.Context, fillingLevels []domain.FillingLevel, sender SenderFunc) error {
	for _, f := range fillingLevels {
		b, err := json.Marshal(f)
		if err != nil {
			return err
		}

		err = sender(ctx, a.diwiseApiUrl, b)
		if err != nil {
			return err
		}
	}

	return nil
}

type SenderFunc = func(context.Context, string, []byte) error

func HttpPost(ctx context.Context, url string, b []byte) error {
	var err error

	ctx, span := tracer.Start(ctx, "send-fillingLevel")
	defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()

	httpClient := http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/senml+json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	} else if resp.StatusCode != http.StatusCreated {
		err = fmt.Errorf("unexpected response code %d", resp.StatusCode)
	}

	return err
}
