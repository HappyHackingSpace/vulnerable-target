<div align="center">

```text
               __                     __    __        __                       __
 _  __ ____ __/ /____  ___ _______ _/ /   / /_____ _/ /____ ____ _____ ____  / /_
| |/ // / // / // _ \/ -_) __/ _ `/ _ \ / __/ _ `/ __/ _ `/ -_) __/ _ `/ -_)/ __/
|___/\_,_/_//_//_//_/\__/_/  \_,_/_.__/ \__/\_,_/_/  \_, /\__/\__/\_, /\__/ \__/
                                                    /___/        /___/
```

**Create intentionally vulnerable environments for security testing, education, and research**

[![Go Version](https://img.shields.io/github/go-mod/go-version/HappyHackingSpace/vt?style=flat-square)](https://go.dev/)
[![License](https://img.shields.io/github/license/HappyHackingSpace/vt?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/HappyHackingSpace/vt?style=flat-square)](https://github.com/HappyHackingSpace/vt/releases)
[![Discord](https://img.shields.io/badge/Discord-Join-7289DA?style=flat-square&logo=discord&logoColor=white)](https://discord.happyhacking.space)

</div>

> [!Important]
> **This project is in active development.** Expect breaking changes with releases. Review the [release changelog](https://github.com/HappyHackingSpace/vt/releases) before updating. **vt** creates intentionally vulnerable environments - always run in isolated networks (VMs/sandboxes) and never expose to the internet.

---

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Templates](#templates)
- [What can you do with vt?](#what-can-you-do-with-vt)
- [Documentation](#documentation)
- [Community](#community)
- [License](#license)

---

## Features

| | Feature | Description |
|:--:|---------|-------------|
| 🐳 | **Docker Compose** | Container orchestration for vulnerable environments |
| 📦 | **Templates** | Community-curated vulnerable targets from [vt-templates](https://github.com/HappyHackingSpace/vt-templates) |
| 🏷️ | **Tag Filtering** | Find templates by vulnerability type (sqli, xss, ssrf, etc.) |
| 📊 | **State Tracking** | Track and manage running deployments |
| 🔄 | **Auto-Update** | Sync templates from remote repository |

---

## Installation

### Prerequisites

- Go 1.24+
- Docker & Docker Compose

### Install with Go

```bash
go install github.com/happyhackingspace/vulnerable-target/cmd/vt@latest
```

### Build from Source

```bash
git clone https://github.com/HappyHackingSpace/vt.git
cd vt
go build -o vt cmd/vt/main.go
mv vt /usr/local/bin/  # Optional: add to PATH
```

---

## Quick Start

```bash
# 1. Browse available templates
vt template --list

# 2. Start a vulnerable environment
vt start --id vt-dvwa

# 3. Access the target at http://localhost:80
```

---

## Usage

<details>
<summary><b>Command Reference</b></summary>

| Command | Description |
|---------|-------------|
| `vt template --list` | List all available templates |
| `vt template --list --filter <tag>` | Filter templates by tag |
| `vt template --update` | Update templates from remote repository |
| `vt start --id <template-id>` | Start a vulnerable environment |
| `vt start --tags <tag1,tag2>` | Start all templates matching tags |
| `vt ps` | List running environments |
| `vt stop --id <template-id>` | Stop an environment |
| `vt stop --tags <tag1,tag2>` | Stop all templates matching tags |
| `vt -v debug <command>` | Run with debug verbosity |

</details>

### Examples

```bash
# List templates with SQL injection vulnerabilities
vt template --list --filter sqli

# Start DVWA (Damn Vulnerable Web App)
vt start --id vt-dvwa

# Start all XSS-related labs
vt start --tags xss

# Check running environments
vt ps

# Stop a specific environment
vt stop --id vt-dvwa
```

---

## Templates

Templates are automatically cloned to `~/vt-templates` on first run.

| Template | Type | Description |
|----------|:----:|-------------|
| `vt-dvwa` | Lab | Damn Vulnerable Web Application |
| `vt-juice-shop` | Lab | OWASP Juice Shop |
| `vt-webgoat` | Lab | OWASP WebGoat |
| `vt-bwapp` | Lab | Buggy Web Application |
| `vt-mutillidae-ii` | Lab | OWASP Mutillidae II |

> **Want more?** Check out the [vt-templates repository](https://github.com/HappyHackingSpace/vt-templates) for all available templates and contribution guidelines.

---

## What can you do with vt?

| Use Case | Template |
|----------|----------|
| Practice SQL Injection | [vt-dvwa](https://github.com/HappyHackingSpace/vt-templates/tree/main/labs/vt-dvwa) |
| Learn XSS Exploitation | [vt-dvwa](https://github.com/HappyHackingSpace/vt-templates/tree/main/labs/vt-dvwa) |
| Test OWASP Top 10 | [vt-juice-shop](https://github.com/HappyHackingSpace/vt-templates/tree/main/labs/vt-juice-shop) |
| Exploit Real CVEs | [vt-2025-29927](https://github.com/HappyHackingSpace/vt-templates/tree/main/cves/vt-2025-29927) |
| API Security Testing | [vt-webgoat](https://github.com/HappyHackingSpace/vt-templates/tree/main/labs/vt-webgoat) |
| Train Security Teams | [vt-mutillidae-ii](https://github.com/HappyHackingSpace/vt-templates/tree/main/labs/vt-mutillidae-ii) |

---

## Documentation

| | Resource | Description |
|:--:|----------|-------------|
| 📖 | [Wiki](https://github.com/HappyHackingSpace/vt/wiki) | Full documentation and guides |
| 📦 | [Templates](https://github.com/HappyHackingSpace/vt-templates) | Browse all available templates |
| 🤝 | [Contributing](./CONTRIBUTING.md) | Contribution guidelines |
| 🐛 | [Issues](https://github.com/HappyHackingSpace/vt/issues) | Report bugs or request features |

---

## Community

- 💬 **Discord**: [Join our community](https://discord.happyhacking.space)
- 🐛 **Issues**: [Report bugs](https://github.com/HappyHackingSpace/vt/issues)
- 🤝 **Contributing**: Check out [CONTRIBUTING.md](./CONTRIBUTING.md)

---

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Happy Hacking!** 🎯

</div>
