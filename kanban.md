# TODO

- Get new/edit articles page
  - /editor - new
  - /editor/slug - edit
- Get article page
  - /article/slug
  - delete if owner
  - comments
  - delete comment if owner
- Comments
  - create comment on article
  - delete comment from article

# Doing

- Articles Api
  - [x] Add get route
  - [x] return list of articles
  - [ ] Check if user is following article author
  - [x] Filter articles
- Add seed data import (use gofakeit)
  - [x] add users
  - [x] add articles
  - [ ] add comments
- Profile page
  - [ ] favorite own article?
  - [ ] unfavorite favorited articles

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
