package testhelper

import (
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"github.com/uptrace/go-realworld-example-app/rwe"
)

func Extend(a, b Keys) Keys {
	res := make(Keys)
	for k, v := range a {
		res[k] = v
	}
	for k, v := range b {
		res[k] = v
	}
	return res
}

func TruncateUsersTable() {
	cmd := "TRUNCATE users, favorite_articles, follow_users"
	_, err := rwe.PGMain().Exec(cmd)
	Expect(err).NotTo(HaveOccurred())
}

func TruncateArticlesTable() {
	cmd := "TRUNCATE articles, favorite_articles, article_tags"
	_, err := rwe.PGMain().Exec(cmd)
	Expect(err).NotTo(HaveOccurred())
}
