package main

import "investmentwatcher/repository"

type Overview struct {
	OriginalBalance float64
	CurrentAsset    float64
	Rate            float64
	Change          float64
}

func (app *Config) GetOverview() (*Overview, error) {
	var total []repository.TotalDetail
	total, err := app.DB.SelectTotalDetailGropByRegDt()
	if err != nil {
		app.ErrorLog.Println("failed to get recent data")
		return nil, err
	}

	var original_balance, current, rate, change float64

	if len(total) == 0 {
		original_balance, current, rate, change = 0, 0, 0, 0
	} else if len(total) == 1 {
		original_balance, current, change = total[0].Balance, total[0].Asset, 0
		rate = (current - original_balance) / original_balance * 100
	} else {
		original_balance, current = total[0].Balance, total[0].Asset
		change = (current - total[1].Asset) / total[1].Asset * 100
		rate = (current - original_balance) / original_balance * 100
	}
	var assetInfo = Overview{
		OriginalBalance: original_balance,
		CurrentAsset:    current,
		Rate:            rate,
		Change:          change,
	}

	return &assetInfo, nil

}
