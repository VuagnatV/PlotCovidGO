package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	file, err := os.Open("covid-hospit-clage10-2023-03-31-18h01.csv")
	if err != nil {
		fmt.Println("Erreur en ouvrant le fichier CSV:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.LazyQuotes = true

	allLines, err := reader.ReadAll()

	var dates []string
	var totalHospitalizations []float64
	var totalDay float64 = 0
	var currDate string = ""

	dates = append(dates, allLines[1][2])

	lines := allLines[2:]

	for _, line := range lines {
		if currDate != line[2] {
			currDate = line[2]
			dates = append(dates, line[2])
			totalHospitalizations = append(totalHospitalizations, totalDay)
			totalDay = 0
		}

		num, err := strconv.ParseFloat(line[3], 64)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		totalDay = totalDay + num
	}

	totalHospitalizations = append(totalHospitalizations, totalDay)

	p := plot.New()

	pts := make(plotter.XYs, len(dates))
	for i, date := range dates {

		dateParsed, err := time.Parse("2006-01-02", date)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		pts[i].X = float64(dateParsed.Unix())
		pts[i].Y = totalHospitalizations[i]
	}
	line, err := plotter.NewLine(pts)
	if err != nil {
		fmt.Println("Error creating plot line:", err)
		return
	}
	p.Add(line)

	p.Title.Text = "COVID-related Hospitalizations per Day"
	p.X.Label.Text = "Date"
	p.Y.Label.Text = "Total Hospitalizations"

	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02"}

	if err := p.Save(8*vg.Inch, 4*vg.Inch, "hospitalizations_plot.png"); err != nil {
		fmt.Println("Error saving plot:", err)
		return
	}
	fmt.Println("Plot created successfully.")

}
