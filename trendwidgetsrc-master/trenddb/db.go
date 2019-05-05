package trenddb

/*
Package built from scratch using standard Golang packages

Note : for this to work a local sql server must be running on 127.0.0.1:3306
A user must have the following permissions:

mysql -u root

CREATE USER 'trend'@'localhost';
CREATE DATABASE trend;
GRANT ALL ON trend.* TO 'trend'@'localhost';
FLUSH PRIVILEGES;
USE trend;

*/

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"trendwidgetsrc/trendapi"
)

const (
	DB_HOST = "tcp(127.0.0.1:3306)"
	DB_NAME = "trend"
	DB_USER = /*"root"*/ "trend"
	DB_PASS = /*""*/ ""
)

// For passing data through DB
type Measurements struct {
	CheckTime          string
	CurrentCondition   string
	CurrentTempCelcius int
	HighTempCelcius    int
	LowTempCelcius     int
	Humidity           int
	Windspeed          int
	Maxwindspeed       int
	WindDirection      string
	Index              int
}

// Holds Measurements
type Request struct {
	Sheet    string
	Token    string
	Srange   string
	Data     map[string]*Measurements
	Location string
}

//	Connect, drop, create and ping then return connection to pool
func Start() {
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	// Defer waits until function completes
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Println("FAIL: Database error!")
	} else {
		fmt.Println("PASS: Database Ping success")
	}
	// Automatically drops and creates structure
	DropTables(db)
	time.Sleep(100 * time.Millisecond)
	CreateTables(db)
}

//	Update history based on request from logic - make call to DB
func UpdateWeatherHistory(newcall Request) (data Request) {
	db := CurrentDb()
	//	If the database is out of date... call the API sub-routine
	time.Sleep(100 * time.Millisecond)
	if needAPIupdate(db) {
		//	Use a separate thread to ensure caller can return quickly if API hangs
		//	New connection inside function
		apiSubRoutine(newcall)
	}
	time.Sleep(100 * time.Millisecond)
	newcall = UpdateFromDB(newcall, db)
	// Release the connection
	time.Sleep(100 * time.Millisecond)
	db.Close()
	return newcall

}

//	Gets the latest forecast (future)
func UpdateWeatherForecast(newcall Request) (data Request) {
	db := CurrentDb()
	//	If the database is out of date... call the API sub-routine
	newcall = UpdateForecastFromDB(newcall, db)
	// Release the connection
	time.Sleep(100 * time.Millisecond)
	db.Close()
	return newcall

}

// is the DB out of date?
func needAPIupdate(db *sql.DB) (outfdate bool) {
	rows, err := db.Query("SELECT * FROM reading WHERE yyyymmdd BETWEEN NOW()- INTERVAL 1 DAY AND NOW()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	if rows.Next() {
		// There is a row from previous day so no need to update
		return false
	}
	return true
}

