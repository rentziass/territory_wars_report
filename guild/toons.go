package guild

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"github.com/StefanSchroeder/Golang-Roman"
	"strconv"
	"strings"
)

type Toon struct {
	Name string
	GearLevel int
	GP int
	Stars int
	Level int
}

func (m *Member) GetToonRoster() error {
	m.Toons = []*Toon{}

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

			stars := 0
			s.Find(".star").Each(func(i int, star *goquery.Selection) {
				if !(star.HasClass("star-inactive")) {
					stars++
				}
			})

			romanGearLevel := s.Find(".char-portrait-full-gear-level").First().Text()
			gearLevel := roman.Arabic(romanGearLevel)

			stringLevel := s.Find(".char-portrait-full-level").First().Text()
			level, _ := strconv.ParseInt(stringLevel, 10, 64)

			gpTitle, _ := s.Find(".collection-char-gp").First().Attr("title")
			gpParts := strings.Split(gpTitle, "/")
			gpString := strings.Replace(gpParts[0], "Power", "", 1)
			gpString = strings.Replace(gpString, ",", "", 1)
			gpString = strings.TrimSpace(gpString)
			gpInt, _ := strconv.ParseInt(gpString, 10, 64)

			toon := &Toon{
				Name: name,
				GearLevel: gearLevel,
				Level: int(level),
				Stars: stars,
				GP: int(gpInt),
			}
			m.Toons = append(m.Toons, toon)
		}
	})

	return nil
}
