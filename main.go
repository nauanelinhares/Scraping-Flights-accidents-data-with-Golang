package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

// Define a struct to represent an accident
type accident struct {
	Date                string
	Time                string
	Type                string
	Owner_operator      string
	Registration        string
	MSN                 string
	Year_of_manufacture string
	Fatalities          string
	Aircraft_damage     string
	Category            string
	Location            string
	Phase               string
}

// Helper function to check if a string is present in a slice of strings
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	// Create a new CSV file to store the accident data
	csvFile, err := os.Create("accidents.csv")
	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}
	defer csvFile.Close()
	w := csv.NewWriter(csvFile)

	// Write the column headers to the CSV file
	err = w.Write([]string{"Date", "Time", "Type", "Owner/operator", "Registration", "MSN", "Year of manufacture", "Fatalities", "Aircraft damage", "Category", "Location", "Phase"})

	// Define a list of columns to extract from the HTML table
	columns := []string{"Date:", "Time:", "Type:", "Owner/operator:", "Registration:", "MSN:", "Year of manufacture:", "Fatalities:", "Aircraft damage:", "Category:", "Location:", "Phase:"}

	// Create a new collector
	c := colly.NewCollector()
	accidentInfo := c.Clone()
	urlbase := "https://aviation-safety.net/database/year/"

	// On every <a> element with href attribute, call the callback function
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.HasPrefix(link, "/wikibase/") {
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)
			accidentInfo.Visit(e.Request.AbsoluteURL(link))
		}
	})

	// Before making a request, print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Extract accident information from the HTML table
	accidentInfo.OnHTML("table", func(e *colly.HTMLElement) {
		// Create a map to store the accident data
		accidentMap := make(map[string]string)

		// Iterate over each row in the table
		e.ForEach("tr", func(_ int, el *colly.HTMLElement) {
			// Extract the column and value from each row
			column := el.ChildText("td:nth-of-type(1)")
			value := el.ChildText("td:nth-of-type(2)")

			// Check if the column is in the list of columns to extract
			if stringInSlice(column, columns) {
				accidentMap[column] = value
			}
		})

		// Create an accident object using the extracted data
		accident_info := accident{
			Date:                accidentMap["Date:"],
			Time:                accidentMap["Time:"],
			Type:                accidentMap["Type:"],
			Owner_operator:      accidentMap["Owner/operator:"],
			Registration:        accidentMap["Registration:"],
			MSN:                 accidentMap["MSN:"],
			Year_of_manufacture: accidentMap["Year of manufacture:"],
			Fatalities:          accidentMap["Fatalities:"],
			Aircraft_damage:     accidentMap["Aircraft damage:"],
			Category:            accidentMap["Category:"],
			Location:            accidentMap["Location:"],
			Phase:               accidentMap["Phase:"],
		}

		// Write the accident data to the CSV file
		if accident_info.Date != "" {
			w.Write([]string{
				accident_info.Date,
				accident_info.Time,
				accident_info.Type,
				accident_info.Owner_operator,
				accident_info.Registration,
				accident_info.MSN,
				accident_info.Year_of_manufacture,
				accident_info.Fatalities,
				accident_info.Aircraft_damage,
				accident_info.Category,
				accident_info.Location,
				accident_info.Phase})
			w.Flush()
		}
	})

	// Define a list of years to scrape
	years := []string{"2023", "2024"}

	// Visit the URLs for each year and month
	for _, year := range years {
		err = c.Visit(urlbase + year + "/1")
		err = c.Visit(urlbase + year + "/2")
		err = c.Visit(urlbase + year + "/3")
		err = c.Visit(urlbase + year + "/4")
	}

	if err != nil {
		fmt.Println(err)
	}
}
