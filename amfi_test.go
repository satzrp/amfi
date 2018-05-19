package amfi

import (
	"testing"
)

func TestAMFI(t *testing.T) {
	amfi := &AMFI{}
	t.Run("Test Load function", func(t *testing.T) {
		if err := amfi.Load(); err != nil {
			t.Errorf("Test Load function: FAILED, error = %v", err)
		}
	})
	t.Run("Test GetFundCategories function", func(t *testing.T) {
		if fundCategories := amfi.GetFundCategories(); fundCategories == nil {
			t.Error("Test GetFundCategories function: FAILED")
		}
	})
	t.Run("Test GetFundHouses function", func(t *testing.T) {
		if fundHouses := amfi.GetFundHouses(); fundHouses == nil {
			t.Error("Test GetFundHouses function: FAILED")
		}
	})
	t.Run("Test GetFunds function", func(t *testing.T) {
		if funds := amfi.GetFunds(); funds == nil {
			t.Error("Test GetFunds function: FAILED")
		}
	})
	t.Run("Test GetFund function", func(t *testing.T) {
		funds := amfi.GetFunds()
		if _, err := amfi.GetFund(funds[0].Code); err != nil {
			t.Error("Test GetFund function: FAILED")
		}
	})
	t.Run("Test GetFund function - Error", func(t *testing.T) {
		if _, err := amfi.GetFund("ABCD"); err == nil {
			t.Error("Test GetFund function - Error: FAILED")
		}
	})
	t.Run("Test GetLastUpdatedTime function", func(t *testing.T) {
		if lastUpdate := amfi.GetLastUpdatedTime(); lastUpdate.IsZero() {
			t.Error("Test GetLastUpdatedTime function: FAILED")
		}
	})
}
