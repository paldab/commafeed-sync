# Declaritive commafeed feeds sync

## Yaml structure

```YAML
  commafeedSetup:
    - category: Category1
      feeds:
        - name: Commafeed Releases
          url: "https://github.com/Athou/commafeed/releases"
        - name: Reddit golang
          url: "https://www.reddit.com/r/golang/.rss"

    - category: Category2 / News
      feeds:
        - name: BCC
          url: "https://feeds.bbci.co.uk/news/world/rss.xml"
        - name: Nasa
          url: "https://www.nasa.gov/rss/dyn/breaking_news.rss"
          disabled: true # Optional, defaulted false
        - name: Hacker News
          url: "https://hnrss.org/frontpage"
```

