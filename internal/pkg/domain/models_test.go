package domain

import (
	"encoding/json"
	"testing"

	"github.com/matryer/is"
)

func TestAssetStruct(t *testing.T) {
	is := is.New(t)
	var a Asset

	err := json.Unmarshal([]byte(assetJson), &a)
	is.NoErr(err)

	is.Equal("", a.Description)
	is.Equal(1704872214000, a.LastCollection.Timestamp)
}

const assetJson string = `
{
	"latestFullness": 0,
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
}`
