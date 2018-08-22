package geo

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Point represents a geographical latitude/longitude pair
type Point struct {
	Lat float64
	Lon float64
}

// Location is a GeoJSON location
type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"` // GeoJSON requires [lon, lat]
}

// Addresses represent the number of addresses at a Location
type Addresses struct {
	Location Location `json:"location"`
	Count    int      `json:"count"`
}

func stringsToFloat64s(s []string) (f []float64) {
	for _, str := range s {
		float, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
		if err != nil {
			log.Print(err)
			return []float64{}
		}
		f = append(f, float)
	}
	return f
}

// GetAddresses returns all address data
func GetAddresses(mongos *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var addresses []Addresses
		var query bson.M

		// get mongo db session
		session := mongos.Copy()
		defer session.Close()
		c := session.DB("ipv6").C("addresses")

		// check for "box" query parameters
		bbox, ok := r.URL.Query()["bbox"]
		if ok {
			// bbox format = min Lon, min Lat, max Lon, max Lat
			values := strings.Split(bbox[0], ",")
			floatValues := stringsToFloat64s(values)
			bottomLeft := []float64{floatValues[0], floatValues[1]}
			upperRight := []float64{floatValues[2], floatValues[3]}
			coords := [][]float64{bottomLeft, upperRight}
			query = bson.M{"location": bson.M{"$geoWithin": bson.M{"$box": coords}}}
		}

		err := c.Find(query).All(&addresses)
		if err != nil {
			log.Println("Failed get all addresses: ", err)
			return
		}

		var data [][]interface{}

		for a := range addresses {
			data = append(data, []interface{}{
				addresses[a].Location.Coordinates[1],
				addresses[a].Location.Coordinates[0],
				addresses[a].Count,
			})
		}
		json.NewEncoder(w).Encode(data)
	}
}
