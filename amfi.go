package amfi

import (
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

// AMFI includes functions to load nav data from internet, get the list of funds and fund houses
// custom HTTPClient can be used,based on the requirements
type AMFI struct {
	HTTPClient     *http.Client
	funds          map[string]Fund
	fundHouses     []string
	fundCategories []string
	lastUpdated    time.Time
}

const navURL = "https://www.amfiindia.com/spages/NAVAll.txt"

// Load the latest nav data from internet (amfi india server)
func (amfi *AMFI) Load() error {
	var httpClient = amfi.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	request, err := http.NewRequest(http.MethodGet, navURL, nil)
	if err != nil {
		return err
	}
	request.Header.Set("User-Agent", "go-amfi-client")
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	amfi.processNavLines(string(data))
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
		if len(line) == 0 {
			continue
		}
		if strings.Index(line, ";") > -1 {
			// to skip the header line
			if !skipHeader {
				skipHeader = true
				continue
			}
			// building slice of Fund from ; separated lines
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
	// removing duplicate items from fund houses list
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

// GetFunds returns the map of mutual funds with its latest nav data (SchemeCode as key, and Fund as value)
func (amfi *AMFI) GetFunds() map[string]Fund {
	return amfi.funds
}

// GetFund returns the fund details for the give SchemeCode
func (amfi *AMFI) GetFund(schemeCode string) Fund {
	return amfi.funds[schemeCode]
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
