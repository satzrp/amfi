# AMFI - Utility package to get latest NAV data from AMFI

A small utility package to fetch latest NAV(Net Asset Value) of Indian mutual funds published by AMFI.

DISCLAIMER: The package depends completely on the data published by AMFI.

## Adding to your project
```sh
go get github.com/LordOfSati/amfi
```

## Usage
```go
import "github.com/LordOfSati/amfi"

func main() {
  amfi := &amfi.AMFI{}
  err := amfi.Load()
  if err != nil {
    funds := amfi.GetFunds()
    fundHouses := amfi.GetFundHouses()
    fundCategories := amfi.GetFundCategories()
    // to get details of single fund
    fund, err := amfi.GetFund("103155")
  }
}
```
## Sample Fund Details in JSON format
```json
{
  "code": "120518",
  "name": "Aditya Birla Sun Life Balanced '95 Fund - Direct Plan-Dividend",
  "isin": "INF209KA1LH3",
  "isinReinvestment": "INF209K01ZB2",
  "type": "Open Ended Schemes(Balanced)",
  "house": "Aditya Birla Sun Life Mutual Fund",
  "nav": 207.77,
  "repurchaseValue": 205.69,
  "salePrice": 207.77,
  "date": "18-May-2018"
}
```