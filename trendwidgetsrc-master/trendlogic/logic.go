package trendlogic

//os and encoding/csv - one which handles the file interactions the other which converts the data sturcture into csv format.

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	//"strings"
	"trendwidgetsrc/trenddb"
)

var history trenddb.Request
var h_data = map[string]*trenddb.Measurements{}
var forecast trenddb.Request
var f_data = map[string]*trenddb.Measurements{}

type WeatherRange struct {
	Toptemp      string
	Toptempdate  string
	Topwind      string
	Topwinddate  string
	Tophumid     string
	Tophumiddate string
	Lowtemp      string
	Lowtempdate  string
	Lowwind      string
	Lowwinddate  string
	Lowhumid     string
	Lowhumiddate string
	Avtemp       string
	Avwind       string
	Avhumid      string
}

type Forecast struct {
	CheckTime        string
	CurrentTemp      string
	CurrentCondition string
	HighTemp         string
	LowTemp          string
	Humidity         string
	Windspeed        string
	MaxWindspeed     string
	WindDirection    string
}

type Data struct {
	Lastupdate string
	Month      *WeatherRange
	Week       *WeatherRange
	Year       *WeatherRange
	Historic   string
	Forecast   *Forecast
	Future     string
}

//	Initialise the logic
func GetHistory(_sheet string, _srange string, _location string) {
	// Creates tables in db and initialise
	//history.Rtype = "history"
	trenddb.Start()
	history.Sheet = _sheet
	history.Srange = _srange
	history.Data = h_data
	forecast.Data = f_data
	history.Location = _location
	// Updates the database
	UpdateWeather()
}

//	Initialise the logic
func GetForecast(_sheet string, _srange string, _location string) {
	// Creates tables in db and initialise
	//forecast.Rtype = "forecast"
	trenddb.Start()
	forecast.Sheet = _sheet
	forecast.Srange = _srange
	forecast.Data = f_data
	forecast.Location = _location
	// Updates the database
	UpdateWeather()
}

//	helper function to ensure all DB requests are in range
func GetMaxDays(days int) (max int) {
	db := trenddb.CurrentDb()
	id, _ := trenddb.ExecuteValQuery(365, "id", "MAX", db)

	val, _ := strconv.Atoi(id)

	if days > val {

		days = val
	}
	return days
}

func GetRanges() (returned Data) {
	returned.Month = new(WeatherRange)
	returned.Week = new(WeatherRange)
	returned.Year = new(WeatherRange)
	returned.Month = updateRangeData(30, returned.Month)
	returned.Week = updateRangeData(7, returned.Week)
	returned.Year = updateRangeData(365, returned.Year)
	returned.Historic = historicCopy(returned)
	returned.Future = forecastCopy(returned)
	return returned
}

func historicCopy(returned Data) (copy string) {
	text := `Over the past month  the average morning temperature has been ` + returned.Month.Avtemp + `. The coldest morning was on ` + returned.Month.Lowtempdate + ` and warmest on ` + returned.Month.Toptempdate + `. The fastest windspeed was on ` + returned.Month.Topwinddate + ` at ` + returned.Month.Topwind + ` MPH. The most humid day was ` + returned.Month.Tophumiddate + ` hitting ` + returned.Month.Tophumid + `% relative humidity. Looking at the past week, the highest temperate has been ` + returned.Week.Toptemp + `, with lows of ` + returned.Week.Lowtemp + ` on ` + returned.Week.Lowtempdate +
		`The hottest morning of the past year was ` + returned.Year.Toptempdate + ` with a reading of ` + returned.Year.Toptemp + `. The coldest was ` + returned.Year.Lowtemp + ` on ` + returned.Year.Lowtempdate + `. It was very windy on ` + returned.Year.Topwinddate + ` reaching a speed of ` + returned.Year.Topwind + `.`

	return text
}

//	This function is still calling historic data
func forecastCopy(returned Data) (copy string) {
	text := `Today's forecast is ` + returned.Month.Avtemp + `degrees. We expect lows of ` + returned.Month.Lowtemp + ` with fastest wind speed of ` + returned.Month.Topwind + ` MPH. `
	return text
}

