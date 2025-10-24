
````markdown
# Domsnif

Domsnif is a CLI tool that helps web developers find gigs by discovering new websites marked as coming soon or under development, extracting associated emails, verifying them, and forwarding them to the developer. Originally built for personal use, it is now open to the public.

---

## Features

- Fetch newly registered or upcoming websites.
- Extract possible emails associated with those domains.
- Validate emails to ensure they are reachable.
- Export validation results.
- Autocompletion support for your shell.

---

## Installation

Make sure you have [Go](https://golang.org/) installed.

Clone the repository:

```bash
git clone https://github.com/mburu72/domsniff.git
cd domsnif
````

Build the executable:

```bash
go build -o domsnif

./domsnif
```

Or install globally:

```bash
go install .
```

---

## Usage

```bash
domsnif [command]
```

Available commands:

* `completion` – Generate autocompletion script for your shell
* `export` – Export validation results
* `fetch` – Fetch newly registered domains
* `validate` – Validate possible emails for each domain
* `help` – Show help for commands

Global flags:

* `--config string` – Specify config file (default is `$HOME/.domsnif.yaml`)
* `-h, --help` – Show help
* `-t, --toggle` – Example toggle flag

Check the version:

```bash
domsnif --version
```

---

## Contributing

Contributions are welcome! Here's how you can help:

1. **Fork the repository**.
2. **Create a new branch** for your feature or bugfix:

```bash
git checkout -b feature/your-feature-name
```

3. **Make your changes** and commit them:

```bash
git commit -m "Add your feature"
```

4. **Push to your fork**:

```bash
git push origin feature/your-feature-name
```

5. **Open a Pull Request** on the main repository.

Please follow the existing code style and include clear commit messages.

---

## License

This project is open-source. See the [LICENSE](LICENSE) file for details.

---

## Contact

For questions, reach out to **[edupablo72@gmail.com](mailto:edupablo72@gmail.com)**.


