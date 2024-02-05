package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

type Asset struct {
	LatestFullness    int     `json:"latestFullness"`
	FullnessThreshold int     `json:"fullnessThreshold"`
	Status            string  `json:"status"`
	SerialNumber      int64   `json:"serialNumber"`
	Longitude         float64 `json:"longitude"`
	Latitude          float64 `json:"latitude"`
	Disposition       string  `json:"disposition"`
	Reason            string  `json:"reason"`
	Description       string  `json:"description"`

	LastCollection struct {
		PercentFull int   `json:"percentFull"`
		Timestamp   int64 `json:"timestamp"`
	} `json:"lastCollection"`
}

type Assets struct {
	Assets []Asset `json:"assets"`
}

const prefix = "urn:oma:lwm2m:ext"

type DeviceInfo struct {
	ID_        string    `lwm2m:"-"`
	Timestamp_ time.Time `lwm2m:"-"`
}

func NewFillingLevel(deviceID string, actualFillingPercentage float64, containerFull, containerEmpty bool, highThreshold float64, ts time.Time) FillingLevel {
	return FillingLevel{
		DeviceInfo: DeviceInfo{
			ID_:        deviceID,
			Timestamp_: ts,
		},
		ActualFillingPercentage: actualFillingPercentage,
		ContainerFull:           containerFull,
		ContainerEmpty:          containerEmpty,
		HighThreshold:           highThreshold,
	}
}

type FillingLevel struct {
	DeviceInfo
	ActualFillingPercentage float64 `json:"actualFillingPercentage"`
	ContainerFull           bool    `json:"containerFull"`
	ContainerEmpty          bool    `json:"containerEmpty"`
	HighThreshold           float64 `json:"highThreshold"`
}

type WasteContainer struct {
}

func (fl FillingLevel) ID() string {
	return fl.ID_
}
func (fl FillingLevel) Timestamp() time.Time {
	return fl.Timestamp_
}
func (fl FillingLevel) ObjectID() string {
	return "3435"
}
func (fl FillingLevel) ObjectURN() string {
	return fmt.Sprintf("%s:%s", prefix, fl.ObjectID())
}
func (fl FillingLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(fl)
}
