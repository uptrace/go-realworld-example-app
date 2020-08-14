package blog

import (
	"encoding/json"
	"time"

	"github.com/uptrace/go-realworld-example-app/org"
)

type Comment struct {
	tableName struct{} `pg:"comments,alias:c"`

	ID   uint64 `json:"id"`
	Body string `json:"body"`

	Author   *org.Profile `json:"author"`
	AuthorID uint64       `json:"-"`

	ArticleID uint64 `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Comment) MarshalJSON() ([]byte, error) {
	type Alias Comment

	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(c),
		CreatedAt: c.CreatedAt.UTC().Format(TimeFormatStr),
		UpdatedAt: c.UpdatedAt.UTC().Format(TimeFormatStr),
	})
}
