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
	SchemeCode       string  `json:"schemeCode"`
	Isin             string  `json:"isin"`
	IsinReinvestment string  `json:"isinReinvestment"`
	SchemeName       string  `json:"schemeName"`
	NetAssetValue    float64 `json:"nav"`
	RepurchaseValue  float64 `json:"repurchaseValue"`
	SalePrice        float64 `json:"salePrice"`
	Date             string  `json:"date"`
}

// AMFI holds the list of funds and fund houses
// includes functions to load nav data from internet, get the list of funds and fund houses
// the network timeout is set to 2 seconds, default
type AMFI struct {
	Timeout        time.Duration
	funds          map[string]Fund
	fundHouses     []string
	fundCategories []string
}

const navURL = "https://www.amfiindia.com/spages/NAVAll.txt"

// Load the latest nav data from internet (amfi india server)
func (amfi *AMFI) Load() error {
	var timeout = amfi.Timeout
	if timeout == 0 {
		timeout = 2
	}
	httpClient := &http.Client{
		Timeout: time.Second * timeout,
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
	var navLines []string
	var tempFundHouses []string
	for _, line := range strings.Split(data, "\r\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.Index(line, ";") > -1 {
			navLines = append(navLines, line)
		} else {
			if strings.HasPrefix(line, "Open Ended") || strings.HasPrefix(line, "Close Ended") {
				amfi.fundCategories = append(amfi.fundCategories, line)
			} else {
				tempFundHouses = append(tempFundHouses, line)
			}
		}
	}
	// removing duplicate items from fund houses list
	amfi.fundHouses = append(amfi.fundHouses, removeDuplicates(tempFundHouses)...)
	// buildind slice of Fund from ; separated lines
	// ignoring the header line
	amfi.buildFundList(navLines[1:])
}

// function to parse the nav lines
func (amfi *AMFI) buildFundList(lines []string) {
	amfi.funds = make(map[string]Fund)
	for _, line := range lines {
		var fund Fund
		values := strings.Split(line, ";")
		fund.SchemeCode = values[0]
		fund.Isin = values[1]
		fund.IsinReinvestment = values[2]
		fund.SchemeName = values[3]
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
		amfi.funds[fund.SchemeCode] = fund
	}
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
