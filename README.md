# 🛠️ Steele

**Steele** is a developer-first, AI-enhanced CLI for bootstrapping **framework-agnostic** projects with built-in support for Claude, MCP (Model–Context–Protocol), Docker, and best-practice architecture.  

Start fast. Stay scalable. Think with context.

---

## 🚀 Quick Start

### 📦 Install

If you have Go installed:

```bash
go install github.com/yourname/steele@latest
```

Or clone and build manually:

```bash
git clone https://github.com/yourname/steele
cd steele
go build -o steele .
```

Then move it into your path (optional):

```bash
mv steele /usr/local/bin/
```

---

### 🧱 Create a New Project

```bash
mkdir my-app && cd my-app
steele start laravel:12
```

This sets up:

```
my-app/
├── src/                # Laravel 12 source code
├── ai/                 # AI + MCP context system
├── infra/              # Docker environment
├── steele.json         # Framework metadata
└── README.md
```

---

## 💡 What Is Steele?

**Steele** helps you scaffold new application projects with:

- 🚀 Language/framework starter templates (Laravel, Node, etc.)
- 🧠 MCP-ready Claude context (`ai/context/context.yaml`)
- 🐳 Dev infrastructure (`infra/docker/`)
- 🧪 Dev-friendly commands (`steele docker`, `steele artisan`, etc.)
- ✨ AI integration points (prompts, agents, Claude interface)

Think of it as a smarter `create-react-app`, but for **any stack** — and AI-aware from day one.

---

## 🧰 Core CLI Commands

### 🆕 Start a New Project
```bash
steele start laravel:12
```

Scaffolds a Laravel 12 project with context and infra.

---

### 🐳 Docker Commands
```bash
steele docker up       # Start containers
steele docker bash     # Open shell into app container
```

---

### 🧪 Artisan Passthrough
```bash
steele artisan migrate:fresh
```

Runs Laravel `artisan` commands inside Docker transparently.

---

## 🔮 Coming Soon

- `steele context edit`: edit your Claude context file
- `steele claude "Generate a service for onboarding users"`: injects context and prompts Claude
- `steele start node:20`, `steele start react`, etc.
- `steele generate test`, `steele prompt --from src/Service.php`

---

## 🌍 Why Use Steele?

- 🔁 Reuse across all your projects
- 🧠 Keeps your Claude and LLM tooling separated from code
- 🚫 No more polluting your projects with boilerplate or AI agents
- 💬 Prepares every project to use AI as a co-developer, not a side tool

---

## 🧱 Philosophy

> Tools should get out of your way and *into your headspace*.

**Steele separates the CLI from your projects.**  
It bootstraps your app, manages Docker, feeds AI agents with context, and gives you smart commands — without embedding anything unnecessary into your codebase.

You install Steele once. You use it everywhere.

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
