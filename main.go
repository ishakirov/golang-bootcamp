package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html/charset"
)

func main() {
	currencyPtr := flag.String("currency", "RUB", "currency")
	valuePtr := flag.Int("value", 1, "value")
	flag.Parse()

	type Valute struct {
		NumCode  string
		CharCode string
		Nominal  int
		Name     string
		Value    float32
	}

	type Result struct {
		ValCurs []Valute `xml:"Valute"`
	}

	ratesURL := "http://www.cbr.ru/scripts/XML_daily.asp"

	resp, err := http.Get(ratesURL)

	if err != nil {
		log.Fatal("bad request: ", err)
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal("error reading response: ", err)
		os.Exit(1)
	}
	bodyString := string(body)
	bodyString = strings.Replace(bodyString, ",", ".", -1)
	v := Result{}
	d := xml.NewDecoder(strings.NewReader(bodyString))
	d.CharsetReader = charset.NewReaderLabel
	err = d.Decode(&v)

	if err != nil {
		log.Fatal("error parsing xml: ", err)
		os.Exit(1)
	}

	for _, curr := range v.ValCurs {
		if curr.CharCode == *currencyPtr {
			fmt.Printf("%.2f RUB\n", curr.Value*float32(*valuePtr)/float32(curr.Nominal))
			os.Exit(0)
		}
	}

	fmt.Println("currency not found")
}
