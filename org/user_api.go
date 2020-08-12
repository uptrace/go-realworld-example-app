package org

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"golang.org/x/crypto/bcrypt"

	"github.com/uptrace/go-realworld-example-app/rwe"
)

var errUserNotFound = errors.New("Not Registered email or invalid password")

func setUserToken(user *User) error {
	token, err := CreateUserToken(user.ID, 24*time.Hour)
	if err != nil {
		return err
	}

	user.Token = token
	return nil
}

func currentUser(c *gin.Context) {
	user := c.MustGet("user").(*User)
	c.JSON(200, gin.H{"user": user})
}

func createUser(c *gin.Context) {
	user := new(User)
	if err := c.BindJSON(user); err != nil {
		return
	}

	var err error
	user.PasswordHash, err = hashPassword(user.Password)
	if err != nil {
		c.Error(err)
		return
	}

	if _, err = rwe.PGMain().
		ModelContext(c, user).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	if err = setUserToken(user); err != nil {
		c.Error(err)
		return
	}

	user.Password = ""
	c.JSON(200, gin.H{"user": user})
}

func hashPassword(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func loginUser(c *gin.Context) {
	var in struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&in); err != nil {
		return
	}

	user := new(User)
	if err := rwe.PGMain().
		ModelContext(c.Request.Context(), user).
		Where("email = ?", in.Email).
		Select(); err != nil {
		if err == pg.ErrNoRows {
			err = errUserNotFound
		}

		c.Error(err)
		return
	}

	if err := comparePasswords(user.PasswordHash, in.Password); err != nil {
		c.Error(err)
		return
	}

	err := setUserToken(user)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"user": user})
}

func comparePasswords(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return errUserNotFound
	}
	return nil
}

func updateUser(c *gin.Context) {
	user := new(User)
	if err := c.BindJSON(user); err != nil {
		return
	}

	var err error
	user.PasswordHash, err = hashPassword(user.Password)
	if err != nil {
		c.Error(err)
		return
	}

	authUser := c.MustGet("user").(*User)
	if _, err = rwe.PGMain().
		ModelContext(c, authUser).
		Set("email = ?", user.Email).
		Set("username = ?", user.Username).
		Set("password_hash = ?", user.PasswordHash).
		Set("image = ?", user.Image).
		Set("bio = ?", user.Bio).
		Where("id = ?", authUser.ID).
		Returning("*").
		Update(); err != nil {
		c.Error(err)
		return
	}

	user.Password = ""
	c.JSON(200, gin.H{"user": authUser})
}

func showProfile(c *gin.Context) {
	followingColumn := func(q *orm.Query) (*orm.Query, error) {
		u, _ := c.Get("user")
		authUser, ok := u.(*User)

		if !ok {
			q = q.ColumnExpr("false AS following")
		} else {
			subq := rwe.PGMain().Model((*FollowUser)(nil)).
				Where("fu.followed_user_id = u.id").
				Where("fu.user_id = ?", authUser.ID)

			q = q.ColumnExpr("EXISTS (?) AS following", subq)
		}

		return q, nil
	}

	user := new(User)
	if err := rwe.PGMain().
		ModelContext(c, user).
		ColumnExpr("u.*").
		Apply(followingColumn).
		Where("username = ?", c.Param("username")).
		Select(); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"profile": NewProfile(user)})
}

func followUser(c *gin.Context) {
	authUser := c.MustGet("user").(*User)

	user, err := SelectUserByUsername(c, c.Param("username"))
	if err != nil {
		c.Error(err)
		return
	}

	followUser := &FollowUser{
		UserID:         authUser.ID,
		FollowedUserID: user.ID,
	}
	if _, err := rwe.PGMain().
		ModelContext(c, followUser).
		Insert(); err != nil {
		c.Error(err)
		return
	}

	user.Following = true
	c.JSON(200, gin.H{"profile": NewProfile(user)})
}

func unfollowUser(c *gin.Context) {
	authUser := c.MustGet("user").(*User)

	user, err := SelectUserByUsername(c, c.Param("username"))
	if err != nil {
		c.Error(err)
		return
	}

	if _, err := rwe.PGMain().
		ModelContext(c, (*FollowUser)(nil)).
		Where("user_id = ?", authUser.ID).
		Where("followed_user_id = ?", user.ID).
		Delete(); err != nil {
		c.Error(err)
		return
	}

	user.Following = false
	c.JSON(200, gin.H{"profile": NewProfile(user)})
}
