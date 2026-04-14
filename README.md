# portwatch

A lightweight CLI daemon that monitors port usage and alerts on unexpected bindings or conflicts.

---

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git
cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon and watch for unexpected port activity:

```bash
portwatch start --config config.yaml
```

Watch specific ports and get alerted on new bindings:

```bash
portwatch watch --ports 8080,3306,5432 --alert stdout
```

Run a one-time scan of currently bound ports:

```bash
portwatch scan
```

### Example Config (`config.yaml`)

```yaml
watch:
  ports: [80, 443, 8080, 3306]
  alert: stdout
interval: 5s
```

---

## Features

- Detects new or unexpected port bindings in real time
- Alerts via stdout, log file, or webhook
- Minimal resource footprint
- Configurable scan intervals and port allowlists

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)