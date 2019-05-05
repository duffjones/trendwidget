package trendweb

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"net/url"
	"github.com/gorilla/mux"
	"time"
	"trendwidgetsrc/trenddb"
	"trendwidgetsrc/trendlogic"
)

func Start() {

	fmt.Println("PASS: Web server working")

	log.Printf("main: starting HTTP server")

	//  Adding go insures this runs in a separate routine
	srv := startHttpServer()
	//  Run server for how long?
	time.Sleep(100000 * time.Millisecond)
	// now close the server gracefully ("shutdown")
	defer srv.Shutdown(nil)
	// timeout could be given instead of nil as a https://golang.org/pkg/context/
	if err := srv.Shutdown(nil); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

	log.Printf("main: done. exiting")
}

func startHttpServer() *http.Server {

	router := mux.NewRouter()
	srv := &http.Server{Addr: ":8080", Handler: router}
	//Makes the CSS and related assets visible
	http.FileServer(http.Dir("filestore"))
	//Make the CSV files accessible
	router.HandleFunc("/weather", index)
	//	end point for POST of user data
	router.HandleFunc("/emailuser", emailuser)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("filestore")))
	fmt.Println("PASS: FILESERVER RUNNING>>")
	fmt.Println("PASS: SERVER RUNNING>>")

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
	// returning reference so caller can call Shutdown()
	return srv
}

//Function returns the index template (secure)
func index(w http.ResponseWriter, req *http.Request) {
	// Update the weather...
	trendlogic.UpdateWeather()
	// Get the template
	tpl, err := template.ParseFiles("filestore/templates/dynamicindex_xhtml.xhtml")
	if err != nil {
		log.Fatalln("error parsing template", err)
	}
	// Note must start with caps to be visible outside package
	payload := trendlogic.GetRanges()
	fmt.Println("PASS: Pulling correct averages, Top temp: " + payload.Month.Toptemp + " on " + payload.Month.Toptempdate)
	//Execution request returns the template from the appropraite struct
	err = tpl.ExecuteTemplate(w, "dynamicindex_xhtml.xhtml", payload)
	if err != nil {
		log.Fatalln("error executing template", err)
	}
}

//Function returns the index template (secure)
func emailuser(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Fatalln(err)
		return
	}
	name := req.PostFormValue("username")
	useremail := req.PostFormValue("useremail")
	frequency := req.PostFormValue("frequency")
	//selection := req.PostFormValue("interestingcharts")

	db := trenddb.CurrentDb()
	trenddb.UserAlert(name, useremail, "charts", frequency, db)
	// reload template
	index(w, req)

}
