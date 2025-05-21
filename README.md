# CommafeedSync

**CommafeedSync** is a lightweight tool written in Go that syncs a declarative YAML configuration of feed categories and subscriptions into a running [CommaFeed](https://github.com/Athou/commafeed) instance.

This tool is useful for teams and individuals who want to manage CommaFeed categories and feeds using a GitOps-friendly YAML file — ideal for automation in CI/CD pipelines or homelabs.

---

## 📦 Features

- Automatically creates nested categories based on `/` in the `category` field
- Adds feeds under the correct category structure
- Idempotent: safely re-applies the same config
- Declarative: driven by a single YAML file
- Minimal configuration via environment variables

---

## 📄 Example

Given this YAML config:

```yaml
commafeedSetup:
  - category: Tools
    feeds:
      - name: K9s Releases
        url: "https://github.com/derailed/k9s/releases"

  - category: Dev / Kubernetes
    feeds:
      - name: K3s Releases
        url: "https://github.com/k3s-io/k3s/releases"
      - name: Kubernetes Releases
        url: "https://github.com/kubernetes/kubernetes/releases"

