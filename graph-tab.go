package main

import (
	"bytes"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"github.com/wcharczuk/go-chart/v2"
)

func (app *Config) graphTab() *fyne.Container {
	chart := app.getChart()
	chartContainer := container.NewVBox(chart)
	app.GraphContainer = chartContainer
	return chartContainer
}

func (app *Config) getChart() *canvas.Image {

	var img *canvas.Image
	buf, err := app.makeChart()
	if err != nil {
		img = canvas.NewImageFromResource(resourceUnreachablePng)
	} else {
		img = canvas.NewImageFromReader(buf, "RateGraph")
	}

	img.SetMinSize(fyne.Size{
		Width:  770,
		Height: 410,
	})

	// img.FillMode = canvas.ImageFillOriginal

	return img

}

func (app *Config) makeChart() (*bytes.Buffer, error) {
	values, err := app.getData()
	if err != nil {
		app.ErrorLog.Print("failed to get latest assets by accounts")
		return nil, err
	}
	pie := chart.PieChart{
		Width:  770,
		Height: 410,
		Values: values,
	}
	// Render the bar chart to a byte array.
	buf := new(bytes.Buffer)
	err = pie.Render(chart.PNG, buf)
	if err != nil {
		app.ErrorLog.Print("failed to render bar chart: %v", err)
		return nil, err
	}
	return buf, nil

}

func (app *Config) getData() ([]chart.Value, error) {
	assets, err := app.DB.SelectLatestAssetGroupbyAccount()
	if err != nil {
		app.ErrorLog.Print("failed to get latest assets by accounts")
		return nil, err
	}

	var sum float64
	for _, asset := range assets {
		sum += asset.Asset
	}
	var values []chart.Value

	for _, asset := range assets {
		value := chart.Value{Value: asset.Asset, Label: fmt.Sprintf("%s (%.2f%%)", asset.Name, asset.Asset/sum*100)}
		values = append(values, value)
	}
	return values, nil
}
