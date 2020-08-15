# Go Gin + go-pg realworld example application

[![CircleCI](https://circleci.com/gh/uptrace/go-realworld-example-app.svg?style=svg)](https://circleci.com/gh/uptrace/go-realworld-example-app)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/uptrace/go-realworld-example-app)](https://pkg.go.dev/github.com/uptrace/go-realworld-example-app)

## Introduction

This project was created to demonstrate how to use:

- [Gin Web framework](https://github.com/gin-gonic/gin).
- [go-pg PostgreSQL client and ORM](https://github.com/go-pg/pg).
- [go-pg/migrations](https://github.com/go-pg/migrations).

It implements JSON API as specified in [RealWorld](https://github.com/gothinkster/realworld) spec.

## Project structure

Project consists of following packages:

- [rwe](rwe) global package parses configs, establishes DB connections etc.
- [org](org) package manages users and tokens.
- [blog](blog) package manages articles and comments.
- [app](app) folder contains application resources such as config.
- [cmd/api](cmd/api) runs HTTP server with JSON API.
- [cmd/migrate_db](cmd/migrate_db) command that runs SQL migrations.

The most interesting part for go-pg users is probably [article filter](blog/article_filter.go).

## Project bootstrap

First of all you need to create a config file changing defaults as needed:

```
cp app/config/dev.yml.default app/config/dev.yml
```

Project comes with a `Makefile` that contains following recipes:

- `make db_reset` drops existing database and creates a new one.
- `make test` runs unit tests.
- `make api_test` runs API tests provided by RealWorld.

After checking that tests are passing you can start API HTTP server:

```shell
go run cmd/api/*.go -env=dev
```
