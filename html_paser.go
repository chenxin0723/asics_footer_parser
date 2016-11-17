package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"os"
	"strings"
)

var (
	EcAsics = []string{"http://www.asics.com/gb/en-gb", "http://www.asics.com/fr/fr-fr", "http://www.asics.com/es/es-es", "http://www.asics.com/it/it-it"}

	NoEcAsics = []string{"https://www.asics.co.za/", "https://www.asics.pl", "https://de.asics.ch"}
)

type ParseData struct {
	Locale string
	Header []Link
	Footer []Link
}

type Link struct {
	Title string
	Url   string
	Index int
	Img   string
	Items []ItemLink
}

type ItemLink struct {
	Title string
	Url   string
	Index int
	Img   string
}

func ParseECHtml(url_string string) {

	target_url, _ := url.Parse(url_string)
	doc, err := goquery.NewDocument(url_string)
	if err != nil {
		log.Fatal(err)
	}

	parsedata := ParseData{
		Locale: strings.Split(target_url.Path, "/")[2],
	}

	//parse the header
	// doc.Find("header #main-menu").Children().Each(func(i int, s *goquery.Selection) {

	//parse the footer
	doc.Find("footer").Children().Filter(".tiger-clearfix-toggle").Each(func(i int, s *goquery.Selection) {
		link := Link{
			Title: s.Children().Filter("h4").First().Text(),
			Url:   "#",
			Index: i + 1,
		}

		s.Children().Filter("ul").Children().Each(func(i int, s *goquery.Selection) {
			itemlink := ItemLink{
				Title: s.Children().First().Text(),
				Url:   target_url.Scheme + "://" + target_url.Host + s.Children().First().AttrOr("href", "#"),
				Index: i + 1,
			}
			link.Items = append(link.Items, itemlink)
		})
		parsedata.Footer = append(parsedata.Footer, link)
	})

	jsondata, _ := json.Marshal(parsedata)

	file, _ := os.OpenFile(strings.Split(target_url.Path, "/")[2]+".json", os.O_WRONLY|os.O_CREATE, 0666)
	file.Write(jsondata)

	fmt.Println(string(jsondata))

}

func ParseNoECHtml(url_string string) {

	target_url, _ := url.Parse(url_string)
	doc, err := goquery.NewDocument(url_string)
	if err != nil {
		log.Fatal(err)
	}

	locale := strings.Split(target_url.Host, ".")[0] + "-" + strings.Split(target_url.Host, ".")[2]
	parsedata := ParseData{
		Locale: locale,
	}

	//parse the footer
	doc.Find(".footer #tertiary section").Each(func(i int, s *goquery.Selection) {
		link := Link{
			Title: s.Children().First().Text(),
			Url:   "#",
			Index: i + 1,
		}
		fmt.Println(s)

		s.Children().Filter("ul").Children().Each(func(i int, s *goquery.Selection) {
			itemlink := ItemLink{
				Title: s.Children().First().Text(),
				Url:   s.Children().First().AttrOr("href", "#"),
				Index: i + 1,
			}
			link.Items = append(link.Items, itemlink)
		})
		parsedata.Footer = append(parsedata.Footer, link)
	})

	jsondata, _ := json.Marshal(parsedata)

	file, _ := os.OpenFile(locale, os.O_WRONLY|os.O_CREATE, 0666)
	file.Write(jsondata)

	fmt.Println(string(jsondata))

}

func main() {
	for _, url := range EcAsics {
		ParseECHtml(url)
	}
	for _, url := range NoEcAsics {
		ParseNoECHtml(url)
	}
}
