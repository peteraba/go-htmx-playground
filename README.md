# go|htmx

Fun project with [Go](https://go.dev/), [Fiber](https://github.com/gofiber/fiber), and [htmx](https://htmx.org/docs/#trigger-modifiers).

Nothing in here is production ready. Copy-paste at your own risk.

## Installation instructions

### Install Zitadel

Follow one of the following guides:

- [k8s](https://zitadel.com/docs/self-hosting/deploy/kubernetes)
- [docker](https://zitadel.com/docs/self-hosting/deploy/docker)
- [linux](https://zitadel.com/docs/self-hosting/deploy/linux)

### Install go-htmx

```bash
git clone https://github.com/peteraba/go-htmx-playground
cd go-htmx-playground
sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin
task install
```

## Tools being used in this project

- [Task](https://taskfile.dev/) - Task runner (like Make, but nicer)
- [Air](https://github.com/cosmtrek/air) - rebuilding project on file change
- 

## TODO

- [ ] Project health
  - [ ] Badges
  - [ ] Process to upgrade tools
  - [ ] Process to upgrade vendor libraries
- [x] Fix bugs
  - [x] Fix back buttons
  - [x] Fix JS being lost on re-loading parts
  - [x] Fix check-all remaining hidden on generating / adding movies
  - [x] Fix search design
  - [x] Make sure search is highlighted on reload / load only the table???
- [ ] Speed improvements
  - [ ] [Compress](https://docs.gofiber.io/api/middleware/compress)
  - [ ] [ETag](https://docs.gofiber.io/api/middleware/etag)
  - [ ] [Server-Side Caching](https://docs.gofiber.io/api/middleware/cache) (IDEA)
  - [ ] ~~CDN (IDEA)~~
- [ ] Interop improvements
  - [ ] [CORS](https://docs.gofiber.io/api/middleware/cors) (IDEA)
  - [ ] JSON response from endpoints
    - [x] GET /movies
    - [x] POST /movies
    - [ ] DELETE /movies
    - [ ] POST /film-generators
    - [x] GET /directors
    - [ ] GET /colors
    - [ ] OpenAPI 3.x definition
- [ ] Ops improvements
  - [x] Better logging ([slog](https://github.com/samber/slog-fiber))
  - [ ] Tracing
  - [ ] [CORS](https://docs.gofiber.io/api/middleware/cors) (IDEA)
  - [x] [Metrics](https://docs.gofiber.io/api/middleware/monitor)
  - [ ] [Pprof](https://docs.gofiber.io/api/middleware/pprof)
  - [x] [Recover](https://docs.gofiber.io/api/middleware/recover)
- [x] Code quality
  - [x] golangci-lint
  - [x] Refactoring for more modular code
- [x] Colors module
- [x] Notifications
  - [x] htmx-style (SSE?)
  - [ ] ~~htmx-style (Websockets?)~~
  - [ ] ~~htmx-style (polling?)~~
  - [ ] ~~Load/reload on notification~~
- [ ] Support for microservice-like deployment (IDEA)
- [x] Arch changes
  - [x] Use templ
  - [x] Use urfave/cli to support multiple commands
  - [x] Use alpineJS or similar
- [ ] Harden the application
  - [ ] Testing race conditions
  - [ ] Fuzz testing
  - [ ] Performance testing
  - [ ] Authorization
    - [Fiber-Casbin](https://github.com/gofiber/contrib/tree/main/casbin)
    - [Casbin](https://github.com/casbin/casbin)
    - [gorbac](https://github.com/mikespook/gorbac)
    - [topaz](https://github.com/aserto-dev/topaz)
  - [x] JWT-integration (IDEA)
    - [zitadel](https://github.com/zitadel/zitadel-go) ([AMR](https://analyzemyrepo.com/analyze/zitadel/zitadel))
    - [hanko](https://www.hanko.io/) ([AMR](https://analyzemyrepo.com/analyze/teamhanko/hanko))
    - [authelia](https://www.authelia.com/) ([AMR](https://analyzemyrepo.com/analyze/authelia/authelia))
    - [hydra](https://github.com/ory/hydra) ([AMR](https://analyzemyrepo.com/analyze/ory/hydra))
    - [jwt](https://github.com/golang-jwt/jwt)
  - [ ] [CSRF](https://docs.gofiber.io/api/middleware/csrf)
  - [ ] [Encrypt Cookie](https://docs.gofiber.io/api/middleware/encryptcookie) (probably not relevant)
  - [ ] [Limiter](https://docs.gofiber.io/api/middleware/limiter)
  - [ ] [Helmet](https://docs.gofiber.io/api/middleware/helmet)
  - [x] [Idempotency](https://docs.gofiber.io/api/middleware/idempotency)
  - [x] [Capslock integration](https://github.com/google/capslock)
  - [x] [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
  - [ ] [OpenSSF](https://securityscorecards.dev/)
  - [ ] Semantic Versioning
    - [ ] Automatic tagging (IDEA)
    - [ ] Easy releases (IDEA)
    - [ ] Tooling to enforce semantic versioning
  - [ ] [AnalyzeMyRepo](https://analyzemyrepo.com/analyze/teamhanko/hanko)
  - [x] Integrity sums added for assets

# Notes
JS goodies:
- https://github.com/fabiospampinato/cash#fn

Go goodies:
- https://github.com/bugbytes-io/htmx-go-demo/blob/master/index.html
- https://github.com/donseba/go-htmx
- [Templating Cheatsheet](https://docs.google.com/document/d/17-eD5SO8ChKi4a4DXJq24SxOgb8AYdBeEpW9pcqj1Ok/edit)

HOT JavaScript frameworks to check:
- https://alpinejs.dev/start-here
- https://unpoly.com/

HOT CSS frameworks to check:
- https://daisyui.com/
- https://ui.shadcn.com/
- https://cirrus-ui.netlify.app/
- https://www.radix-ui.com/
- https://picocss.com/
- https://www.patternfly.org/
- https://fomantic-ui.com/
- https://bulma.io/
- https://mantine.dev/
- https://open-props.style/
- https://vanillaframework.io/
- https://picturepan2.github.io/
-
- https://github.com/uikit/uikit
- https://tailwindcss.com/
- https://jenil.github.io/chota/
- https://purecss.io/
- https://tachyons.io/
- https://milligram.io/
- https://watercss.kognise.dev/
- https://andybrewer.github.io/mvp/
- https://www.blazeui.com/

Go stuff to check:
- [TODO app done with Go + HTMX](https://github.com/paganotoni/todox/tree/main)

Interesting HTMX plugins:
- [loading-states](https://htmx.org/extensions/loading-states/) - helps showing/hiding load indicators
- [multi-swap](https://htmx.org/extensions/multi-swap/) - helps loading replacing multiple parts of the page
- [response-targets](https://htmx.org/extensions/response-targets/) - helps displaying error responses in different elements from success responses
- [web-sockets](https://htmx.org/extensions/web-sockets/)
- [server-sent-events](https://htmx.org/extensions/server-sent-events/)

Alternative to HTMX:
- https://hotwired.dev/
- https://vreco.fly.dev/blog/post/HTMX%20in%20Go

Tooling:
- go install github.com/cosmtrek/air@latest (watch)
- sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin