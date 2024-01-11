package domain

type Assets struct {
	LatestFullness    string  `json:"latestFullness"`
	FullnessThreshold string  `json:"fullnessThreshold"`
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

type FillingLevel struct {
	//	Mandatory, Filling level
	ActualFillingPercentage float64 `json:"actualFillingPercentage"`
	ContainerFull           bool    `json:"containerFull"`
	ContainerEmpty          bool    `json:"containerEmpty"`

	//	Optional, Filling level
	ActualFillingLevel int `json:"actualFillingLevel"`
}
