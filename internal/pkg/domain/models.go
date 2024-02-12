package domain

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/farshidtz/senml/v2"
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

type Lwm2mObject interface {
	ID() string
	ObjectID() string
	ObjectURN() string
	Timestamp() time.Time
	MarshalJSON() ([]byte, error)
}

type DeviceInfo struct {
	ID_        string    `lwm2m:"-"`
	Timestamp_ time.Time `lwm2m:"-"`
}

func NewFillingLevel(deviceID string, actualFillingPercentage float64, containerFull, containerEmpty bool, ts time.Time) FillingLevel {
	return FillingLevel{
		DeviceInfo: DeviceInfo{
			ID_:        deviceID,
			Timestamp_: ts,
		},
		ActualFillingPercentage: actualFillingPercentage,
		ContainerFull:           containerFull,
		ContainerEmpty:          containerEmpty,
	}
}

type FillingLevel struct {
	DeviceInfo
	ActualFillingPercentage float64 `json:"actualFillingPercentage" lwm2m:"2"`
	ContainerFull           bool    `json:"containerFull" lwm2m:"5"`
	ContainerEmpty          bool    `json:"containerEmpty" lwm2m:"7"`
	HighThreshold           float64 `json:"highThreshold" lwm2m:"4"`
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
	return marshalJSON(fl)
}

func marshalJSON(lwm2mObject Lwm2mObject) ([]byte, error) {
	t := reflect.TypeOf(lwm2mObject)
	v := reflect.ValueOf(lwm2mObject)

	p := senml.Pack{senml.Record{
		BaseName:    fmt.Sprintf("%s/%s/", lwm2mObject.ID(), lwm2mObject.ObjectID()),
		BaseTime:    float64(lwm2mObject.Timestamp().Unix()),
		Name:        "0",
		StringValue: lwm2mObject.ObjectURN(),
	}}

	for i := 0; i < t.NumField(); i++ {
		var tagName, tagUnit = getTags(t.Field(i))

		if tagName == "" || tagName == "-" {
			continue
		}

		value := v.Field(i)

		r := senml.Record{
			Name: tagName,
		}

		if addValue(&r, value) {
			if tagUnit != "" {
				r.Unit = tagUnit
			}
			p = append(p, r)
		}
	}

	return json.Marshal(p)
}

func addValue(r *senml.Record, value reflect.Value) bool {
	kind := value.Kind()

	if kind == reflect.Float32 || kind == reflect.Float64 {
		v := value.Float()
		r.Value = &v
		return true
	}

	if kind == reflect.Int || kind == reflect.Int8 || kind == reflect.Int16 || kind == reflect.Int32 || kind == reflect.Int64 {
		v := float64(value.Int())
		r.Value = &v
		return true
	}

	if kind == reflect.Uint || kind == reflect.Uint8 || kind == reflect.Uint16 || kind == reflect.Uint32 || kind == reflect.Uint64 {
		v := float64(value.Uint())
		r.Value = &v
		return true
	}

	if kind == reflect.String {
		v := value.String()
		r.StringValue = v
		return true
	}

	if kind == reflect.Bool {
		v := value.Bool()
		r.BoolValue = &v
		return true
	}

	if kind == reflect.Ptr || kind == reflect.Pointer {
		if value.IsNil() {
			return false
		}
		return addValue(r, value.Elem())
	}

	return false
}

func getTags(f reflect.StructField) (string, string) {
	tag := f.Tag.Get("lwm2m")
	tags := strings.Split(tag, ",")

	if len(tags) > 1 {
		return tags[0], tags[1]
	}
	return tags[0], ""
}
