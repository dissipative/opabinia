# Opabinia ![logo](layout/assets/favicon/favicon-32x32.png)

Opabinia is a simple markdown server and a compiler for static sites, written in Go.
I created it as the engine for [my personal website](https://artemfrolov.me).

## What's Inside

- A simple web server with the [`Chi`](https://github.com/go-chi/chi) router for serving static files and markdown
  pages.
- [Go Markdown parser](https://github.com/gomarkdown/markdown) for pages, enhanced with custom rendering extensions.
- Structured logging via [slog](https://pkg.go.dev/golang.org/x/exp/slog).
- My custom implementation of an in-memory LFU cache for rendered pages.
- Configuration structure using [`gopkg.in/yaml.v3`](https://pkg.go.dev/gopkg.in/yaml.v3).
- Code formatting with [`gofumpt`](https://github.com/mvdan/gofumpt), set as a git hook.

## Working with Opabinia

```shell
# Install the opabinia engine
go install github.com/dissipative/opabinia

# Create a directory for your markdown files and assets:
opabinia --init <project name>

# Start the engine with:
cd <project name>
opabinia --serve

# Compile a static website with:
opabinia --compile
```

Where to put your files:

- Put entry point `index.md` to project root directory.
- Put other markdown files to `pages/`.
- Put images and other assets to `assets/`.

Edit the template at `templates/default.tmpl`. You can also create your own one and specify it in the configuration
file `opabinia.yml`.

You can find an example of the project structure and settings in the `layout` directory.

## To Do

- [x] Compiler: handle resources in HTML tags.
- [x] Ensure link consistency for both the compiler and server.
- [x] Migrate to Go 1.21.
- [ ] Increase test coverage.
- [ ] Implement cache TTL.
- [ ] Implement live reloading for the `--serve` mode.
