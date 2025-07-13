# Atempo

**Atempo** is a command-line tool for bootstrapping modern, AI-enabled development environments.

It scaffolds framework-native codebases (e.g. Laravel, Node, Django) using official installers, then layers in a Claude-ready AI context system and optional Docker-based infrastructure. Atempo is designed for speed, clarity, and long-term reuse.

[![Go Version](https://img.shields.io/badge/go-1.22+-brightgreen.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Atempo%20OSL-blue)](./LICENSE)
[![Build](https://img.shields.io/badge/build-passing-brightgreen)]()
[![CLI](https://img.shields.io/badge/cli-atempo-informational)]()

---

## Overview

Atempo creates clean, opinionated project structures with:

- Framework-native source code (via official installers)
- AI-first Claude context system
- Optional Docker setup for local development
- Project metadata via `atempo.json`

---

## Installation

You’ll need Go 1.22+ installed:

```bash
go install github.com/yourname/atempo@latest
```

Alternatively, build from source:

```bash
git clone https://github.com/yourname/atempo
cd atempo
go build -o atempo .
```

(Optional) Move it to your path:

```bash
mv atempo /usr/local/bin/
```

### DNS Requirements

Atempo uses DNSmasq for local DNS management to provide clean domain names for your projects (e.g., `myproject.local` instead of `localhost:3000`).

**Required Setup:**

```bash
# Install DNSmasq (macOS)
brew install dnsmasq

# Start DNSmasq service
sudo brew services start dnsmasq

# Create resolver directory
sudo mkdir -p /etc/resolver
```

**What this enables:**
- Access projects via clean URLs like `myproject.local`
- Automatic service subdomains like `api.myproject.local`
- No port conflicts between multiple projects
- Professional local development environment

**DNS Configuration:**
Atempo automatically configures:
- DNSmasq configuration files in `/opt/homebrew/etc/dnsmasq.d/`
- macOS resolvers in `/etc/resolver/`
- Nginx reverse proxy for domain routing

---

## Usage

### Create a new Laravel project

```bash
# For production projects
mkdir my-app && cd my-app
atempo create laravel:12

# For testing and experimentation
atempo create laravel:12 testing/my-test-app
```

This sets up:

```
my-app/
├── src/                # Laravel 12 source code
├── ai/                 # AI + MCP context system
├── infra/              # Docker environment
├── atempo.json         # Framework metadata
└── README.md
```

### **CRITICAL: Testing Directory Convention**

**ALL testing and experimental projects MUST be created in the `/testing` directory!**

```bash
# ✅ CORRECT: Use testing/ directory for test projects
atempo create laravel testing/my-test-app
atempo create django testing/api-experiment

# ❌ WRONG: Never create test projects in root directory
atempo create laravel my-test-app       # Creates clutter
```

**Why this matters:**
- Keeps the workspace clean and organized
- Easy to ignore all test projects with `.gitignore`
- Prevents accidental commits of test code
- Clear separation between production and experimental work

## What Is Atempo?

**Atempo** helps you scaffold new application projects with:

- Language/framework starter templates (Laravel, Node, etc.)
- MCP-ready Claude context (`ai/context/context.yaml`)
- Dev infrastructure (`infra/docker/`)
- Dev-friendly commands (`atempo docker`, `atempo artisan`, etc.)
- AI integration points (prompts, agents, Claude interface)

Think of it as a smarter `create-react-app`, but for **any stack** — and AI-aware from day one.

---

## Core CLI Commands

### Start a New Project
```bash
atempo create laravel:12
```

Scaffolds a Laravel 12 project with context and infra.
=======
This will:

- Run Laravel's official installer in a Docker container
- Populate your `/src` directory
- Add Claude context (`/ai/context.yaml`)
- Add optional Docker environment (`/infra/docker/`)
- Create `atempo.json` for future commands

---

## Project Structure

After running `atempo start`, you’ll get:

### Docker Commands
```bash
atempo docker up       # Start containers
atempo docker bash     # Open shell into app container
```

---

### Artisan Passthrough
```bash
atempo artisan migrate:fresh
```

---

## Commands

- `atempo context edit`: edit your Claude context file
- `atempo claude "Generate a service for onboarding users"`: injects context and prompts Claude
- `atempo create node:20`, `atempo create react`, etc.
- `atempo generate test`, `atempo prompt --from src/Service.php`

---

## Why Use Atempo?

- Reuse across all your projects
- Keeps your Claude and LLM tooling separated from code
- No more polluting your projects with boilerplate or AI agents
- Prepares every project to use AI as a co-developer, not a side tool

---

## Philosophy

> Tools should get out of your way and *into your headspace*.

**Atempo separates the CLI from your projects.**  
It bootstraps your app, manages Docker, feeds AI agents with context, and gives you smart commands — without embedding anything unnecessary into your codebase.

You install Atempo once. You use it everywhere.

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

Licensed under the **Atempo Open Source License v1.0**  
See [LICENSE](./LICENSE) for details.
