# TODO

- Get new/edit articles page
  - /editor - new
  - /editor/slug - edit
- Get article page
  - [x] /article/slug
  - delete if owner
  - comments
  - delete comment if owner
- Use fastHttp context in sql requests
- Comments
  - create comment on article
  - delete comment from article
- Add seed data import (use gofakeit)
  - [x] add users
  - [x] add articles
  - [ ] add comments
- Articles Api
  - [ ] Filter articles
    - [ ] Add tag
    - [ ] add favorites
  - [x] Add get route
  - [x] return list of articles
  - [x] add author details
  - [x] Check if user is following article author
  - [x] move createAt out of repo, into service
    - This should allow us to manipulate created at during seeding
- Profile page
  - [ ] favorite own article?
  - [ ] unfavorite favorited articles

# Doing

- Move to templ
  - [ ] move article edit to templ
  - [x] Install
  - [x] render to page
  - [x] create render func
  - [x] move settings to templ
  - [x] move profile to templ
  - [x] move article to templ

# Done

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
- ArticlesRepo
  - [x] Unfavorite
  - [x] Insert tags
  - [x] Get popular tags
- Home
  - [x] Fix: home page should load articles
  - [x] Load tags
