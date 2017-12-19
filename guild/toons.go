package guild

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func (m *Member) GetToonRoster() error {
	m.Toons = []string{}

	res, err := http.DefaultClient.Get(m.URL + "collection/")
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return err
	}
	doc.Find(".collection-char").Each(func(i int, s *goquery.Selection) {
		if !(s.HasClass("collection-char-missing")) {
			name := s.Find(".collection-char-name-link").First().Text()
			m.Toons = append(m.Toons, name)
		}
	})

	return nil
}
