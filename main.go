package main

import (
	"log"
	"net/http"

	"github.com/alexcooperse/ipv6world/csv"
	"github.com/alexcooperse/ipv6world/geo"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

func main() {
	// get the most up to date data
	locations := csv.ReadFile("./csv/data/GeoLite2-City-Blocks-IPv6.csv")

	// connect to mongodb for persistent storage
	mongoSession, dbErr := mgo.Dial("localhost")
	if dbErr != nil {
		log.Panic(dbErr)
	}
	defer mongoSession.Close()

	// remove existing data from the collection to ensure a fresh start
	collection := mongoSession.DB("ipv6").C("addresses")
	_, err := collection.RemoveAll(nil)
	if err != nil {
		log.Panic(err)
	}

	// refresh the collection with the most recent data
	for point, count := range locations {
		coords := []float64{point.Lon, point.Lat} // GeoJSON requires [lon, lat]
		data := geo.Addresses{
			Location: geo.Location{"Point", coords},
			Count:    count,
		}
		collection.Insert(data)
	}

	// start the REST server
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/addresses", geo.GetAddresses(mongoSession)).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	corsObj := handlers.AllowedOrigins([]string{"*"})

	log.Println("Listening on :8000")
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(corsObj)(router)))
}
