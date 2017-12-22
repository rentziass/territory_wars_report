package guild

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"net/http"
	"strings"
)

type Ship struct {
	Name string
	GP int
	Stars int
	Level int
}

func (m *Member) GetShipRoster() error {
	m.Ships = []*Ship{}

	res, err := http.DefaultClient.Get(m.URL + "ships/")
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return err
	}
	doc.Find(".collection-ship").Each(func(i int, s *goquery.Selection) {
		if !(s.HasClass("collection-ship-missing")) {
			name := s.Find(".collection-ship-name-link").First().Text()

			stars := 0
			s.Find(".ship-portrait-full-star").Each(func(i int, star *goquery.Selection) {
				if !(star.HasClass("ship-portrait-full-star-inactive")) {
					stars++
				}
			})

			stringLevel := s.Find(".ship-portrait-full-frame-level").First().Text()
			level, _ := strconv.ParseInt(stringLevel, 10, 64)

			gpTitle, _ := s.Find(".collection-char-gp").First().Attr("title")
			gpParts := strings.Split(gpTitle, "/")
			gpString := strings.Replace(gpParts[0], "Power", "", 1)
			gpString = strings.Replace(gpString, ",", "", 1)
			gpString = strings.TrimSpace(gpString)
			gpInt, _ := strconv.ParseInt(gpString, 10, 64)

			ship := &Ship{
				Name: name,
				Level: int(level),
				Stars: stars,
				GP: int(gpInt),
			}
			m.Ships = append(m.Ships, ship)
		}
	})

	return nil
}
