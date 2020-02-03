package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

func findStationIDAndName(mobileNo string) (string, string) {
	resp, err := http.Get(urlBusStationService +
		fmt.Sprintf("?serviceKey=%s&keyword=%s", getServiceKey(), mobileNo))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var sr BusStationResponse
	xmlDec := xml.NewDecoder(resp.Body)
	xmlDec.Decode(&sr)
	if sr.MsgHeader.ResultCode != "0" {
		log.Println(sr)
		// log.Println(sr.ComMsgHeader.ErrMsg
		// log.Println(sr.MsgHeader.ResultMessage)
		panic("somthing wrong in query station")
	}

	return sr.BusStationList.StationID, sr.BusStationList.StationName
}

func findBusNo(routeID string) string {
	resp, err := http.Get(urlBusRouteInfoService +
		fmt.Sprintf("?serviceKey=%s&routeId=%s", getServiceKey(), routeID))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var sr BusRouteInfoResponse
	xmlDec := xml.NewDecoder(resp.Body)
	xmlDec.Decode(&sr)
	if sr.MsgHeader.ResultCode != "0" {
		log.Println(sr)
		// log.Println(sr.ComMsgHeader.ErrMsg
		// log.Println(sr.MsgHeader.ResultMessage)
		panic("somthing wrong in query bus routeID")
	}

	return sr.BusRouteInfoItem.RouteName
}
