# Vulnerable Target


<img width="530" height="137" alt="image" src="https://github.com/user-attachments/assets/5601f3a8-7037-4fe5-8115-eb4182d8ab2e" />



Vulnerable Target (VT) is a specialized tool designed for security professionals, researchers, and educators that creates intentionally vulnerable environments across multiple platforms.

> [!CAUTION]
> **SECURITY WARNING: DO NOT RUN ON UNTRUSTED NETWORKS**
> This tool creates intentionally vulnerable environments. Running this on a public server or an insecure network can expose you to severe security risks. Use only in an isolated, local environment (sandbox/VM).

## Features
- CLI for managing vulnerable environments
- Docker Compose provider for container orchestration
- Community-curated templates from [vt-templates](https://github.com/HappyHackingSpace/vt-templates)
- Template filtering by tags
- Deployment state tracking

## Prerequisites
- Go 1.24+
- Docker & Docker Compose

## Installation

1. Clone the repository
```bash
git clone https://github.com/HappyHackingSpace/vulnerable-target.git
cd vulnerable-target
```

2. Install dependencies
```bash
go mod download
```

3. Build the binary
```bash
go build -o vt cmd/vt/main.go
```

4. (Optional) Move to your PATH
```bash
mv vt /usr/local/bin/
```

## Usage

```bash
# List available templates
vt template --list

# Filter templates by tag
vt template --list --filter sql

# Update templates from remote repository
vt template --update

# Start a vulnerable environment
vt start --id <template-id> --provider docker-compose

# List running environments
vt ps

# Stop an environment
vt stop --id <template-id> --provider docker-compose

# Set verbosity level
vt -v debug <command>
```

## Templates

Templates are automatically cloned to `~/vt-templates` on first run. To contribute new vulnerable target templates, visit the [vt-templates repository](https://github.com/HappyHackingSpace/vt-templates).

## Documentation
Check the full documentation here: [Vulnerable Target Wiki](https://github.com/HappyHackingSpace/vulnerable-target/wiki)

## How to Contribute
Hack! don't forget to follow [CONTRIBUTING.md](./CONTRIBUTING.md)

## Disclaimer
### Use with caution and additional security measures.

---
Happy Hacking!
