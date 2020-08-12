package testhelper

import (
	"github.com/uptrace/go-realworld-example-app/rwe"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

func ExtendKeys(a, b Keys) Keys {
	res := make(Keys)
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}

func TruncateDB() {
	cmds := []string{
		"TRUNCATE users, favorite_articles, follow_users, comments",
		"TRUNCATE articles, article_tags, favorite_articles, comments",
	}

	for _, cmd := range cmds {
		_, err := rwe.PGMain().Exec(cmd)
		Expect(err).NotTo(HaveOccurred())
	}
}
