# prefix-filenames-with-date

Walks through a directory and checks all files if they begin with "yyyy-mm-dd ".
If not, renames the files to begin with the current date. For example:

```bash
test/some-invoice.pdf => test/2020-05-31 some-invoice.pdf
```

I run this program from Alfred to tag all my tax related files with receive dates.

## Build

```bash
go build prefix-filenames-with-date.go
```

## Use

Usage:

```bash
prefix-filenames-with-date <directory>
```
