import os
from datetime import datetime

import requests

email = os.getenv("CLOUDFLARE_EMAIL")
api_key = os.getenv("CLOUDFLARE_API_KEY")

# TODO: get from api?¿
account_id = os.getenv("CLOUDFLARE_ACCOUNT_ID")

# TODO: get all projects from api
project_name = os.getenv("CLOUDFLARE_PAGES_PROJECT")

url = "https://api.cloudflare.com/client/v4/accounts/{0}/pages/projects/{1}/deployments".format(
    account_id, project_name)

headers = {
    "X-Auth-Email": email,
    "X-Auth-Key": api_key
}

resp = requests.get(url, headers=headers)
if resp.status_code != 200:
    # TODO: log
    exit(-1)
deployments = resp.json()

result = deployments["result"]
for d in result:
    timedelta = datetime.now() - \
        datetime.strptime(d["created_on"], "%Y-%m-%dT%H:%M:%S.%fZ")

    if timedelta.days < 30:
        print("🚧 Skipped deployment id={0}".format(d["id"]))
        continue

    resp = requests.delete(url + "/" + d["id"], headers=headers)
    if resp.status_code != 200:
        print("❌ Failed to delete deployment id={0} request={1}".format(
            d["id"], resp.text))
        continue
    print("🧹 Deleted deployment id={0}".format(d["id"]))
