package repositories

import (
	"bufio"
	"encoding/csv"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Data struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	City      string
}

func (d *Data) ToString() string {
	return strings.Join([]string{d.ID, d.FirstName, d.LastName, d.Email, d.City}, ",")
}

func CommaStringToData(s string) Data {
	row := strings.Split(s, ",")
	return Data{
		ID:        row[0],
		FirstName: row[1],
		LastName:  row[2],
		Email:     row[3],
		City:      row[4],
	}
}

func LoadCSV() []Data {
	dataset := make([]Data, 0)

	// Load CSV file and populate the 'dataset' slice
	csvFile, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	records = records[1:]
	records = Shuffle(records)

	for _, rec := range records {
		dataset = append(dataset, CommaStringToData(rec[0]))
	}

	return dataset
}

func Shuffle(data [][]string) [][]string {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
	return data
}
