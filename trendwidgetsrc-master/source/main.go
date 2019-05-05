package main

import (
	"fmt"
	"trendwidgetsrc/trendlogic"
	"trendwidgetsrc/trendweb"
)

//	Tool set up for the following presently:
const weatherHistoryID string = "1fQFz1qsG1TVjAHbhJ0tK8z3KAcvCr7euQcR6TofePYI"
const weatherForecastID string = "1F9v9TE7cQ3GEKkyP28z6AG--4Iupn4uNt1zndmiQHhY"
const HistoryRange string = "A1:D"
const ForecastRange string = "A1:P"
const location string = "Bristol"

func main() {

	// Initialise the API and gather data -> put into database
	trendlogic.GetHistory(weatherHistoryID, HistoryRange, location)
	trendlogic.GetForecast(weatherForecastID, ForecastRange, location)
	//	go create the CSVs to power the JS connectors
	generateCSVs()
	//  Start serving
	trendweb.Start()

}

func generateCSVs() {
	// Go get the files to serve
	csv0 := trendlogic.PreviousPeriod(365, "CurrentTempCelcius")
	// Go get the files to serve
	csv1 := trendlogic.PreviousPeriod(30, "CurrentTempCelcius")
	// Go get the files to serve
	csv2 := trendlogic.PreviousPeriod(30, "Humidity")
	// Go get the files to serve
	csv3 := trendlogic.PreviousPeriod(30, "Windspeed")
	// Go get the files to serve
	csv4 := trendlogic.PreviousPeriod(5, "CurrentTempCelcius")
	// Go get the files to serve
	csv5 := trendlogic.PreviousPeriod(5, "Humidity")
	// Go get the files to serve
	csv6 := trendlogic.PreviousPeriod(5, "Windspeed")
	// Go get the files to serve
	csv7 := trendlogic.PreviousPeriod(1, "CurrentTempCelcius")
	// Go get the files to serve
	csv8 := trendlogic.PreviousPeriod(1, "Humidity")
	// Go get the files to serve
	csv9 := trendlogic.PreviousPeriod(1, "Windspeed")

	fmt.Println("\n\n**File created: " + "localhost:8080/" + csv0)
	fmt.Println("\n\n**File created: " + "localhost:8080/" + csv1)
	fmt.Println("**File created: " + "localhost:8080/" + csv2)
	fmt.Println("**File created: " + "localhost:8080/" + csv3)
	fmt.Println("\n\n**File created: " + "localhost:8080/" + csv4)
	fmt.Println("**File created: " + "localhost:8080/" + csv5)
	fmt.Println("**File created: " + "localhost:8080/" + csv6)
	fmt.Println("\n\n**File created: " + "localhost:8080/" + csv7)
	fmt.Println("**File created: " + "localhost:8080/" + csv8)
	fmt.Println("**File created: " + "localhost:8080/" + csv9 + "\n\n")

}
