// Package engine defines interfaces each supported storage should implement.
// Includes default implementation with boltdb
package engine

import (
	"sort"
	"strings"

	"github.com/umputun/remark/app/store"
)

//go:generate sh -c "mockery -inpkg -name Interface -print > file.tmp && mv file.tmp engine_mock.go"

// Interface combines all store interfaces
type Interface interface {
	Accessor
	Admin
}

// Accessor defines all usual access ops avail for regular user
type Accessor interface {
	Create(comment store.Comment) (commentID string, err error)                 // create new comment, avoid dups by id
	Get(locator store.Locator, commentID string) (store.Comment, error)         // get comment by id
	Put(locator store.Locator, comment store.Comment) error                     // update comment, mutable parts only
	Find(locator store.Locator, sort string) ([]store.Comment, error)           // find comments for locator
	Last(siteID string, limit int) ([]store.Comment, error)                     // last comments for given site, sorted by time
	User(siteID string, userID string, limit int) ([]store.Comment, int, error) // comments by user, sorted by time
	Count(locator store.Locator) (int, error)                                   // number of comments for the post
	List(siteID string, limit int, skip int) ([]store.PostInfo, error)          // list of commented posts
}

// Admin defines all store ops avail for admin only
type Admin interface {
	Delete(locator store.Locator, commentID string) error     // delete comment by id
	DeleteAll(siteID string) error                            // delete all data from site
	SetBlock(siteID string, userID string, status bool) error // block or unblock  user
	IsBlocked(siteID string, userID string) bool              // check if user blocked
	Blocked(siteID string) ([]store.BlockedUser, error)       // get list of blocked users
}

// sortComments is for engines can't sort data internally
func sortComments(comments []store.Comment, sortFld string) []store.Comment {
	sort.Slice(comments, func(i, j int) bool {
		switch sortFld {
		case "+time", "-time", "time", "+active", "-active", "active":
			if strings.HasPrefix(sortFld, "-") {
				return comments[i].Timestamp.After(comments[j].Timestamp)
			}
			return comments[i].Timestamp.Before(comments[j].Timestamp)

		case "+score", "-score", "score":
			if strings.HasPrefix(sortFld, "-") {
				if comments[i].Score == comments[j].Score {
					return comments[i].Timestamp.Before(comments[j].Timestamp)
				}
				return comments[i].Score > comments[j].Score
			}
			if comments[i].Score == comments[j].Score {
				return comments[i].Timestamp.Before(comments[j].Timestamp)
			}
			return comments[i].Score < comments[j].Score

		default:
			return comments[i].Timestamp.Before(comments[j].Timestamp)
		}
	})
	return comments
}
