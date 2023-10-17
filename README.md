# Real World App

## The mother of all demo apps

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/raahii/golang-grpc-realworld-example/blob/master/LICENSE)

> ### Go/Templ/HTMX/Sqlc codebase containing RealWorld examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec as much as possible while sticking to HATEOAS principals.

### [Demo](https://github.com/gothinkster/realworld)&nbsp;&nbsp;&nbsp;&nbsp;[RealWorld](https://github.com/gothinkster/realworld)

This codebase was created to demonstrate a fully fledged fullstack application built with golang/HTMX including CRUD operations, authentication, routing, pagination, and more.

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

See [kanban](./kanban.md)
