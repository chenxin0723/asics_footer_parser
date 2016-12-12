package main

import (
	"github.com/chenxin0723/asics_footer_parser/asics_parser"
)

func main() {
	asics_parser.WirteDir = "./asics_parser/"

	asics_parser.ParseECHtml("http://www.asics.com/us/en-us/")

}
