# 🛠️ Atempo

**Atempo** is a developer-first, AI-enhanced CLI for bootstrapping **framework-agnostic** projects with built-in support for Claude, MCP (Model–Context–Protocol), Docker, and best-practice architecture.  

Start fast. Stay scalable. Think with context.

---

## 🚀 Quick Start

### 📦 Install

If you have Go installed:

```bash
go install github.com/yourname/atempo@latest
```

Or clone and build manually:

```bash
git clone https://github.com/yourname/atempo
cd atempo
go build -o atempo .
```

Then move it into your path (optional):

```bash
mv atempo /usr/local/bin/
```

---

### 🧱 Create a New Project

```bash
mkdir my-app && cd my-app
atempo start laravel:12
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

---

## 💡 What Is Atempo?

**Atempo** helps you scaffold new application projects with:

- 🚀 Language/framework starter templates (Laravel, Node, etc.)
- 🧠 MCP-ready Claude context (`ai/context/context.yaml`)
- 🐳 Dev infrastructure (`infra/docker/`)
- 🧪 Dev-friendly commands (`atempo docker`, `atempo artisan`, etc.)
- ✨ AI integration points (prompts, agents, Claude interface)

Think of it as a smarter `create-react-app`, but for **any stack** — and AI-aware from day one.

---

## 🧰 Core CLI Commands

### 🆕 Start a New Project
```bash
atempo start laravel:12
```

Scaffolds a Laravel 12 project with context and infra.

---

### 🐳 Docker Commands
```bash
atempo docker up       # Start containers
atempo docker bash     # Open shell into app container
```

---

### 🧪 Artisan Passthrough
```bash
atempo artisan migrate:fresh
```

Runs Laravel `artisan` commands inside Docker transparently.

---

## 🔮 Coming Soon

- `atempo context edit`: edit your Claude context file
- `atempo claude "Generate a service for onboarding users"`: injects context and prompts Claude
- `atempo start node:20`, `atempo start react`, etc.
- `atempo generate test`, `atempo prompt --from src/Service.php`

---

## 🌍 Why Use Atempo?

- 🔁 Reuse across all your projects
- 🧠 Keeps your Claude and LLM tooling separated from code
- 🚫 No more polluting your projects with boilerplate or AI agents
- 💬 Prepares every project to use AI as a co-developer, not a side tool

---

## 🧱 Philosophy

> Tools should get out of your way and *into your headspace*.

**Atempo separates the CLI from your projects.**  
It bootstraps your app, manages Docker, feeds AI agents with context, and gives you smart commands — without embedding anything unnecessary into your codebase.

You install Atempo once. You use it everywhere.

---

## 🧩 Framework Roadmap

| Framework   | Status   |
|-------------|----------|
| Laravel     | ✅ Ready |
| Node.js     | 🔜       |
| Django      | 🔜       |
| React       | 🔜       |
| Astro       | 🔜       |

---

## 📄 License

MIT – © 2025 Rhys May
