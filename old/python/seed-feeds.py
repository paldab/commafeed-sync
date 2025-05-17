from typing import Dict, List, TypedDict
from requests import request
from requests.models import Response
import os, re, yaml


class InputFeed(TypedDict):
    name: str
    url: str


class InputCategory(TypedDict):
    name: str
    urls: List[InputFeed]


class Node(TypedDict):
    id: int
    parentId: str | None
    parentName: str
    name: str
    children: List
    feeds: List
    expanded: bool
    position: int


def call_api(endpoint: str, method: str, json=None, data=None) -> Response:
    return request(method, f"{url}{endpoint}", auth=auth, json=json, data=data)


def get_categories_from_api() -> Dict:
    endpoint = r"/rest/category/get"

    res = call_api(endpoint, "GET")
    return res.json()


def create_category_via_api(name: str, parent: str) -> Response:
    endpoint = "/rest/category/add"
    payload = {"name": name, "parentId": parent}

    res = call_api(endpoint, "POST", json=payload)

    if res.status_code != 200:
        raise Exception(f"Failed to create category {name}: {res.text}")
    else:
        print(f"Created category {name}")

    return res


def subscribe_to_feed_via_api(name: str, feed_url: str, category_id: str) -> Response:
    endpoint = "/rest/feed/subscribe"
    payload = {"url": feed_url, "title": name, "categoryId": category_id}

    res = call_api(endpoint, "POST", json=payload)

    if res.status_code != 200:
        raise Exception(f"Failed to subscribe to feed {feed_url}: {res.text}")
    else:
        print(f"Subscribing to feed {name} with url: {feed_url}")

    return res


def unsubscribe_feed_via_api(id: int) -> Response:
    endpoint = "/rest/feed/unsubscribe"
    payload = {"id": id}

    return call_api(endpoint, "POST", json=payload)


def delete_category_via_api(id: int) -> Response:
    endpoint = r"/rest/category/delete"
    payload = {"id": id}

    return call_api(endpoint, "POST", json=payload)


def build_declared_config_tree(input_categories: InputCategory, existing_data_tree):

    def find_child_by_name(children: List, name: str) -> Dict | None:
        for child in children:
            if child["name"].strip().lower() == name.strip().lower():
                return child
        return None

    if "name" not in input_categories or not input_categories["name"]:
        raise Exception(
            "Could not parse category, no property of 'name' found in feeds"
        )

    category_name = input_categories["name"]
    feed_urls = input_categories["urls"] or []

    current_children = existing_data_tree.get("children", [])
    current_parent_id = "all"
    is_nested_category = "/" in category_name

    if is_nested_category:
        parts = category_name.split("/")

        for part in parts:
            exists = find_child_by_name(current_children, part)

            if exists:
                # Existing category
                current_parent_id = exists["id"]
                current_children = exists.get("children", [])
            else:
                # Create Category
                res = None
                try:
                    res = create_category_via_api(part, current_parent_id)
                except Exception as e:
                    print(f"Exception while creating category '{part}': {e}")
                    return

                # New id
                current_parent_id = res.json()

        create_feeds(feed_urls, current_parent_id)
    else:
        category_exists = False
        for category in existing_data_tree.get("children", []):
            category_exists = category["name"] == category_name

        if not category_exists:
            try:
                res = create_category_via_api(category_name, current_parent_id)
            except Exception as e:
                print(f"Exception while creating category '{category_name}': {e}")
                return

        create_feeds(feed_urls, current_parent_id)


def create_feeds(feeds: List[InputFeed], parent: str):
    for feed in feeds:
        url = feed["url"]
        name = feed["name"]

        try:
            subscribe_to_feed_via_api(name, url, parent)
        except Exception as e:
            print(f"Exception while subscribe to feed '{name}': {e}")


