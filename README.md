# Opabinia ![logo](layout/assets/favicon/favicon-32x32.png)

Opabinia is simple markdown server and compiler of static sites, written on go.
I started it as the engine of [my personal website](https://artemfrolov.me).

## What's inside

- A simple web server with [`Chi`](https://github.com/go-chi/chi/v5) router for serving static files and markdown
  pages.
- [Go Markdown parser](https://github.com/gomarkdown/markdown) for pages, complemented with custom rendering extensions.
- Structured logging via [slog](https://pkg.go.dev/golang.org/x/exp/slog).
- My custom implementation of in-memory LFU cache for rendered pages.
- Configuration structure using [`gopkg.in/yaml.v3`](https://pkg.go.dev/gopkg.in/yaml.v3)
- Code formatting using [`gofumpt`](https://github.com/mvdan/gofumpt) (set as a git hook).

## Working with Opabinia

```shell
# Install opabinia engine
go install github.com/dissipative/opabinia

# Create directory with your markdown files and assets:
opabinia --init <project name>

# Start engine with
cd <project name>
opabinia --serve

# Compile static website with
opabinia --compile
```

Where to put your files:

- put entry point `index.md` to project root directory
- put other markdown files to `pages/`
- put images and other assets to `assets/`

Edit template at `templates/default.tmpl`. You can create your own and specify it in configuration file `opabinia.yml`.

You can find example of project structure and settings in `layout` directory.

## To Do

- [x] Compiler: handle resources in html tags
- [x] Links consistency for compiler and server
- [x] Migrate to Go 1.21.
- [ ] Increase test coverage.
- [ ] Implement cache TTL.
- [ ] Live reloading for serve mode
