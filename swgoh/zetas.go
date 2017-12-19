package swgoh

import (
	"net/http"
	"github.com/PuerkitoBio/goquery"
)

var zetaURL = "https://swgoh.gg/zeta-report/"

func GetGameZetas() ([]string, error) {
	res, err := http.DefaultClient.Get(zetaURL)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}

	zetas := []string{}

	doc.Find(".media.list-group-item.p-0.character").Each(func(i int, s *goquery.Selection) {
		name := s.Find("h5").First().Text()
		zetas = append(zetas, name)
	})

	return zetas, nil
}
