/*
Copyright © 2020 Alex Watt <alex@alexcwatt.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package data

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Hymn struct {
	Number int
	Title  string
	Author string
}

var Hymns []Hymn
var ZipPath string

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func getMP3s() {
	zipURL := "https://www.opc.org/books/TH/TH2/MP3/Th2_MP3.zip"
	ZipPath = filepath.Join("data", "TH2_MP3.zip")
	if _, err := os.Stat(ZipPath); os.IsNotExist(err) {
		confirm := askForConfirmation("Would you like to download the MP3 audio for the hymnal? (~2.2 GB)")
		if confirm {
			_ = DownloadFile(ZipPath, zipURL)
		}
	}
}

func loadHymns() {
	// TODO: This doesn't work if the command `hymnal` is run from a different spot
	// Look into https://github.com/GeertJohan/go.rice
	csvpath := filepath.Join("data", "index.csv")
	csvfile, err := os.Open(csvpath)
	defer csvfile.Close()
	if err != nil {
		log.Fatalln("Couldn't open hymnal index")
	}

	r := csv.NewReader(csvfile)
	// read header
	if _, err := r.Read(); err != nil {
		panic(err)
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		number, _ := strconv.Atoi(record[0])
		Hymns = append(Hymns, Hymn{
			Number: number,
			Title:  record[1],
			Author: record[2],
		})
	}
}

func init() {
	loadHymns()
	getMP3s()
}
