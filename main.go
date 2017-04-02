package main

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nenadl/atom"
)

func ParsePage() (atom.Feed, error) {
	feed := atom.Feed{Logo: "https://www.openrightsgroup.org/assets/site/org/images/logo.png",
		Icon:    "https://www.openrightsgroup.org/assets/site/org/images/favicon.ico",
		ID:      "tag:openrightsgroup.org,2017-04-02:/blog/",
		Title:   "Open Rights Group",
		Updated: atom.Time(time.Now())}

	feed.Link = []atom.Link{atom.Link{Href: "https://www.openrightsgroup.org/blog/",
		Rel:      "alternate",
		Type:     "text/html",
		HrefLang: "en-gb"},
		atom.Link{Href: "https://tml.betterfeeds.org/org.atom",
			Rel:      "self",
			Type:     "application/atom+xml",
			HrefLang: "en-gb"}}

	doc, err := goquery.NewDocument("https://www.openrightsgroup.org/blog/")
	if err != nil {
		return atom.Feed{}, err
	}

	doc.Find(".container .post").Each(func(i int, s *goquery.Selection) {
		titleElement := s.Find("h2 a")

		entry := atom.Entry{Title: titleElement.Text()}

		entry.ID, _ = titleElement.Attr("href")
		entry.Link = []atom.Link{atom.Link{Href: entry.ID,
			Rel:      "alternate",
			Type:     "text/html",
			HrefLang: "en-gb"},
			atom.Link{Href: entry.ID + "#quip-topofcomments-qcom",
				Rel:      "replies",
				Type:     "text/html",
				HrefLang: "en-gb"}}

		postTime, _ := s.Find(".info span").Attr("datetime")
		entry.Updated = atom.TimeStr(postTime)

		authorName := s.Find(".info").Text()
		authorName = strings.Split(authorName, "|")[1]
		authorName = strings.Replace(authorName, "\n   ", "", 1)
		entry.Author = []atom.Person{atom.Person{Name: authorName}}

		content, _ := s.Find(".text").Html()
		entry.Content = &atom.Text{Type: "html", Body: content}

		feed.Entry = append(feed.Entry, entry)
	})

	return feed, nil
}

func main() {
	feed, _ := ParsePage()

	buffer, err := xml.MarshalIndent(feed, "", "	")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("\n" + xml.Header + string(buffer))
}
