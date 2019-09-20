package db

import (
	"ChiliOverFlow/pkg/cache"
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
)

func ReadCSV() []Chili {
	log.Print("Reading csv")
	csvFile, _ := os.Open("ChiliDB.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '|'
	var chilis []Chili
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		c := Chili{
			Name:              line[0],
			Origin:            line[3],
			Family:            line[2],
			Heat:              line[13]+" ("+line[14]+")",
			Info:              line[4],
			Taste: line[5]+"\n"+line[6],
			Color:             line[7],
			Size:              line[8],
			Fruit:             line[9],
			Plant_Form:        line[10],
			Plant_Size:        line[11],
			Guide:             "Time: " + line[15] + "\n" + "Temperature: " + line[16],
			Seeds_Per_Request: 2,
		}

		if(line[1] != "") {
			target := strings.Replace(c.Name, "/", "_", -1)+".img"
			err := cache.Image("https://"+line[1], target)

			if err == nil {
				c.Image = target;
			}
		}
		chilis = append(chilis, c)
	}
	return chilis
}
