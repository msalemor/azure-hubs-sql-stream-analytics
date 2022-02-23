package common

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	RaiseAnomaly    = 3
	AC_VENT         = 1
	GENERATOR_EVENT = 2
	MOTOR_EVENT     = 3
)

func getRandom(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func GetRandomEvent() string {
	anomaly := getRandom(1, 6)
	eventType := getRandom(1, 4)

	airTemperature := float64(getRandom(600, 800)) / 10.0
	airFlow := float64(getRandom(30, 40)) / 10.0
	coolantTemperature := float64(getRandom(200, 400)) / 10.0
	gasPercentage := float64(getRandom(1, 1000)) / 10.0
	voltage := float64(getRandom(2300, 2450)) / 10.0
	motorTemp := float64(getRandom(1800, 2000)) / 10.0
	motorRevolutions := float64(getRandom(2000, 5000)) / 10.0
	hertz := float64(getRandom(580, 650)) / 10.0
	amps := float64(getRandom(150, 250)) / 10.0

	if anomaly == RaiseAnomaly {
		voltage = 0
		motorTemp = 0
		motorRevolutions = 0
		gasPercentage = 10
		airFlow = 0
		airTemperature = 90
	}

	var event Event
	if eventType == AC_VENT {
		event = NewACEvent(airFlow, airTemperature, coolantTemperature)
		jsonBytes, _ := json.Marshal(event)
		return string(jsonBytes)

	} else if eventType == GENERATOR_EVENT {
		event = NewGeneratorEvent(hertz, amps, voltage, gasPercentage)
		jsonBytes, _ := json.Marshal(event)
		return string(jsonBytes)

	} else {
		event = NewMotorEvent(motorTemp, motorRevolutions)
		jsonBytes, _ := json.Marshal(event)
		return string(jsonBytes)
	}
}
