# go|htmx

Fun project with [Go](https://go.dev/), [Fiber](https://github.com/gofiber/fiber), and [htmx](https://htmx.org/docs/#trigger-modifiers).

Nothing in here is production ready. Copy-paste at your own risk.

# TODO
- [ ] Harden the application
  - [ ] JWT-integration
  - [ ] CSRF
  - [ ] Encrypt Cookie (probably not relevant)
  - [ ] Limiter
  - [ ] Helmet
  - [/] Idempotency
- [ ] Speed improvements
  - [ ] Compress
  - [ ] ETag
  - [ ] Server-Side Caching (IDEA)
  - [ ] CDN (IDEA)
- [ ] Interop improvements
  - [ ] CORS (IDEA)
  - [ ] JSON response from endpoints
      - [ ] GET /
      - [ ] GET /colors
      - [ ] GET /films
      - [ ] POST /films
      - [ ] DELETE /colors
      - [ ] POST /generators/films/:num<min(5);max(50)>
      - [ ] GET /directors
- [ ] Ops improvements
  - [ ] CORS (IDEA)
  - [/] Metrics
  - [ ] Pprof
  - [/] Recover
- [ ] Colors module
- [ ] Errors
  - [ ] htmx-style (polling?)
  - [ ] htmx-style (polling?)
- [ ] Support for microservice-like deployment
- [ ] Login (IDEA)

# Notes
JS goodies:
- [jQuery -> VanillaJS](https://tobiasahlin.com/blog/move-from-jquery-to-vanilla-javascript/)
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