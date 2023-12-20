# Real World App

## The mother of all demo apps

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/raahii/golang-grpc-realworld-example/blob/master/LICENSE)

> ### Go/Templ/HTMX/Sqlc codebase containing RealWorld examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec as much as possible while sticking to HATEOAS principals.

### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged backend application built with golang/HTMX including CRUD operations, authentication, routing, pagination, and more.

## How it works

- Using **Go** to implement realworld backend server.

  - Fiber: [fiber](https://github.com/gofiber/fiber)
  - SQLC: [sqlc](https://github.com/sqlc-dev/sqlc)
  - Templ: [templ](https://templ.guide/)
  - HTMX: [htmx](https://htmx.org/)
  - Hyperscript: [hyperscript](https://hyperscript.org/)

- Using **SQLite** to store data.
- Using **Nix Flakes** to build and deploy

## Getting started

The app listens and serves on `0.0.0.0:3000`.

Do drop into a nix shell (you will need the nix package manager installed)

```bash
nix develop
```

This will install the binaries needed as well as display useful information to get started developing.

## Project status

In Progress...

### Todo

- User Api
  - Register should return error if user already exists
  - feed should show following author articles...
- Update dates in templates
- Refactor to use structured errors
  - repo's should return raw error
  - services should return user facing errors
  - controllers should format error for templates
- Refactor to Echo
  - Fiber doesn't have a clean way to extend ctx, which is a main stay of a express like framework.
- Sqlite issue
  - Sqlite often crashes when making multiple request in parallel. I believe this might be a go-sql issu. I create a branch with libsql to see if the issue persists.
- Docker Build
  - a build to let non-nix users run server
- Seed
  - add default user by cli flag

#### Doing

- Get new/edit articles page
  - [x] edit tags
  - [ ] /editor/slug - edit
    - [ ] Load article and prefill page
    - [ ] update links to patch
  - [x] /editor - new

#### Done

- test login/register
- update header for authed user
- Settings page
  - [x] Get settings page
  - [x] make errors dynamic
  - [x] add post for form
  - [x] add logout handler
  - [x] Update settings
    - [x] image
    - [x] username
    - [x] bio
    - [x] email
    - [x] password
- Profile page
  - [x] Get profile page
  - [x] Make content live
  - [x] link to settings
  - [x] hide follow if own profile
  - [x] show favorited articles
- Home
  - [x] Fix: home page should load articles
  - [x] Load tags
- Move to templ
  - [x] move article edit to templ
  - [x] Install
  - [x] render to page
  - [x] create render func
  - [x] move settings to templ
  - [x] move profile to templ
  - [x] move article to templ
- Articles Api
  - [x] Filter articles
    - [x] by author
    - [x] Add tag
    - [x] add favorites
  - [x] Unfavorite
  - [x] Insert tags
  - [x] Get popular tags
- Views
  - [x] add banners
  - [x] add error message on htmx 500
- fix session stored alerts
- Use fastHttp context in sql requests (this caused a lot of disk io issues with sqlite 🤷)
- Add seed data import (use gofakeit)
  - [x] add comments
  - [x] add users
  - [x] add articles
  - [x] Add get route
  - [x] return list of articles
  - [x] add author details
  - [x] Check if user is following article author
  - [x] move createAt out of repo, into service
    - This should allow us to manipulate created at during seeding
- Comments
  - API
    - [x] create comment on article
    - [x] delete comment from article
    - [x] get articles by slug
  - Repo
    - [x] create
    - [x] read
      - [x] by id
      - [x] by author
      - [x] by article
    - [x] delete
  - Service
    - [x] create
    - [x] read
    - [x] delete
    - [x] mark if owner
- Profile page
  - [x] follow author
- Get article page
  - [x] delete if owner
  - [x] Follow/Unfollow author
  - [x] Favorite/Unfavorite article
  - [x] comments
    - [x] Get on load
    - [x] mark if owner
  - [x] delete comment if owner
  - [x] /article/:slug
