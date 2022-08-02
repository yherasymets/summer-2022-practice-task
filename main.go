package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

// змінна виконує роль пакету "encoding/json"
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	jsoniter.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		t, err := time.ParseInLocation("15:04:05", iter.ReadString(), time.UTC)
		if err != nil {
			iter.Error = err
			return
		}
		*((*time.Time)(ptr)) = t
	})

	var departureStation, arrivalStation, criteria string
	fmt.Print("Departure Station:")
	fmt.Scan(&departureStation)
	fmt.Print("Arrival Station:")
	fmt.Scan(&arrivalStation)
	fmt.Print("Criteria:")
	fmt.Scan(&criteria)

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		log.Fatal(err)
	}

	for _, train := range result {
		fmt.Printf("%#v \n", train)
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	if departureStation == "" {
		return nil, fmt.Errorf("empty departure station")
	}
	if arrivalStation == "" {
		return nil, fmt.Errorf("empty arrival station")
	}
	departureStationId, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, fmt.Errorf("bad departure station input")
	}
	arrivalStationId, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, fmt.Errorf("bad arrival station input")
	}

	body, err := os.ReadFile("./data.json")
	if err != nil {
		return nil, fmt.Errorf("read file error: %v", err)
	}
	var allTrains Trains
	err = json.Unmarshal(body, &allTrains)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling file error: %v", err)
	}

	var trains Trains
	for i := range allTrains {
		if arrivalStationId == allTrains[i].ArrivalStationID &&
			departureStationId == allTrains[i].DepartureStationID {
			trains = append(trains, allTrains[i])
		}
	}
	switch strings.ToLower(criteria) {
	case "price":
		sort.Slice(trains, func(i, j int) bool {
			if trains[i].Price == trains[j].Price {
				return trains[i].TrainID < trains[j].TrainID
			}
			return trains[i].Price < trains[j].Price
		})
	case "arrival-time":
		sort.Slice(trains, func(i, j int) bool {
			if trains[i].ArrivalTime.Equal(trains[j].ArrivalTime) {
				return trains[i].TrainID < trains[j].TrainID
			}
			return trains[i].ArrivalTime.Before(trains[j].ArrivalTime)
		})
	case "departure-time":
		sort.Slice(trains, func(i, j int) bool {
			if trains[i].DepartureTime.Equal(trains[j].DepartureTime) {
				return trains[i].TrainID < trains[j].TrainID
			}
			return trains[i].DepartureTime.Before(trains[j].DepartureTime)
		})
	default:
		return nil, fmt.Errorf("unsupported criteria")
	}

	if len(trains) < 3 {
		return trains, nil
	}
	return trains[:3], nil
}
