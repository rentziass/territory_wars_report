package guild

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Guild struct {
	Name string
	URL string
	Members []*Member
}

type Member struct {
	Name string
	URL string
	GP int64
	Zetas []string
	Toons []string
}

func GetGuild(url string) (*Guild, error) {
	res, err := http.DefaultClient.Get(url + "zetas/")
	if err != nil {
		return nil, err
	}

	guild := &Guild{URL: url, Members: []*Member{}}

	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return nil, err
	}
	guildName := doc.Find("h3").First().Text()
	guildName = strings.Replace(guildName, "Zetas", "", -1)
	guildName = strings.TrimSpace(guildName)

	guild.Name = guildName


	doc.Find("tbody > tr").Each(func(i int, s *goquery.Selection) {
		memberName := strings.TrimSpace(s.Find("td").First().Text())
		memberURL, _ := s.Find("td").First().Find("a").First().Attr("href")
		memberZetas := []string{}

		// All zetas
		s.Find(".guild-member-zeta").Each(func(i int, zc *goquery.Selection) {
			toonName, _ := zc.Find(".char-portrait").First().Attr("title")
			zc.Find("img.guild-member-zeta-ability").Each(func(i int, za *goquery.Selection) {
				name, _ := za.Attr("title")
				memberZetas = append(memberZetas, toonName + " - " + name)
			})
		})

		//var memberGP int64
		//s.Find("td").Each(func(i int, s *goquery.Selection) {
		//	if i == 1 {
		//		memberGP, _ = strconv.ParseInt(s.Text(), 10, 32)
		//	}
		//})

		member := &Member{
			Name: memberName,
			URL: "https://swgoh.gg" + memberURL,
			Zetas: memberZetas,
			//GP: memberGP,
		}

		guild.Members = append(guild.Members, member)
	})

	return guild, nil
}
