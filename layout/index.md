# Hello, Opabinia ![logo](assets/favicon/favicon-32x32.png)

This is a sample markdown page with image and _formatting_[^1].

Check this [hyperlink](pages/other_page.md)!

---

<br>

## Working with Opabinia

1. Install opabinia engine

       go install github.com/dissipative/opabinia

2. Create directory with your markdown files and assets:
    - generate new project with `opabinia --init <project name>`
    - put .md files to `<project dir>/pages/`
    - put images and other files to `<project dir>/assets/`
    - edit template at `<project dir>/templates/default.tmpl`
3. Start engine with `opabinia --serve`
4. Compile static website with `opabinia --compile`

[^1]: Footnotes work too.