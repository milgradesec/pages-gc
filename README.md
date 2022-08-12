# pages-gc

[![CI](https://github.com/milgradesec/pages-gc/actions/workflows/ci.yml/badge.svg)](https://github.com/milgradesec/pages-gc/actions/workflows/ci.yml)

Simple CLI tool to delete deployments from Cloudflare Pages.

## Usage

```shell
pages-gc [OPTIONS]

Options:
  -account string
        Cloudflare account ID
  -apikey string
        Cloudflare API key
  -email string
        Cloudflare account email
  -project string
        Pages project name (default "all")
```

Example:

```shell
export CLOUDFLARE_EMAIL='user@example.com'
export CLOUDFLARE_API_KEY='XXXXXXXXXXXX'
export CLOUDFLARE_ACCOUNT_ID='XXXXXXXXX'

./pages-gc \
      -email $CLOUDFLARE_EMAIL \
      -apikey $CLOUDFLARE_API_KEY \
      -account $CLOUDFLARE_ACCOUNT_ID \
      -project 'project-name'
```
