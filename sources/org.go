package sources

import (
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/nenadl/atom"
)

type Org struct {
}

func (org Org) CreateFeed(id string, page int) (atom.Feed, error) {
	feed := atom.Feed{Logo: "https://www.openrightsgroup.org/assets/site/org/images/logo.png",
		Icon:    "https://www.openrightsgroup.org/assets/site/org/images/favicon.ico",
		ID:      "tag:openrightsgroup.org,2017-04-02:/blog/",
		Title:   "Open Rights Group",
		Updated: atom.Time(time.Now())}

	feed.Link = []atom.Link{atom.Link{Href: "https://www.openrightsgroup.org/blog/?page=" + strconv.Itoa(page),
		Rel:      "alternate",
		Type:     "text/html",
		HrefLang: "en-gb"},
		atom.Link{Href: "https://tml.betterfeeds.org/org/org-" + strconv.Itoa(page) + ".atom",
			Rel:      "self",
			Type:     "application/atom+xml",
			HrefLang: "en-gb"}}

	doc, err := goquery.NewDocument("https://www.openrightsgroup.org/blog/?page=" + strconv.Itoa(page))
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

	lastPage, _ := doc.Find(".paging .pageList li").Last().Find("a").Attr("href")
	lastPageNo := strings.Split(lastPage, "=")[1]
	lastPageNoInt, _ := strconv.Atoi(lastPageNo)

	feed.Link = append(feed.Link, atom.Link{Href: "https://tml.betterfeeds.org/org/org-1.atom",
		Rel:      "first",
		Type:     "application/atom+xml",
		HrefLang: "en-gb"}, atom.Link{Href: "https://tml.betterfeeds.org/org/org-" + lastPageNo + ".atom",
		Rel:      "last",
		Type:     "application/atom+xml",
		HrefLang: "en-gb"})

	if page > 1 {
		feed.Link = append(feed.Link, atom.Link{Href: "https://tml.betterfeeds.org/org/org-" +
			strconv.Itoa(page-1) + ".atom",
			Rel:      "previous",
			Type:     "application/atom+xml",
			HrefLang: "en-gb"})
	}
	if page < lastPageNoInt {
		feed.Link = append(feed.Link, atom.Link{Href: "https://tml.betterfeeds.org/org/org-" +
			strconv.Itoa(page+1) + ".atom",
			Rel:      "next",
			Type:     "application/atom+xml",
			HrefLang: "en-gb"})
	}

	return feed, nil
}
