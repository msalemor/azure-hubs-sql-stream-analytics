package common

import "time"

type Event interface{}

type acEvent struct {
	Ts                 time.Time `json:"ts"`
	Type               string    `json:"type"`
	DeviceID           string    `json:"deviceId"`
	AirFlow            float64   `json:"airflow"`
	AirTemperature     float64   `json:"airTemperature"`
	CoolantTemperature float64   `json:"coolantTemperature"`
}

type generatorEvent struct {
	Ts            time.Time `json:"ts"`
	Type          string    `json:"type"`
	DeviceID      string    `json:"deviceId"`
	Hertz         float64   `json:"hertz"`
	Amps          float64   `json:"amps"`
	Voltage       float64   `json:"voltage"`
	GasPercentage float64   `json:"gasPercentage"`
}

type motorEvent struct {
	Ts          time.Time `json:"ts"`
	Type        string    `json:"type"`
	DeviceID    string    `json:"deviceId"`
	Temperature float64   `json:"temperature"`
	Revolutions float64   `json:"revolutions"`
}

type EventsRequest struct {
	Count int  `json:"count"`
	Delay int  `json:"delay"`
	Batch bool `json:"batch"`
}

type EventsResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
	Delay   int    `json:"delay"`
	Batch   bool   `json:"batch"`
}

//{"ts":"2022-02-23T18:52:40.8653687Z","deviceId":"100","eventType":"MotorAnomality","property":"revolutions","value":0}
type AnomalyEvent struct {
	Ts        time.Time `json:"ts"`
	DeviceID  string    `json:"deviceId"`
	EventType string    `json:"eventType"`
	Property  string    `json:"Property"`
	Value     float64   `json:"Value"`
}

func NewACEvent(airFlow, airTemperature, coolantTemperature float64) Event {
	return acEvent{
		Ts:                 time.Now().UTC(),
		Type:               "ACEvent",
		DeviceID:           "100",
		AirFlow:            airFlow,
		AirTemperature:     airTemperature,
		CoolantTemperature: coolantTemperature,
	}
}

func NewGeneratorEvent(hertz, amps, voltage, gasPercentage float64) Event {
	return generatorEvent{
		Ts:            time.Now().UTC(),
		Type:          "GeneratorEvent",
		DeviceID:      "100",
		Hertz:         hertz,
		Amps:          amps,
		Voltage:       voltage,
		GasPercentage: gasPercentage,
	}
}

func NewMotorEvent(revolutions, temperature float64) Event {
	return motorEvent{
		Ts:          time.Now().UTC(),
		Type:        "MotorEvent",
		DeviceID:    "100",
		Temperature: temperature,
		Revolutions: revolutions,
	}
}
