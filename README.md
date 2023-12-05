# go|htmx

Fun project with [Go](https://go.dev/), [Fiber](https://github.com/gofiber/fiber), and [htmx](https://htmx.org/docs/#trigger-modifiers).

Nothing in here is production ready. Copy-paste at your own risk.

# TODO
- [/] Delete individual movies
- [/] Check all movies on screen
- [ ] Search
  - [ ] Search movies
  - [ ] Search films
- [ ] Speed improvements
  - [ ] [Compress](https://docs.gofiber.io/api/middleware/compress)
  - [ ] [ETag](https://docs.gofiber.io/api/middleware/etag)
  - [ ] [Server-Side Caching](https://docs.gofiber.io/api/middleware/cache) (IDEA)
  - [ ] CDN (IDEA)
- [ ] Interop improvements
  - [ ] [CORS](https://docs.gofiber.io/api/middleware/cors) (IDEA)
  - [ ] JSON response from endpoints
    - [ ] GET /
    - [ ] GET /colors
    - [ ] GET /films
    - [ ] POST /films
    - [ ] DELETE /colors
    - [ ] POST /generators/films/:num<min(5);max(50)>
    - [ ] GET /directors
    - [ ] OpenAPI 3.x definition
- [ ] Ops improvements
  - [ ] Better logging ([slog](https://github.com/samber/slog-fiber))
  - [ ] Tracing
  - [ ] [CORS](https://docs.gofiber.io/api/middleware/cors) (IDEA)
  - [/] [Metrics](https://docs.gofiber.io/api/middleware/monitor)
  - [ ] [Pprof](https://docs.gofiber.io/api/middleware/pprof)
  - [/] [Recover](https://docs.gofiber.io/api/middleware/recover)
- [ ] Code quality
  - [ ] 
- [/] Colors module
- [/] Notifications
  - [/] htmx-style (SSE?)
  - [ ] htmx-style (Websockets?)
  - [ ] htmx-style (polling?)
  - [ ] Load/reload on notification
- [ ] Support for microservice-like deployment (IDEA)
- [ ] Harden the application
  - [ ] JWT-integration (IDEA)
    - [Fiber-Casbin](https://github.com/gofiber/contrib/tree/main/casbin)
    - [Casbin](github.com/casbin/casbin)
    - [zitadel](https://github.com/zitadel/zitadel-go) ([AMR](https://analyzemyrepo.com/analyze/zitadel/zitadel))
    - [hanko](https://www.hanko.io/) ([AMR](https://analyzemyrepo.com/analyze/teamhanko/hanko))
    - [authelia](https://www.authelia.com/) ([AMR](https://analyzemyrepo.com/analyze/authelia/authelia))
    - [hydra](https://github.com/ory/hydra) ([AMR](https://analyzemyrepo.com/analyze/ory/hydra))
    - [gorbac](https://github.com/mikespook/gorbac)
    - [jwt](https://github.com/golang-jwt/jwt)
    - [topaz](https://github.com/aserto-dev/topaz)
  - [ ] [CSRF](https://docs.gofiber.io/api/middleware/csrf)
  - [ ] [Encrypt Cookie](https://docs.gofiber.io/api/middleware/encryptcookie) (probably not relevant)
  - [ ] [Limiter](https://docs.gofiber.io/api/middleware/limiter)
  - [ ] [Helmet](https://docs.gofiber.io/api/middleware/helmet)
  - [/] [Idempotency](https://docs.gofiber.io/api/middleware/idempotency)
  - [ ] [Capslock integration](https://github.com/google/capslock)
  - [ ] [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)
  - [ ] [OpenSSF](https://securityscorecards.dev/)
  - [ ] Semantic Versioning
    - [ ] Automatic tagging (IDEA)
    - [ ] Easy releases (IDEA)
    - [ ] Tooling to enforce semantic versioning
  - [ ] [AnalyzeMyRepo](https://analyzemyrepo.com/analyze/teamhanko/hanko)

# Notes
JS goodies:
- https://github.com/fabiospampinato/cash#fn
- 
Go goodies:
- https://github.com/bugbytes-io/htmx-go-demo/blob/master/index.html
- https://github.com/donseba/go-htmx
- [Templating Cheatsheet](https://docs.google.com/document/d/17-eD5SO8ChKi4a4DXJq24SxOgb8AYdBeEpW9pcqj1Ok/edit)

HOT JavaScript frameworks to check:
- https://alpinejs.dev/
- https://stimulus.hotwired.dev/

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