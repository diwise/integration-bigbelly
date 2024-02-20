package application

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diwise/integration-bigbelly/internal/pkg/domain"
	"github.com/matryer/is"
)

func TestGetAssets(t *testing.T) {
	is := is.New(t)

	var assets []domain.Asset

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, assetsResponse)
	}))
	defer svr.Close()

	app := App{
		bigBellyApiUrl: svr.URL,
		xToken:         "nisse",
	}

	assets, err := app.GetAssets(context.Background())
	is.NoErr(err)

	is.Equal(3, len(assets))

}

func TestMapToFillingLevels(t *testing.T) {
	is := is.New(t)

	r := BigBellyResponse{}
	err := json.Unmarshal([]byte(assetsResponse), &r)
	is.NoErr(err)

	app := App{}
	fillingLevels, err := app.MapToFillingLevels(context.Background(), r.Assets)
	is.NoErr(err)

	is.Equal(3, len(fillingLevels))
	fmt.Println(len(fillingLevels))

	is.Equal(float64(60), fillingLevels[1].ActualFillingPercentage)
}

func TestSend(t *testing.T) {
	is := is.New(t)

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, assetsResponse)
	}))
	defer svr.Close()

	app := App{
		diwiseApiUrl: svr.URL,
		xToken:       "nisse",
	}
	r := BigBellyResponse{}
	err := json.Unmarshal([]byte(assetsResponse), &r)
	is.NoErr(err)

	fillingLevels, err := app.MapToFillingLevels(context.Background(), r.Assets)
	is.NoErr(err)

	i := 0

	testSender := func(context.Context, string, []byte) error {
		i++
		return nil
	}

	err = app.Send(context.Background(), fillingLevels, testSender)
	is.NoErr(err)

	is.Equal(3, i)
	fmt.Println(fillingLevels)
}

const assetsResponse string = `{
    "assets": [
        {
            "latestFullness": 8,
            "reason": "NOT_READY",
            "serialNumber": 1,
            "accountName": "Test kommun",
            "latitude": 61.39204040369734,
            "stationSerialNumber": 1,
            "description": "Torget",
            "ageThreshold": 0,
            "fullnessThreshold": 6,
            "lastCall": 1704949151000,
            "accountId": 1,
            "disposition": "OPERATIONAL",
            "streamType": "TRASH",
            "groupIds": [
                10157
            ],
            "lastCollection": {
                "percentFull": 0,
                "timestamp": 1704872214000
            },
            "position": "center",
            "longitude": 18.30680549824277,
            "status": "IN_SERVICE"
        },
        {
            "latestFullness": 6,
            "reason": "NOT_READY",
            "serialNumber": 2,
            "accountName": "Test kommun",
            "latitude": 61.38743005399236,
            "stationSerialNumber": 2,
            "description": "Skolan",
            "ageThreshold": 0,
            "fullnessThreshold": 6,
            "lastCall": 1704952202000,
            "accountId": 1,
            "disposition": "OPERATIONAL",
            "streamType": "TRASH",
            "groupIds": [
                10144
            ],
            "lastCollection": {
                "percentFull": 8,
                "timestamp": 1700644191000
            },
            "position": "center",
            "longitude": 18.297631047783796,
            "status": "IN_SERVICE"
        },
        {
            "latestFullness": 0,
            "reason": "NOT_READY",
            "serialNumber": 3,
            "accountName": "Test kommun",
            "latitude": 61.3930171865015,
            "stationSerialNumber": 3,
            "description": "Parken",
            "ageThreshold": 0,
            "fullnessThreshold": 6,
            "lastCall": 1704945847000,
            "accountId": 1,
            "disposition": "OPERATIONAL",
            "streamType": "TRASH",
            "groupIds": [
                10144
            ],
            "lastCollection": {
                "percentFull": 8,
                "timestamp": 1704700513000
            },
            "position": "center",
            "longitude": 18.308756448328495,
            "status": "IN_SERVICE"
        }
    ],
    "errorCode": "STATUS_OK"
}`
