# Steele

[![Go Version](https://img.shields.io/badge/go-1.22+-brightgreen.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Steele%20OSL-blue)](./LICENSE)
[![Build](https://img.shields.io/badge/build-passing-brightgreen)]()
[![CLI](https://img.shields.io/badge/cli-steele-informational)]()

**Steele** is a command-line tool for bootstrapping modern, AI-enabled development environments.

It scaffolds framework-native codebases (e.g. Laravel, Node, Django) using official installers, then layers in a Claude-ready AI context system and optional Docker-based infrastructure. Steele is designed for speed, clarity, and long-term reuse.

---

## Overview

Steele creates clean, opinionated project structures with:

- Framework-native source code (via official installers)
- AI-first Claude context system
- Optional Docker setup for local development
- Project metadata via `steele.json`

---

## Installation

You’ll need Go 1.22+ installed:

```bash
go install github.com/yourname/steele@latest
```

Alternatively, build from source:

```bash
git clone https://github.com/yourname/steele
cd steele
go build -o steele .
```

(Optional) Move it to your path:

```bash
mv steele /usr/local/bin/
```

---

## Usage

### Create a new Laravel project

```bash
mkdir my-app && cd my-app
steele start laravel:12
```

This will:

- Run Laravel's official installer in a Docker container
- Populate your `/src` directory
- Add Claude context (`/ai/context.yaml`)
- Add optional Docker environment (`/infra/docker/`)
- Create `steele.json` for future commands

---

## Project Structure

After running `steele start`, you’ll get:

```
my-app/
├── src/              # Laravel source code
├── ai/               # Claude context
├── infra/            # Docker infrastructure
│   └── docker/
├── steele.json       # Metadata used by Steele
└── README.md
```

---

## Commands

- `steele start <framework>` – Scaffold a new project
- `steele docker up` – Start Docker services (coming soon)
- `steele context edit` – Modify Claude context (coming soon)
- `steele generate <type>` – Generate services, tests, etc. using Claude (future)

---

## Philosophy

Steele helps developers start fast without sacrificing structure. By combining official frameworks with an AI-first toolchain, Steele keeps code clean, reproducible, and easy to scale — without bloated templates or vendor lock-in.

---

## Roadmap

- Laravel: ✅
- Node.js: ⏳
- Django: ⏳
- React: ⏳
- Custom templates: 🔜
- Service/test generation via Claude: 🔜

---

## License

Licensed under the **Steele Open Source License v1.0**  
See [LICENSE](./LICENSE) for details.
