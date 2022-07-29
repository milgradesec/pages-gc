# pages-gc

[![CI](https://github.com/milgradesec/pages-gc/actions/workflows/ci.yml/badge.svg)](https://github.com/milgradesec/pages-gc/actions/workflows/ci.yml)

Small golang tool to delete old deployments from Cloudflare Pages.

## Usage

```shell
export CLOUDFLARE_EMAIL='user@example.com'
export CLOUDFLARE_API_KEY='XXXXXXXXXXXX'

export CLOUDFLARE_ACCOUNT_ID='XXXXXXXXX'
export CLOUDFLARE_PAGES_PROJECT='project-name'

./pages-gc
```
