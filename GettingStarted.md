````markdown
## Getting Started

After installing Domsnif, you can start fetching domains and validating emails with the following commands.

---

### 1. Fetch Newly Registered Domains

The `fetch` command retrieves newly registered domains. You can optionally filter them based on status.

**Basic usage:**

```bash
domsnif fetch
````

**With filtering enabled:**

```bash
domsnif fetch --filter
```

**Available flags:**

* `--filter, -f` – Enable filtering based on TLD, online status, and under-development status.
* `--check-online, -o` – Filter domains based on whether they are online.
* `--check-dev, -d` – Filter domains that are under development.

**Example:**

```bash
domsnif fetch --filter --check-online --check-dev
```

This will fetch newly registered domains, check which ones are online, filter under-development domains, and save the results to `filtered_domains.csv` and `under_dev.csv`.

---

### 2. Validate Emails for Each Domain

The `validate` command checks possible emails associated with your fetched domains.

**Basic usage:**

```bash
domsnif validate
```

**With custom input/output and concurrency:**

```bash
domsnif validate --input domains.json --output verified_emails.csv --w 10
```

**Available flags:**

* `--input, -i` – Input JSON file containing domains (default: `domains.json`).
* `--output, -o` – Output CSV file for verified emails (default: `data/verified_emails.csv`).
* `--w` – Number of concurrent workers for verification (default: 5).

---

### Quick Workflow Example

```bash
# Fetch new domains and filter under-development ones
domsnif fetch --filter --check-dev

# Validate emails from the filtered domains
domsnif validate --input filtered_domains.csv --output verified_emails.csv
```

This workflow allows you to discover new projects and verify contact emails efficiently.

