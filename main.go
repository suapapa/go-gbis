package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"sort"
)

var (
	flagImageOut bool
	flagDebugGG  bool
)

func init() {
	flag.BoolVar(&flagImageOut, "i", false, "set if u want image output")
	flag.BoolVar(&flagDebugGG, "d", false, "draw guide line for gg elements")
	flag.Parse()
}

func main() {
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	mobileNo := flag.Args()[0] // 07-479 (H스퀘어)

	stationID, stationName := findStationIDAndName(mobileNo)
	resp, err := http.Get(urlBusArrivalServiceStation +
		fmt.Sprintf("?serviceKey=%s&stationId=%s", getBusArrivalServiceKey(), stationID))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	var sr BusArrivalStationResponse
	xmlDec := xml.NewDecoder(resp.Body)
	xmlDec.Decode(&sr)

	sort.Sort(sr.BusArrivalList)
	// print result in txt
	if !flagImageOut {
		printBusArrivalInfo(stationName, sr.BusArrivalList)
	} else {
		drawBusArrivalInfo(stationName, sr.BusArrivalList) // TODO : make out.png
	}
}

func printBusArrivalInfo(stationName string, buses []busArrival) {
	fmt.Printf("# %s #\n", stationName)
	for _, b := range buses {
		fmt.Printf("## 버스번호: %s ##\n", findBusNoFrom(b.RouteID))
		if b.PredictTime1 != "" && b.LocationNo1 != "" {
			fmt.Printf("* 다음버스: %s분 후 (%s 정류장 전)\n", b.PredictTime1, b.LocationNo1)
		}
		if b.PredictTime2 != "" && b.LocationNo2 != "" {
			fmt.Printf("* 다다음버스: %s분 후 (%s 정류장 전)\n", b.PredictTime2, b.LocationNo2)
		}
	}
}
