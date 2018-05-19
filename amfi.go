// Package amfi ..
//
// A small utility package to fetch latest NAV(Net Asset Value) of Indian mutual funds published by AMFI.
//
// DISCLAIMER: The package depends completely on the data published by AMFI.
package amfi

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Fund struct for mutual fund details from AMFI
type Fund struct {
	Code             string  `json:"code"`
	Name             string  `json:"name"`
	Isin             string  `json:"isin"`
	IsinReinvestment string  `json:"isinReinvestment"`
	Type             string  `json:"type"`
	Manager          string  `json:"manager"`
	NetAssetValue    float64 `json:"nav"`
	RepurchaseValue  float64 `json:"repurchaseValue"`
	SalePrice        float64 `json:"salePrice"`
	Date             string  `json:"date"`
}

// AMFI includes functions to load nav data from internet, get the list of funds and fund houses.
//
// Custom HTTPClient can also be used, based on the requirements
type AMFI struct {
	HTTPClient     *http.Client
	funds          map[string]Fund
	fundHouses     []string
	fundCategories []string
	lastUpdated    time.Time
}

// AMFI url to fetch latest nav data
const navURL = "https://www.amfiindia.com/spages/NAVAll.txt"

// Load the latest nav data from internet (amfi server: https://www.amfiindia.com/spages/NAVAll.txt)
func (amfi *AMFI) Load() error {
	if amfi.HTTPClient == nil {
		amfi.HTTPClient = http.DefaultClient
	}
	request, err := http.NewRequest(http.MethodGet, navURL, nil)
	if err != nil {
		return err
	}
	response, err := amfi.HTTPClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	amfi.processNavLines(string(data))
	amfi.lastUpdated = time.Now()
	return nil
}

// function to process the lines and categorize different types of lines
func (amfi *AMFI) processNavLines(data string) {
	var (
		tempFundHouses []string
		currentManager string
		currentType    string
		skipHeader     bool
	)
	amfi.funds = make(map[string]Fund)
	for _, line := range strings.Split(data, "\r\n") {
		line = strings.TrimSpace(line)
		// to skip the empty lines
		if len(line) == 0 {
			continue
		}
		if strings.Index(line, ";") > -1 {
			// to skip the header line
			if !skipHeader {
				skipHeader = true
				continue
			}
			// building Fund object from ; separated lines
			fund := amfi.buildFund(line, currentType, currentManager)
			amfi.funds[fund.Code] = fund
		} else {
			if strings.HasPrefix(line, "Open Ended") || strings.HasPrefix(line, "Close Ended") {
				currentType = line
				amfi.fundCategories = append(amfi.fundCategories, line)
			} else {
				currentManager = line
				tempFundHouses = append(tempFundHouses, line)
			}
		}
	}
	// remove duplicate items from fund houses list and assign
	amfi.fundHouses = append(amfi.fundHouses, removeDuplicates(tempFundHouses)...)
}

// function to parse the nav lines
func (amfi *AMFI) buildFund(line, currentType, currentManager string) Fund {
	values := strings.Split(line, ";")
	var fund Fund
	fund.Code = values[0]
	fund.Isin = values[1]
	fund.IsinReinvestment = values[2]
	fund.Name = values[3]
	if nav, err := strconv.ParseFloat(values[4], 64); err == nil {
		fund.NetAssetValue = nav
	}
	if repurchasePrice, err := strconv.ParseFloat(values[5], 64); err == nil {
		fund.RepurchaseValue = repurchasePrice
	}
	if salePrice, err := strconv.ParseFloat(values[6], 64); err == nil {
		fund.SalePrice = salePrice
	}
	fund.Date = values[7]
	fund.Manager = currentManager
	fund.Type = currentType
	return fund
}

// GetFundCategories returns the list of different categories of mutual funds
func (amfi *AMFI) GetFundCategories() []string {
	return amfi.fundCategories
}

// GetFundHouses returns the list of mutual fund houses
func (amfi *AMFI) GetFundHouses() []string {
	return amfi.fundHouses
}

// GetFunds returns a list of mutual funds with its latest nav data
func (amfi *AMFI) GetFunds() []Fund {
	var funds []Fund
	for _, fund := range amfi.funds {
		funds = append(funds, fund)
	}
	return funds
}

// GetFund returns fund details for the given SchemeCode.
//
// Returns an error ("Invalid Code"), if the input schemeCode is invalid
func (amfi *AMFI) GetFund(schemeCode string) (Fund, error) {
	fund, exists := amfi.funds[schemeCode]
	if exists {
		return fund, nil
	}
	return fund, errors.New("Invalid Code")
}

// GetLastUpdatedTime returns the timestamp of data fetched from the internet
func (amfi *AMFI) GetLastUpdatedTime() time.Time {
	return amfi.lastUpdated
}

// utility function to remove duplicate values in a string slice
func removeDuplicates(input []string) []string {
	keys := make(map[string]bool)
	output := []string{}
	for _, item := range input {
		if _, value := keys[item]; !value {
			keys[item] = true
			output = append(output, item)
		}
	}
	return output
}