def clean_up_data(input_categories: List[InputCategory]):

    def normalize_category_name(name: str):
        return re.sub(r"\s*/\s*", "/", name.strip())

    # Returns if full_path is no longer declared in config and does not have children
    def should_cleanup_category(full_category_path: str, node: Node) -> bool:
        full_category_path = normalize_category_name(full_category_path)
        return (
            full_category_path not in valid_categories  # Category isn't in config
            and not node.get("children")
        )

    def cleanup_category(full_category_path: str, node: Node, name: str):
        deleted.append(full_category_path)

        # Remove feeds if exist
        feeds = node.get("feeds")
        if feeds:
            for feed in feeds:
                feed_id = feed["id"]
                feed_name = feed["name"]

                print(f"Deleting feed: {feed_name}")
                unsubscribe_feed_via_api(feed_id)

        # Remove category
        print(f"Deleting Category: {name}")
        delete_category_via_api(node["id"])

    def cleanup_feeds(node: Node, full_category_path: str):
        existing_feeds = node.get("feeds", [])
        declared_categories = list(
            filter(
                lambda x: normalize_category_name(x["name"]) == full_category_path,
                input_categories,
            )
        )

        result = len(declared_categories) == 1
        if result:
            declared_target_feeds = declared_categories[0]["urls"]
            declared_target_feed_urls = list(
                map(lambda feed: feed["url"], declared_target_feeds)
            )
            for feed in existing_feeds:
                item_feed_url: str = feed.get("feedUrl", "").removesuffix(".atom")

                # Existing item not in declared config
                if item_feed_url and item_feed_url not in declared_target_feed_urls:
                    item_id = feed["id"]
                    unsubscribe_feed_via_api(item_id)
                    print(f"Removed feed with id: {item_id} and url: {item_feed_url}")

    def recurse(node: Node, parent_path=""):
        name = node["name"].strip()
        full_category_path = f"{parent_path}/{name}" if parent_path else name

        # Recurse into children first
        children = node.get("children", [])
        new_children = []
        for child in children:
            if not recurse(child, full_category_path):  # If child is kept, retain it
                new_children.append(child)
        node["children"] = new_children

        # Check if we should delete the category
        if should_cleanup_category(full_category_path, node):
            cleanup_category(full_category_path, node, name)
            return True  # Signal to parent to remove this node

        # # Remove not declared feeds
        cleanup_feeds(node, full_category_path)

        return False

    # Renew existing data tree after data has been created
    existing_data = get_categories_from_api()
    valid_categories = {
        normalize_category_name(category["name"]) for category in input_categories
    }

    deleted = []

    # Start recursion from top-level categories
    for top in existing_data.get("children", []):
        recurse(top)

    return deleted


def validated_commafeed_url() -> str:
    url = os.getenv("COMMAFEED_URL")

    if not url:
        print("variable 'COMMAFEED_URL' is empty. missing commafeed url")
        exit(1)

    if not "http" in url or not "https" in url:
        print(f"Missing http protocol in {url}")
        exit(1)

    return url


# Start script

# Load config file from mounted volume
CONFIG_PATH = os.getenv("CONFIG_PATH", "/commafeed/feeds.yaml")
print(f"Loading config from {CONFIG_PATH}")

try:
    with open(CONFIG_PATH) as f:
        config = yaml.safe_load(f)
except Exception as e:
    print(f"Failed to read config: {e}")
    exit(1)

url = validated_commafeed_url()
print(f"Found commafeed at {url}")

auth = (
    os.getenv("COMMAFEED_USER", "admin"),
    os.getenv("COMMAFEED_PASS", "admin"),
)

if "feeds" not in config or not config["feeds"]:
    print("Could not find 'feeds' property or it is empty")
    exit(1)

input_data = config["feeds"]
existing_data_tree = get_categories_from_api()


print("Removing undeclared feeds")
clean_up_data(input_data)

print("Adding declared feeds")
for category in input_data:
    build_declared_config_tree(category, existing_data_tree)

print("Finished setting up feeds!")
