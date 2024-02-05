package application

import (
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
}

type BigBellyResponse struct {
	Assets    []domain.Asset `json:"assets"`
	ErrorCode string         `json:"errorCode"`
}

func New(bigBellyApiUrl string, xToken string) App {
	return App{
		bigBellyApiUrl: bigBellyApiUrl,
		xToken:         xToken,
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

	return bigBellyResponse.Assets, nil
}

func (a *App) MapToFillingLevels(ctx context.Context, assets []domain.Asset) ([]domain.FillingLevel, error) {

	var fillingLevels []domain.FillingLevel

	// TODO: implementera logik för att mappa från asset till filling level

	for _, asset := range assets {

		// finns det någon status i asset som gör att vi vill hoppa över en asset?
		// om ja, hoppa över den här asseten och gå vidare till nästa

		if asset.Status == "OUT_OF_SERVICE" {
			continue
		} else {
			//fl := domain.NewFillingLevel(asset.ID, asset.ActualFillingPercentage, asset.ContainerFull, asset.ContainerEmpty, asset.LastCollection.Timestamp)
			// beräkna fillingpercentage mha vad??? LatestFullness *10?
			//actualFillingPercentage := 100;

			// Om fillingPercentage > 100 -> containerFull true annars false
			containerFull := false
			containerEmpty := true
			var actualFillingPercentage float64
			var highThreshold float64

			actualFillingPercentage = float64(asset.LatestFullness) * 10
			highThreshold = float64(asset.FullnessThreshold) * 10

			if actualFillingPercentage >= 100 {
				containerFull = true
			} else {

			}

			// Viken tid ska sättas i time
			fl := domain.NewFillingLevel(strconv.Itoa(int(asset.SerialNumber)), actualFillingPercentage, containerFull, containerEmpty, highThreshold, time.Now().UTC())
			fillingLevels = append(fillingLevels, fl)

		}

		// om nej, lägg till en ny filling level i listan med filling levels

		// är det mer på fl som ska fyllas i än de som är mandatory? gör det i sådana fall
		// fl.AnnanProperty = asset.AnnanProperty

	}

	return fillingLevels, nil
}
