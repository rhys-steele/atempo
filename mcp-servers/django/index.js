#!/usr/bin/env node

import { Server } from "@modelcontextprotocol/sdk/server/index.js";
import { StdioServerTransport } from "@modelcontextprotocol/sdk/server/stdio.js";
import {
  CallToolRequestSchema,
  ListToolsRequestSchema,
  McpError,
  ErrorCode,
} from "@modelcontextprotocol/sdk/types.js";
import { exec } from "child_process";
import { promisify } from "util";
import fs from "fs/promises";
import path from "path";

const execAsync = promisify(exec);

class DjangoMCPServer {
  constructor() {
    this.server = new Server(
      {
        name: "steele-django-mcp-server",
        version: "1.0.0",
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.setupToolHandlers();
    this.setupErrorHandling();
  }

  setupErrorHandling() {
    this.server.onerror = (error) => console.error("[MCP Error]", error);
    process.on("SIGINT", async () => {
      await this.server.close();
      process.exit(0);
    });
  }

  setupToolHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => ({
      tools: [
        {
          name: "django_manage",
          description: "Run Django management commands",
          inputSchema: {
            type: "object",
            properties: {
              command: {
                type: "string",
                description: "The management command to run (without 'python manage.py')",
              },
              args: {
                type: "array",
                items: { type: "string" },
                description: "Additional arguments for the command",
                default: [],
              },
            },
            required: ["command"],
          },
        },
        {
          name: "django_startapp",
          description: "Create a new Django app",
          inputSchema: {
            type: "object",
            properties: {
              name: {
                type: "string",
                description: "Name of the Django app",
              },
              directory: {
                type: "string",
                description: "Directory to create the app in",
                default: "",
              },
            },
            required: ["name"],
          },
        },
        {
          name: "django_migrations",
          description: "Handle Django migrations",
          inputSchema: {
            type: "object",
            properties: {
              action: {
                type: "string",
                enum: ["makemigrations", "migrate", "showmigrations", "sqlmigrate"],
                description: "Migration action to perform",
              },
              app: {
                type: "string",
                description: "Specific app name (optional)",
                default: "",
              },
              options: {
                type: "array",
                items: { type: "string" },
                description: "Additional options",
                default: [],
              },
            },
            required: ["action"],
          },
        },
        {
          name: "django_shell",
          description: "Execute Python code in Django shell environment",
          inputSchema: {
            type: "object",
            properties: {
              code: {
                type: "string",
                description: "Python code to execute",
              },
            },
            required: ["code"],
          },
        },
        {
          name: "django_test",
          description: "Run Django tests",
          inputSchema: {
            type: "object",
            properties: {
              path: {
                type: "string",
                description: "Specific test file, class, or method",
                default: "",
              },
              options: {
                type: "array",
                items: { type: "string" },
                description: "Test options (e.g., --verbosity=2, --keepdb)",
                default: [],
              },
            },
          },
        },
        {
          name: "django_collectstatic",
          description: "Collect static files",
          inputSchema: {
            type: "object",
            properties: {
              options: {
                type: "array",
                items: { type: "string" },
                description: "Collectstatic options (e.g., --noinput, --clear)",
                default: ["--noinput"],
              },
            },
          },
        },
        {
          name: "django_runserver",
          description: "Start Django development server",
          inputSchema: {
            type: "object",
            properties: {
              host: {
                type: "string",
                description: "Host to bind to",
                default: "0.0.0.0",
              },
              port: {
                type: "string",
                description: "Port to bind to",
                default: "8000",
              },
            },
          },
        },
        {
          name: "django_check",
          description: "Check Django project for issues",
          inputSchema: {
            type: "object",
            properties: {
              options: {
                type: "array",
                items: { type: "string" },
                description: "Check options (e.g., --deploy, --tag=security)",
                default: [],
              },
            },
          },
        },
        {
          name: "django_dbshell",
          description: "Access Django database shell",
          inputSchema: {
            type: "object",
            properties: {
              database: {
                type: "string",
                description: "Database alias",
                default: "default",
              },
            },
          },
        },
        {
          name: "django_createsuperuser",
          description: "Create Django superuser",
          inputSchema: {
            type: "object",
            properties: {
              username: {
                type: "string",
                description: "Username for superuser",
              },
              email: {
                type: "string",
                description: "Email for superuser",
              },
              noinput: {
                type: "boolean",
                description: "Skip interactive prompts",
                default: false,
              },
            },
          },
        },
        {
          name: "celery_worker",
          description: "Start Celery worker",
          inputSchema: {
            type: "object",
            properties: {
              loglevel: {
                type: "string",
                enum: ["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"],
                description: "Log level",
                default: "INFO",
              },
              concurrency: {
                type: "string",
                description: "Number of worker processes",
                default: "4",
              },
            },
          },
        },
        {
          name: "celery_beat",
          description: "Start Celery beat scheduler",
          inputSchema: {
            type: "object",
            properties: {
              loglevel: {
                type: "string",
                enum: ["DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"],
                description: "Log level",
                default: "INFO",
              },
            },
          },
        },
        {
          name: "pip_install",
          description: "Install Python packages",
          inputSchema: {
            type: "object",
            properties: {
              packages: {
                type: "array",
                items: { type: "string" },
                description: "Package names to install",
              },
              requirements: {
                type: "boolean",
                description: "Install from requirements.txt",
                default: false,
              },
            },
          },
        },
        {
          name: "pip_freeze",
          description: "List installed packages",
          inputSchema: {
            type: "object",
            properties: {
              format: {
                type: "string",
                enum: ["freeze", "list"],
                description: "Output format",
                default: "freeze",
              },
            },
          },
        },
      ],
    }));

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case "django_manage":
            return await this.runManage(args.command, args.args || []);
          case "django_startapp":
            return await this.startApp(args.name, args.directory || "");
          case "django_migrations":
            return await this.handleMigrations(args.action, args.app || "", args.options || []);
          case "django_shell":
            return await this.runShell(args.code);
          case "django_test":
            return await this.runTests(args.path || "", args.options || []);
          case "django_collectstatic":
            return await this.collectStatic(args.options || ["--noinput"]);
          case "django_runserver":
            return await this.runServer(args.host || "0.0.0.0", args.port || "8000");
          case "django_check":
            return await this.checkProject(args.options || []);
          case "django_dbshell":
            return await this.dbShell(args.database || "default");
          case "django_createsuperuser":
            return await this.createSuperuser(args.username, args.email, args.noinput || false);
          case "celery_worker":
            return await this.startCeleryWorker(args.loglevel || "INFO", args.concurrency || "4");
          case "celery_beat":
            return await this.startCeleryBeat(args.loglevel || "INFO");
          case "pip_install":
            return await this.pipInstall(args.packages || [], args.requirements || false);
          case "pip_freeze":
            return await this.pipFreeze(args.format || "freeze");
          default:
            throw new McpError(ErrorCode.MethodNotFound, `Unknown tool: ${name}`);
        }
      } catch (error) {
        throw new McpError(ErrorCode.InternalError, `Tool execution failed: ${error.message}`);
      }
    });
  }

  async runManage(command, args = []) {
    const fullCommand = ["docker-compose", "exec", "-T", "web", "python", "manage.py", command, ...args].join(" ");
    const { stdout, stderr } = await execAsync(fullCommand);
    return {
      content: [
        {
          type: "text",
          text: `Command: ${fullCommand}\n\nOutput:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
        },
      ],
    };
  }

  async startApp(name, directory = "") {
    const args = directory ? [name, directory] : [name];
    return await this.runManage("startapp", args);
  }

  async handleMigrations(action, app = "", options = []) {
    const args = app ? [app, ...options] : options;
    return await this.runManage(action, args);
  }

  async runShell(code) {
    const tempFile = `/tmp/django_shell_${Date.now()}.py`;
    await fs.writeFile(tempFile, code);
    
    try {
      const { stdout, stderr } = await execAsync(
        `docker-compose exec -T web python manage.py shell < ${tempFile}`
      );
      return {
        content: [
          {
            type: "text",
            text: `Django Shell Output:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
          },
        ],
      };
    } finally {
      await fs.unlink(tempFile).catch(() => {});
    }
  }

  async runTests(path = "", options = []) {
    const args = path ? [path, ...options] : options;
    return await this.runManage("test", args);
  }

  async collectStatic(options = ["--noinput"]) {
    return await this.runManage("collectstatic", options);
  }

  async runServer(host, port) {
    return await this.runManage("runserver", [`${host}:${port}`]);
  }

  async checkProject(options = []) {
    return await this.runManage("check", options);
  }

  async dbShell(database) {
    const args = database !== "default" ? [`--database=${database}`] : [];
    return await this.runManage("dbshell", args);
  }

  async createSuperuser(username, email, noinput = false) {
    const args = [];
    if (username) args.push(`--username=${username}`);
    if (email) args.push(`--email=${email}`);
    if (noinput) args.push("--noinput");
    
    return await this.runManage("createsuperuser", args);
  }

  async startCeleryWorker(loglevel, concurrency) {
    const command = [
      "docker-compose", "exec", "-T", "worker", 
      "celery", "-A", "config", "worker", 
      "-l", loglevel.toLowerCase(),
      "--concurrency=" + concurrency
    ].join(" ");
    
    const { stdout, stderr } = await execAsync(command);
    return {
      content: [
        {
          type: "text",
          text: `Celery Worker:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
        },
      ],
    };
  }

  async startCeleryBeat(loglevel) {
    const command = [
      "docker-compose", "exec", "-T", "beat",
      "celery", "-A", "config", "beat",
      "-l", loglevel.toLowerCase()
    ].join(" ");
    
    const { stdout, stderr } = await execAsync(command);
    return {
      content: [
        {
          type: "text",
          text: `Celery Beat:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
        },
      ],
    };
  }

  async pipInstall(packages = [], requirements = false) {
    let command;
    if (requirements) {
      command = ["docker-compose", "exec", "-T", "web", "pip", "install", "-r", "requirements.txt"].join(" ");
    } else {
      command = ["docker-compose", "exec", "-T", "web", "pip", "install", ...packages].join(" ");
    }
    
    const { stdout, stderr } = await execAsync(command);
    return {
      content: [
        {
          type: "text",
          text: `Pip Install:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
        },
      ],
    };
  }

  async pipFreeze(format) {
    const command = format === "list" ? "pip list" : "pip freeze";
    const fullCommand = ["docker-compose", "exec", "-T", "web", ...command.split(" ")].join(" ");
    
    const { stdout, stderr } = await execAsync(fullCommand);
    return {
      content: [
        {
          type: "text",
          text: `Installed Packages:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
        },
      ],
    };
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error("Django MCP server running on stdio");
  }
}

const server = new DjangoMCPServer();
server.run().catch(console.error);