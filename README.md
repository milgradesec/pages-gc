# pages-gc

[![CI](https://github.com/milgradesec/pages-gc/actions/workflows/ci.yml/badge.svg)](https://github.com/milgradesec/pages-gc/actions/workflows/ci.yml)

Small golang tool to delete old deployments from Cloudflare Pages.

## Usage

```shell
pages-gc \
    -email <CLOUDFLARE_ACCOUNT_EMAIL> \
    -key <CLOUDFLARE_API_KEY> \
    -account <CLOUDFLARE_ACCOUNT_ID> \
    -project <PAGES_PROJECT_NAME>
```
