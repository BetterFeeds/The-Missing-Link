package sources

import "github.com/nenadl/atom"

type Source interface {
	CreateFeed(id string, page int) (atom.Feed, error)
}
