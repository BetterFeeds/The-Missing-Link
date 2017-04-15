package sources

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nenadl/atom"
)

type Imdb struct {
}

func (imdb Imdb) CreateFeed(id string, page int) (atom.Feed, error) {
	feed := atom.Feed{Icon: "http://ia.media-imdb.com/images/G/01/imdb/images/favicon-2165806970._CB522736556_.ico",
		ID:      "tag:tml.betterfeeds.org,2017-04-02:/imdb/",
		Title:   "IMDb Trailers",
		Updated: atom.Time(time.Now())}

	feed.Link = []atom.Link{atom.Link{Href: "http://www.imdb.com/trailers#recAddTab",
		Rel:      "alternate",
		Type:     "text/html",
		HrefLang: "en-gb"},
		atom.Link{Href: "https://tml.betterfeeds.org/imdb.atom",
			Rel:      "self",
			Type:     "application/atom+xml",
			HrefLang: "en-gb"}}

	feed.Extension = []interface{}{Complete{}}

	doc, err := goquery.NewDocument("http://www.imdb.com/trailers")
	if err != nil {
		return atom.Feed{}, err
	}

	doc.Find("#recAddTab .gridlist .trailer-item").Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find(".trailer-caption a")

		title := titleElement.Text()[1:]
		title = title[:len(title)-1]

		href, _ := titleElement.Attr("href")
		imdbID := strings.Split(href, "?")[0][7:]
		videoId, _ := s.Attr("data-videoid")

		entry := atom.Entry{Title: title}

		entry.ID = "http://www.imdb.com/video/imdb/" + videoId
		entry.Link = []atom.Link{
			atom.Link{Href: entry.ID,
				Rel:      "alternate",
				Type:     "text/html",
				HrefLang: "en-gb",
				Title:    "IMDb - Trailer"},
			atom.Link{Href: "https://trakt.tv/search/imdb?query=" + imdbID,
				Rel:      "related",
				Type:     "text/html",
				HrefLang: "en-gb",
				Title:    "Trakt - " + title},
			atom.Link{Href: "http://imdb.com/title/" + imdbID,
				Rel:      "related",
				Type:     "text/html",
				HrefLang: "en-gb",
				Title:    "IMDb - " + title}}
		entry.Updated = atom.Time(time.Now()) // ... TODO ...

		entry.Extension = []interface{}{Player{
			Url:    entry.ID + "/imdb/embed",
			Height: 480,
			Width:  854}}

		feed.Entry = append(feed.Entry, entry)
	})

	return feed, nil
}

type Player struct {
	XMLName xml.Name `xml:"http://search.yahoo.com/mrss/ player"`

	Url    string `xml:"url,attr"`
	Height int    `xml:"height,attr,omitempty"`
	Width  int    `xml:"width,attr,omitempty"`
}

type Complete struct {
	XMLName xml.Name `xml:"http://purl.org/syndication/history/1.0 complete"`
}