func updateRangeData(days int, request *WeatherRange) (returned *WeatherRange) {
	//	make sure we have the number of days for request
	days = GetMaxDays(days)
	db := trenddb.CurrentDb()

	request.Toptemp, request.Toptempdate = trenddb.ExecuteValQuery(days, "temperature", "MAX", db)
	request.Topwind, request.Topwinddate = trenddb.ExecuteValQuery(days, "windspeed", "MAX", db)
	request.Tophumid, request.Tophumiddate = trenddb.ExecuteValQuery(days, "humidity", "MAX", db)
	request.Avtemp, _ = trenddb.ExecuteValQuery(days, "temperature", "AVG", db)
	request.Avwind, _ = trenddb.ExecuteValQuery(days, "windspeed", "AVG", db)
	request.Avhumid, _ = trenddb.ExecuteValQuery(days, "humidity", "AVG", db)
	request.Lowtemp, request.Lowtempdate = trenddb.ExecuteValQuery(days, "temperature", "MIN", db)
	request.Lowwind, request.Lowwinddate = trenddb.ExecuteValQuery(days, "windspeed", "MIN", db)
	request.Lowhumid, request.Lowhumiddate = trenddb.ExecuteValQuery(days, "humidity", "MIN", db)

	return request

}

//	Makes a call to the database which decides whether to call the API
func UpdateWeather() {
	history = trenddb.UpdateWeatherHistory(history)
	forecast = trenddb.UpdateWeatherForecast(forecast)
}

// Turns a request into a string for parsing into CSV
func RequestToString(days int, request string) (retstring string) {
	//	make sure we have the number of days for request
	days = GetMaxDays(days)
	returnstring := CSVtoString(PreviousPeriod(days, request))
	return returnstring

}

func CSVtoString(path string) (retstring string) {

	var filepath io.Reader

	filepath, err := os.Open(path)

	if err != nil {
		fmt.Println("FAIL: Can't read file!")
	}

	file, err := ioutil.ReadAll(filepath)

	if err != nil {
		fmt.Println("FAIL: reading error!")
	}

	str := fmt.Sprintf("PASS: %s created", file)

	return str
}

func PreviousPeriod(days int, request string) (path string) {
	//	make sure we have the number of days for request
	days = GetMaxDays(days)

	location := "filestore/csv/"
	filepath := location + request + "_" + strconv.Itoa(days) + "days_" + "trend.csv"

	CreateCSV(createDayArray(days, request), filepath)

	return filepath
}

func ForecastPath(days int, request string) (path string) {
	//	make sure we have the number of days for request
	days = GetMaxDays(days)

	location := "filestore/csv/"
	filepath := location + request + "_" + strconv.Itoa(days) + "forecast_" + "trend.csv"

	CreateCSV(createDayArray(days, request), filepath)

	return filepath
}

func createDayArray(days int, request string) (values [][]string) {
	//	make sure we have the number of days for request
	days = GetMaxDays(days)

	retvalues := [][]string{}

	//adds Date header as first value of csv
	//request is the title of the content wanted ie "humidity"
	header := []string{"Date", request}
	retvalues = append(retvalues, header)

	//UpdateWeather()

	most_recent := 0
	most_recent_key := ""

	//	Find most recent
	for k, v := range history.Data {
		if v.Index > most_recent {
			most_recent = v.Index
			most_recent_key = k
		}
	}

	dates := []string{}

	for k, v := range history.Data {
		if v.Index > (most_recent - days) {
			//Standardise the date formatting
			k = trenddb.ParseDate(k, false)
			//k = parseDate(k, false)
			dates = append(dates, k)
			v.CheckTime = k
			sort.Strings(dates)
		}
	}

	for i := 0; i < len(dates); i++ {
		for _, v := range history.Data {
			if v.CheckTime == dates[i] {
				if v.Index > (most_recent - days) {
					row := []string{v.CheckTime, strconv.Itoa(v.CurrentTempCelcius)}
					switch request {
					case "CurrentTempCelcius":
						row = []string{v.CheckTime, strconv.Itoa(v.CurrentTempCelcius)}
					case "Humidity":
						row = []string{v.CheckTime, strconv.Itoa(v.Humidity)}
					case "Windspeed":
						row = []string{v.CheckTime, strconv.Itoa(v.Windspeed)}
					}
					retvalues = append(retvalues, row)
				}
			}

		}

	}

	fmt.Println("\nPASS:"+request+" data current to = ", most_recent_key+"\n")
	return retvalues
}

// adapted from https://golangcode.com/write-data-to-a-csv-file/
func CreateCSV(data [][]string, path string) {
	file, err := os.Create(path)
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
