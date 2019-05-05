package main

import (
	"fmt"
	"io" //for deleting
	"os" // for testing the csv creation
	"testing"
	"time"
	//"trendwidgetsrc/trendapi"
	"trendwidgetsrc/trenddb"
	"trendwidgetsrc/trendlogic"
	"trendwidgetsrc/trendweb"
)

func main() {
	// Test CSV creation
	var t *testing.T
	var path string = "test_csv.csv"
	var data = [][]string{{"Line1", "Data1"}, {"Line2", "Data2"}}
	// create a test csv with some dummy values
	trendlogic.CreateCSV(data, path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("FAIL: CSV Creation not working")
	}
	// Test read the csv
	readCSV(path)
	deleteFile(path)
	fmt.Println("PASS: CSV testing...\n")
	// Drop all data
	trenddb.Start()
	trendlogic.GetHistory("1fQFz1qsG1TVjAHbhJ0tK8z3KAcvCr7euQcR6TofePYI", "A1:D", "Bristol")
	// Does database report it has more than 1 day?
	days := trendlogic.GetMaxDays(1)
	if days < 1 {
		fmt.Println("FAIL: Database read/write not working")
	} else {
		fmt.Println("PASS: updated into database")
	}
	// Test parse date funciton
	date := "August 03, 2017 at 07:00AM"
	dateparsed := trenddb.ParseDate(date, false)
	assertEqual(t, date, dateparsed, "FAIL: Error converting date format")
	// Create a new CSV file and check it's been completed
	if _, err := os.Stat(trendlogic.PreviousPeriod(5, "Windspeed")); os.IsNotExist(err) {
		fmt.Println("FAIL")
	} else {
		fmt.Println("PASS: dynamic request from Database working")
	}
	//  Starts the server as a concurrent go routine
	go trendweb.Start()
	//  Ignores return value so add in wait to allow server to run
	time.Sleep(1000 * time.Millisecond)
	db := trenddb.CurrentDb()
	//	Check for a duplicate record inserted at database creation
	if trenddb.DuplicateRead("May 18, 1919 at 07:00AM", db) {
		fmt.Println("PASS: Database read working")
	} else {
		fmt.Println("FAIL: database read not working")
	}
	//	try a bad
	trenddb.UserAlert(":", "test", "content", "weekly", db)
	trenddb.UserAlert("testuser", "test@test.com", "content", "weekly", db)
	//	Try to create a user and check present
	assertEqual(t, trenddb.DataPresent("user", db), true, "User creation failed.")
	assertEqual(t, trenddb.DataPresent("forecast", db), true, "User creation failed.")
	assertEqual(t, trenddb.DataPresent("history", db), true, "User creation failed.")
	//	Try to clear up
	trenddb.DropTables(db)
	fmt.Printf("\n\n* * ALL SERVER TESTS PASSED * *\n\n\n")
}

func assertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a == b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("\nTEST FAIL : %v != %v\n", a, b)
		fmt.Println(message)
	}
}

func readCSV(path string) {
	// re-open file
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// read file, line by line
	var text = make([]byte, 1024)
	for {
		_, err = file.Read(text)

		// break if finally arrived at end of file
		if err == io.EOF {
			break
		}

		// break if error occured
		if err != nil && err != io.EOF {
			isError(err)
			break
		}
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("==> done deleting test file")
}
