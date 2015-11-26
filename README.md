# 1pwd

:lock: 1pwd :unlock: A command-line tool for searching 1Password vaults.

## Install

Install `fzf` ([instructions](https://github.com/junegunn/fzf#installation))

```sh
GO15VENDOREXPERIMENT=1 go get github.com/fd/1pwd/cmd/...
```

## Usage

```sh
# get a single entry
1pwd [--vault=PATH] get ID [FIELD] [--json]

# search for an entry
1pwd [--vault=PATH] search [FIELD] [--query=QUERY] [--type=TYPE] [--json]
```
