# ghdl
A simple cli tool for working out how many times your github release assets have been downloaded in total.

**Warning:** This tool is early stages and usage will more than likely change.

currently only works on public repositories.

### Installation
```bash
go get github.com/jakehl/ghdl
```

### Usage

```bash
Usage: ghdl [--terse] OWNER REPO

Positional arguments:
  OWNER                  Repository owner
  REPO                   Repository name

Options:
  --terse, -t            Minimal output mode
  --help, -h             display this help and exit
```
Example
```bash
ghdl JakeHL ghdl
```
output
```
Fetching release downloads for https://github.com/JakeHL/ghdl
Assets have been downloaded 1994 times
```

### TODO
- ~~Add a flag for simple output~~
- Add auth for private repositories