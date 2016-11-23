package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

var (
	EuropeEcAsics = []string{"http://www.asics.com/gb/en-gb", "http://www.asics.com/fr/fr-fr", "http://www.asics.com/es/es-es", "http://www.asics.com/it/it-it"}
	NoEcAsics     = []string{"https://www.asics.co.za/", "https://www.asics.pl", "https://de.asics.ch"}
)

type ParseData struct {
	Locale string
	Header []Link
	Footer []Link
}

type Link struct {
	Title   string
	Url     string
	Name    string
	Index   int
	Img     string `json:"Img,omitempty"`
	ImgSize string `json:"ImgSize,omitempty"`
	//the value of imgsize represent the type of the image, small | medium | big
	SingleRow bool
	//when the value of SingleRow is false, there will have two link in a row.
	Items []interface{}
}

type ItemLink struct {
	Title   string
	Url     string
	Name    string
	Index   int
	Img     string `json:"Img,omitempty"`
	ImgSize string `json:"ImgSize,omitempty"`
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
	doc.Find("#main-menu").Children().Not(".mobile").Not("div").Each(func(i int, s *goquery.Selection) {
		link := Link{
			Title:     s.Children().First().AttrOr("title", ""),
			Url:       paser_url(target_url, s.Children().First().AttrOr("href", "#")),
			Name:      s.Children().First().Text(),
			Index:     i + 1,
			SingleRow: true,
		}

		s.Find("ul").Each(func(i int, s *goquery.Selection) {
			if s.Children().Size() == 2 {

				firstlink := Link{
					Title:     s.Find("h5").First().AttrOr("title", ""),
					Url:       "#",
					Name:      s.Find("h5").First().Text(),
					Index:     1,
					SingleRow: false,
				}

				secondlink := Link{
					Title:     s.Find(".empty-nav-item").Last().AttrOr("title", ""),
					Url:       "#",
					Name:      s.Find(".empty-nav-item").Last().Text(),
					Index:     2,
					SingleRow: false,
				}

				is_first := true
				s.Find(".yCmsComponent").Each(func(i int, s *goquery.Selection) {
					// fmt.Println(i)
					// fmt.Println(s.Children().First().Html())
					if s.Children().First().Is("span") {
						if s.Children().First().Text() != "" {
							is_first = false
						}

					} else {
						itemlink := ItemLink{
							Title: s.Children().First().AttrOr("title", ""),
							Name:  s.Children().First().Text(),
							Url:   paser_url(target_url, s.Children().First().AttrOr("href", "#")),
							Index: i + 1,
						}
						if s.Children().First().Children().First().Is("img") {
							itemlink.Img = paser_url(target_url, s.Children().First().Children().First().AttrOr("src", "#"))
						}

						if is_first {
							firstlink.Items = append(firstlink.Items, itemlink)
						} else {
							secondlink.Items = append(secondlink.Items, itemlink)

						}
					}
				})
				link.Items = append(link.Items, firstlink)
				link.Items = append(link.Items, secondlink)

			} else {

				firstlink := Link{
					Title:     s.Find("h5").First().AttrOr("title", ""),
					Url:       "#",
					Name:      s.Find("h5").First().Text(),
					Index:     i + 1,
					SingleRow: true,
				}
				s.Find(".yCmsComponent").Each(func(i int, s *goquery.Selection) {

					itemlink := ItemLink{
						Title: s.Children().First().AttrOr("title", ""),
						Name:  s.Children().First().Text(),
						Url:   paser_url(target_url, s.Children().First().AttrOr("href", "#")),
						Index: i + 1,
					}
					if s.Children().First().Children().First().Is("img") {
						itemlink.Img = paser_url(target_url, s.Children().First().Children().First().AttrOr("src", "#"))
					}

					firstlink.Items = append(firstlink.Items, itemlink)

				})
				link.Items = append(link.Items, firstlink)
			}
		})
		parsedata.Header = append(parsedata.Header, link)
	})

	//parse the footer
	doc.Find("footer").Children().Filter(".tiger-clearfix-toggle").Each(func(i int, s *goquery.Selection) {
		link := Link{
			Title:     s.Children().Filter("h4").First().AttrOr("title", ""),
			Url:       "#",
			Name:      s.Children().Filter("h4").First().Text(),
			Index:     i + 1,
			SingleRow: true,
		}

		s.Children().Filter("ul").Children().Each(func(i int, s *goquery.Selection) {
			itemlink := ItemLink{
				Title: s.Children().First().AttrOr("title", ""),
				Name:  s.Children().First().Text(),
				Url:   paser_url(target_url, s.Children().First().AttrOr("href", "#")),
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

	//parse the header
	doc.Find("header #asicsAreas").Children().Each(func(i int, s *goquery.Selection) {
		link := Link{
			Title:     s.Children().First().AttrOr("title", ""),
			Name:      s.Children().First().Children().First().Text(),
			SingleRow: true,
			Url:       "#",
			Index:     i + 1,
		}
		if i != 3 {
			doc.Find("header #asicsPanels").Children().Slice(i, i+1).Children().First().Children().Filter(".asicsListing").Each(func(ii int, ss *goquery.Selection) {
				secend_link := Link{
					Title:     doc.Find("header #asicsPanels").Children().Slice(i, i+1).Children().First().Children().Filter(".asicsFeatured").Children().First().Children().Slice(ii, ii+1).Children().End().AttrOr("title", ""),
					Name:      strings.TrimSpace(doc.Find("header #asicsPanels").Children().Slice(i, i+1).Children().First().Children().Filter(".asicsFeatured").Children().First().Children().Slice(ii, ii+1).Children().End().Text()),
					Url:       doc.Find("header #asicsPanels").Children().Slice(i, i+1).Children().First().Children().Filter(".asicsFeatured").Children().First().Children().Slice(ii, ii+1).Children().End().AttrOr("href", "#"),
					Img:       paser_url(target_url, doc.Find("header #asicsPanels").Children().Slice(i, i+1).Children().First().Children().Filter(".asicsFeatured").Children().First().Children().Slice(ii, ii+1).Children().First().Children().First().AttrOr("src", "#")),
					ImgSize:   "small",
					Index:     i + 1,
					SingleRow: true,
				}
				ss.Children().Filter("ul").Children().Each(func(iii int, s *goquery.Selection) {
					itemlink := ItemLink{
						Title: s.Children().First().AttrOr("title", ""),
						Url:   paser_url(target_url, s.Children().First().AttrOr("href", "#")),
						Name:  s.Children().First().Text(),
						Index: iii + 1,
					}

					secend_link.Items = append(secend_link.Items, itemlink)
				})

				link.Items = append(link.Items, secend_link)

			})

		}

		parsedata.Header = append(parsedata.Header, link)
	})

	//parse the footer
	doc.Find(".footer #tertiary section").Each(func(i int, s *goquery.Selection) {
		link := Link{
			Title:     s.Children().First().AttrOr("title", ""),
			Name:      s.Children().First().Text(),
			Url:       "#",
			Index:     i + 1,
			SingleRow: true,
		}

		s.Children().Filter("ul").Children().Each(func(i int, s *goquery.Selection) {
			itemlink := ItemLink{
				Title: s.Children().First().AttrOr("title", ""),
				Url:   paser_url(target_url, s.Children().First().AttrOr("href", "#")),
				Name:  s.Children().First().Text(),
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

func paser_url(target_url *url.URL, url string) (return_url string) {

	if matched, err := regexp.Match(`http`, []byte(url)); matched && err == nil {
		return_url = url
	} else {
		return_url = target_url.Scheme + "://" + target_url.Host + url
	}
	return
}

func main() {
	for _, url := range EuropeEcAsics {
		ParseECHtml(url)
	}
	for _, url := range NoEcAsics {
		ParseNoECHtml(url)
	}

}
