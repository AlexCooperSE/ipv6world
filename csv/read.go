package csv

import (
	"bufio"
	encsv "encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/alexcooperse/ipv6world/geo"
)

// ReadFile takes the path to valid GeoLite2 City Block csv file and returns
// a map of the frequency with which a location appears in the data
func ReadFile(path string) (points map[geo.Point]int) {
	points = make(map[geo.Point]int)

	// TODO: replace hard coded column numbers
	latCol := 7
	lonCol := 8

	// open file
	file, _ := os.Open(path)
	reader := encsv.NewReader(bufio.NewReader(file))

	// read line by line
	for {
		column, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Panic(error)
		}

		// read entry from lat column and convert to float64
		lat, err := strconv.ParseFloat(strings.TrimSpace(column[latCol]), 64)
		if err != nil {
			// note any parsing error and move on to the next line
			log.Print(err)
			continue
		}

		// read entry from lon column and convert to float64
		lon, err := strconv.ParseFloat(strings.TrimSpace(column[lonCol]), 64)
		if err != nil {
			// note any parsing error and move on to the next line
			log.Print(err)
			continue
		}

		pt := geo.Point{lat, lon}

		_, exists := points[pt]
		if exists {
			// location exists already in the map, increment its value
			points[pt]++
		} else {
			// location does not already exist in the map, add it with the value 1
			points[pt] = 1
		}
	}
	return points
}