//	Get latest from database using request struct
func UpdateFromDB(newcall Request, db *sql.DB) (data Request) {
	rows, err := db.Query("SELECT id, stamp, ts, yyyymmdd,temperature,humidity,windspeed FROM reading WHERE yyyymmdd BETWEEN NOW()- INTERVAL 30 DAY AND NOW()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var stamp string
		var ts string
		var temperature int
		//	date is held in uint8 in Maria
		var yyyymmdd []uint8
		var humidity int
		var windspeed int
		var id int

		if err := rows.Scan(&id, &stamp, &ts, &yyyymmdd, &temperature, &humidity, &windspeed); err != nil {
			log.Fatal(err)
		}
		newcall.Data[stamp] = new(Measurements)
		newcall.Data[stamp].CurrentTempCelcius = temperature
		newcall.Data[stamp].Humidity = humidity
		newcall.Data[stamp].Windspeed = windspeed
		newcall.Data[stamp].CheckTime = strconv.Itoa(int(yyyymmdd[1]))
		newcall.Data[stamp].Index = id

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return newcall
}

//	Thread for updating the database from API
func apiSubRoutine(newcall Request) {
	// calls the api for the weather
	// need to get own connection as separate thread
	db := CurrentDb()
	fmt.Println("PASS: Updating database from API...\n")
	fmt.Println("PASS: Location set to " + newcall.Location)
	resp := trendapi.StartAPI(newcall.Sheet, newcall.Srange)
	//	Walks through the map values /responses from the request
	if len(resp.Values) == 0 {
		fmt.Println("FAIL: No data found.")
	} else {
		index := 0

		for _, row := range resp.Values {
			//You need type assertion to give us access to underlying type see: https://golang.org/doc/effective_go.html#interface_conversions
			var stamp string = row[0].(string)
			// note that maps are not safe for concurrent use!
			newcall.Data[stamp] = new(Measurements)
			newcall.Data[stamp].CurrentTempCelcius, _ = strconv.Atoi(row[1].(string))
			newcall.Data[stamp].Humidity, _ = strconv.Atoi(row[2].(string))
			newcall.Data[stamp].Windspeed, _ = strconv.Atoi(row[3].(string))
			newcall.Data[stamp].Index = index
			index++

			if stamp != "Date" {
				var temp int = newcall.Data[stamp].CurrentTempCelcius
				var hum int = newcall.Data[stamp].Humidity
				var wind int = newcall.Data[stamp].Windspeed
				var loc string = newcall.Location
				var date string = stamp

				newReading(stamp, temp, hum, wind, loc, date, db)
				time.Sleep(1 * time.Millisecond)
			}
		}
	}
	return
}

// Converts location string into ID and returns 0 if location not present
func LocationToid(name string, db *sql.DB) (id int) {

	rows, err := db.Query("SELECT name FROM location WHERE location.id = ?;", name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		return id
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return 0
}

//	Pass a reading into the database
func newReading(stamp string, temp int, humidity int, windspeed int, location string, yyyymmdd string, db *sql.DB) {
	// Check for duplicate based on key then add the values
	location_ := 1

	if len(yyyymmdd) > 5 {
		yyyymmdd = ParseDate(yyyymmdd, false)
	} else {
		yyyymmdd = "May 01, 1919 at 07:00AM"
	}

	if DuplicateRead(stamp, db) {
		return
	} else {
		stmt, es := db.Prepare("INSERT INTO Reading (stamp, windspeed, temperature, humidity, location, yyyymmdd) VALUES (?,?,?,?,?,?)")
		defer stmt.Close()
		if es != nil {
			panic(es.Error())
		}
		_, er := stmt.Exec(stamp, windspeed, temp, humidity, location_, yyyymmdd)
		if er != nil {
			panic(er.Error())
		}
		if stmt != nil {
			stmt.Close()
		}

	}
	return
}

//	Function to check whether record exists in DB by key (map)
func DuplicateRead(key string, db *sql.DB) (present bool) {

	rows, err := db.Query("SELECT id FROM Reading WHERE Reading.stamp = ?", key)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		return true
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return false
}

//	Return a pointer to current DB connection from pool beware repeat requests and please reuse
func CurrentDb() (db *sql.DB) {
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME + "?charset=utf8"
	db, err := sql.Open("mysql", dsn)
	//defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	return db

}

//	Create script reads from file
func CreateTables(db *sql.DB) {
	var filepath io.Reader

	filepath, err := os.Open("trenddb/create-drop.sql")

	if err != nil {
		fmt.Println("FAIL: Can't read file!")
		return
	}

	file, err := ioutil.ReadAll(filepath)

	if err != nil {
		fmt.Println("FAIL: Database error!")
		return
	}

	requests := strings.Split(string(file), ";")

	for _, request := range requests {
		if len(request) > 1 {
			_, err := db.Exec(request)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	fmt.Println("PASS: Created all database tables")
}

//	Drops to mirror create
func DropTables(db *sql.DB) {

	t := make([]string, 4)
	t[0] = "DROP TABLE IF EXISTS User"
	t[1] = "DROP TABLE IF EXISTS Forecast"
	t[2] = "DROP TABLE IF EXISTS Reading"
	t[3] = "DROP TABLE IF EXISTS Location"

	for td := range t {
		res, err := db.Exec(t[td])
		if err != nil {
			log.Fatal(err)
		}
		_, err = res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		_, err = res.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("PASS: Database tables dropped")
}

//	If intonly is set string is returned as only integers for dabase sequencing
func ParseDate(currentdate string, intonly bool) (returndate string) {
	newdate := strings.TrimSpace(currentdate)
	newdatesplit := strings.SplitAfter(newdate, ",")
	date, year := newdatesplit[0], newdatesplit[1]
	newyearsplit := strings.SplitAfter(year, "at")
	year = strings.TrimSuffix(newyearsplit[0], " at")
	year = strings.TrimPrefix(year, " ")

	datesplit2 := strings.SplitAfter(date, " ")
	day, month := datesplit2[1], datesplit2[0]
	day = strings.TrimSuffix(day, ",")

	switch month {
	case "January ":
		month = "1"
	case "February ":
		month = "2"
	case "March ":
		month = "3"
	case "April ":
		month = "4"
	case "May ":
		month = "5"
	case "June ":
		month = "6"
	case "July ":
		month = "7"
	case "August ":
		month = "8"
	case "September ":
		month = "9"
	case "October ":
		month = "10"
	case "November ":
		month = "11"
	case "December ":
		month = "12"
	}

	nk := year + "-" + month + "-" + day
	if intonly == true {
		nk = year + month + day
	}

	return nk
}

//	Executes a specific query on database to return count, max, min or av
func ExecuteValQuery(days int, request string, ranger string, db *sql.DB) (ret string, date string) {
	d := strconv.Itoa(days)

	query := "SELECT " + ranger + "(" + request + ") AS val, stamp FROM reading WHERE yyyymmdd BETWEEN NOW()- INTERVAL " + d + " DAY AND NOW()"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
		return "none", "none"
	}
	defer rows.Close()
	var val []uint8
	var stamp string
	for rows.Next() {
		if err := rows.Scan(&val, &stamp); err != nil {
			log.Fatal(err)
			return "none", "none"
		}
	}
	r := string(val)
	if len(d) >= 1 {
		return r, stamp
	}
	return "none", "none"
}

//	Forecast update from DB
func UpdateForecastFromDB(newcall Request, db *sql.DB) (data Request) {
	fmt.Println("PASS: Getting today's weather from API...")

	rows, err := db.Query("SELECT id, stamp, cond, current, hightemp, lowtemp, humidity, windspeed, maxWindspeed, winddirection, yyyymmdd from forecast WHERE ts BETWEEN NOW()- INTERVAL 1 DAY AND NOW()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var stamp string
		var cond string
		var current int
		var hightemp int
		var lowtemp int
		var humidity int
		var windspeed int
		var maxwindspeed int
		var winddirection string
		var yyyymmdd []uint8

		if err := rows.Scan(&id, &stamp, &cond, &current, &hightemp, &lowtemp, &humidity, &windspeed, &maxwindspeed, &winddirection, &yyyymmdd); err != nil {
			log.Fatal(err)
		}
		newcall.Data[stamp] = new(Measurements)
		newcall.Data[stamp].CurrentCondition = cond
		newcall.Data[stamp].CurrentTempCelcius = int(current)
		newcall.Data[stamp].Humidity = int(humidity)
		newcall.Data[stamp].Windspeed = int(windspeed)
		newcall.Data[stamp].CheckTime = strconv.Itoa(int(yyyymmdd[1]))
		newcall.Data[stamp].HighTempCelcius = int(hightemp)
		newcall.Data[stamp].LowTempCelcius = int(lowtemp)
		newcall.Data[stamp].WindDirection = string(lowtemp)
		newcall.Data[stamp].Maxwindspeed = int(lowtemp)
		newcall.Data[stamp].Index = int(id)

		fmt.Println("PASS history updated to: " + stamp)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return newcall
}

//	Add a user to be alerted based on cron job
func UserAlert(username string, email string, content string, frequency string, db *sql.DB) {

	if len(email) < 4 || len(username) < 5 {
		return
	}

	stmt, es := db.Prepare("INSERT INTO User (email, username, frequency, content, location) VALUES (?,?,?,?,?)")
	defer stmt.Close()
	if es != nil {
		panic(es.Error())
	}
	_, er := stmt.Exec(email, username, frequency, content, 1)
	if er != nil {
		panic(er.Error())
	}
	if stmt != nil {
		stmt.Close()
	}

}

// Testing function used in maintest for data validation
func DataPresent(tablename string, db *sql.DB) (present bool) {
	query := "SELECT * FROM " + tablename

	rows, err := db.Query(query)
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		return true
	}
	return true
}
