package main

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"os"
	"time"
)

// Config contains current settings of program
type Config struct {
	BaseInfoServiceKey   string `json:"baseinfoServicekey"`
	BusArrivalServiceKey string `json:"busarrivalServicekey"`
	BaseInfo             struct {
		UpdateDate time.Time `json:"updatedate"`
		Station    string    `json:"station"`
		Route      string    `json:"route"`
		// Area         string `json:"area"`
		// RouteLine    string `json:"routeline"`
		// RouteStation string `json:"routestation"`
	} `json:"baseinfo"`
}

const (
	configFileName = "config.json"
)

var (
	config Config
)

func loadConfig() error {
	if !isConfigValid() {
		resp, err := http.Get(urlBaseInfoService + "?serviceKey=" + getBaseInfoServiceKey())
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		var baseInfoResp BaseInfoResponse
		xmlDec := xml.NewDecoder(resp.Body)
		xmlDec.Decode(&baseInfoResp)

		resp.Body.Close()

		cleanupBaseInfoDir()

		if fPath, err := dlBaseInfo(baseInfoResp.BaseInfoItem.StationDownloadURL); err == nil {
			config.BaseInfo.Station = fPath
		} else {
			return err
		}

		if fPath, err := dlBaseInfo(baseInfoResp.BaseInfoItem.RouteDownloadURL); err == nil {
			config.BaseInfo.Route = fPath
		} else {
			return err
		}

		config.BaseInfo.UpdateDate = time.Now()

		w, err := os.Create(configFileName)
		if err != nil {
			return err
		}
		defer w.Close()

		prettyConfig, err := json.MarshalIndent(config, "", "    ")
		if err != nil {
			return err
		}
		w.Write(prettyConfig)
		return nil
	}

	confR, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	defer confR.Close()
	jDec := json.NewDecoder(confR)
	err = jDec.Decode(&config)
	if err != nil {
		return err
	}

	// TODO: compare config.BaseInfo.UpdateDate with time.Now() and
	// check update in base infos.

	return nil
}

func isConfigValid() bool {
	if !isExist(configFileName) {
		return false
	}

	confR, err := os.Open(configFileName)
	if err != nil {
		panic(err)
	}
	defer confR.Close()
	jDec := json.NewDecoder(confR)
	err = jDec.Decode(&config)
	if err != nil {
		panic(err)
	}

	// if !isExist(config.BaseInfo.Area) {
	// 	return false
	// }
	if !isExist(config.BaseInfo.Station) {
		return false
	}
	if !isExist(config.BaseInfo.Route) {
		return false
	}
	// if !isExist(config.BaseInfo.RouteLine) {
	// 	return false
	// }
	// if !isExist(config.BaseInfo.RouteStation) {
	// 	return false
	// }
	return true
}

func getBaseInfoServiceKey() string {
	serviceKey := os.Getenv("BASEINFOSERVICEKEY")
	if serviceKey != "" {
		return serviceKey
	}

	if config.BaseInfoServiceKey != "" {
		return config.BaseInfoServiceKey
	}

	panic("no servicekey")
}

func getBusArrivalServiceKey() string {
	serviceKey := os.Getenv("BUSARRIVALSERVICEKEY")
	if serviceKey != "" {
		return serviceKey
	}

	if config.BusArrivalServiceKey != "" {
		return config.BusArrivalServiceKey
	}

	panic("no servicekey")
}
