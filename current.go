// Copyright 2015 Brian J. Downs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openweathermap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// CurrentWeatherData struct contains an aggregate view of the structs
// defined above for JSON to be unmarshaled into.
type CurrentWeatherData struct {
	GeoPos  Coordinates        `json:"coord"`
	Sys     Sys                `json:"sys"`
	Base    string             `json:"base"`
	Weather []Weather          `json:"weather"`
	Main    Main               `json:"main"`
	Wind    Wind               `json:"wind"`
	Clouds  Clouds             `json:"clouds"`
	Rain    map[string]float64 `json:"rain"`
	Snow    map[string]float64 `json:"snow"`
	Dt      int                `json:"dt"`
	ID      int                `json:"id"`
	Name    string             `json:"name"`
	Cod     int                `json:"cod"`
	Unit    string
	Lang    string
	Key     string
}

// NewCurrent returns a new CurrentWeatherData pointer with the supplied parameters
func NewCurrent(unit, lang string) (*CurrentWeatherData, error) {
	unitChoice := strings.ToUpper(unit)
	langChoice := strings.ToUpper(lang)

	c := &CurrentWeatherData{}

	if ValidDataUnit(unitChoice) {
		c.Unit = DataUnits[unitChoice]
	} else {
		return nil, errUnitUnavailable
	}

	if ValidLangCode(langChoice) {
		c.Lang = langChoice
	} else {
		return nil, errLangUnavailable
	}

	c.Key = getKey()

	return c, nil
}

// SetLang allows you to set the language responses will be displayed as.  This isn't part of the
// NewCurrent call because it'd keep it easier to go with API defaults and
// adjust if explicitly called.
func (w *CurrentWeatherData) SetLang(lang string) error {
	if !ValidLangCode(lang) {
		return errLangUnavailable
	}
	w.Lang = lang

	return nil
}

// CurrentByName will provide the current weather with the provided
// location name.
func (w *CurrentWeatherData) CurrentByName(location string) error {
	var err error
	var response *http.Response

	response, err = http.Get(fmt.Sprintf(fmt.Sprintf(baseURL, "appid=%s&q=%s&units=%s&lang=%s"), w.Key, url.QueryEscape(location), w.Unit, w.Lang))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err := json.NewDecoder(response.Body).Decode(&w); err != nil {
		return err
	}

	return nil
}

// CurrentByCoordinates will provide the current weather with the
// provided location coordinates.
func (w *CurrentWeatherData) CurrentByCoordinates(location *Coordinates) error {
	response, err := http.Get(fmt.Sprintf(fmt.Sprintf(baseURL, "appid=%s&lat=%f&lon=%f&units=%s&lang=%s"), w.Key, location.Latitude, location.Longitude, w.Unit, w.Lang))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&w); err != nil {
		return err
	}

	return nil
}

// CurrentByID will provide the current weather with the
// provided location ID.
func (w *CurrentWeatherData) CurrentByID(id int) error {
	response, err := http.Get(fmt.Sprintf(fmt.Sprintf(baseURL, "appid=%s&id=%d&units=%s&lang=%s"), w.Key, id, w.Unit, w.Lang))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&w); err != nil {
		return err
	}

	return nil
}

// CurrentByZip will provide the current weather for the
// provided zip code.
func (w *CurrentWeatherData) CurrentByZip(zip int, countryCode string) error {
	response, err := http.Get(fmt.Sprintf(fmt.Sprintf(baseURL, "appid=%s&zip=%d,%s&units=%s&lang=%s"), w.Key, zip, countryCode, w.Unit, w.Lang))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&w); err != nil {
		return err
	}

	return nil
}

// CurrentByArea will provide the current weather for the
// provided area.
func (w *CurrentWeatherData) CurrentByArea() {}
