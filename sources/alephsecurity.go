package sources

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nenadl/atom"
)

type AlephSecurity struct {
}

func (alephsecurity AlephSecurity) CreateFeed(id string, page int) (atom.Feed, error) {
	feed := atom.Feed{
		Logo:    "https://alephsecurity.com/favicon.png",
		ID:      "tag:tml.betterfeeds.org,2017-08-17:/alephsecurity/",
		Title:   "Aleph Security",
		Updated: atom.Time(time.Now())}

	feed.Link = []atom.Link{
		atom.Link{
			Href:     "https://alephsecurity.com/posts/",
			Rel:      "alternate",
			Type:     "text/html",
			HrefLang: "en"},
		atom.Link{
			Href:     "https://tml.betterfeeds.org/alephsecurity.atom",
			Rel:      "self",
			Type:     "application/atom+xml",
			HrefLang: "en"}}

	doc, err := goquery.NewDocument("https://alephsecurity.com/posts/")
	if err != nil {
		return atom.Feed{}, err
	}

	doc.Find(".pagelist ul li").Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("a")
		link, _ := titleElement.Attr("href")

		entry := atom.Entry{Title: titleElement.Text()}

		entry.ID = "https://alephsecurity.com" + link
		entry.Link = []atom.Link{
			atom.Link{
				Href:     entry.ID,
				Rel:      "alternate",
				Type:     "text/html",
				HrefLang: "en"},
			atom.Link{
				Href:     entry.ID + "#disqus_thread",
				Rel:      "replies",
				Type:     "text/html",
				HrefLang: "en"}}

		postDate, _ := time.Parse("02-Jan 2006", s.Find(".pagelist__date span").Text())
		entry.Updated = atom.Time(postDate)

		author := s.Find(".pagelist__subtitle a").First()
		authorLinkRelative, _ := author.Attr("href")
		authorLink := "https://alephsecurity.com" + authorLinkRelative
		entry.Author = []atom.Person{
			atom.Person{
				Name: author.Text(),
				URI:  authorLink}}

		postDoc, err := goquery.NewDocument(entry.ID)
	        if err != nil {
        	        return
	        }

		content, _ := postDoc.Find(".page__content").Html()
                entry.Content = &atom.Text{Type: "html", Body: content}

		feed.Entry = append(feed.Entry, entry)
	})

	return feed, nil
}
