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

class LaravelMCPServer {
  constructor() {
    this.server = new Server(
      {
        name: "steele-laravel-mcp-server",
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
          name: "laravel_artisan",
          description: "Run Laravel Artisan commands",
          inputSchema: {
            type: "object",
            properties: {
              command: {
                type: "string",
                description: "The artisan command to run (without 'php artisan')",
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
          name: "laravel_make",
          description: "Generate Laravel components (controllers, models, migrations, etc.)",
          inputSchema: {
            type: "object",
            properties: {
              type: {
                type: "string",
                enum: ["controller", "model", "migration", "middleware", "request", "resource", "seeder", "factory", "command", "job", "listener", "mail", "notification", "policy", "provider", "rule"],
                description: "Type of component to generate",
              },
              name: {
                type: "string",
                description: "Name of the component",
              },
              options: {
                type: "array",
                items: { type: "string" },
                description: "Additional options (e.g., --resource, --api, --migration)",
                default: [],
              },
            },
            required: ["type", "name"],
          },
        },
        {
          name: "laravel_db",
          description: "Database operations (migrations, seeding, etc.)",
          inputSchema: {
            type: "object",
            properties: {
              action: {
                type: "string",
                enum: ["migrate", "migrate:fresh", "migrate:reset", "migrate:rollback", "db:seed", "migrate:status"],
                description: "Database action to perform",
              },
              options: {
                type: "array",
                items: { type: "string" },
                description: "Additional options (e.g., --force, --step=1)",
                default: [],
              },
            },
            required: ["action"],
          },
        },
        {
          name: "laravel_routes",
          description: "Display Laravel routes",
          inputSchema: {
            type: "object",
            properties: {
              filter: {
                type: "string",
                description: "Filter routes by name, URI, or method",
                default: "",
              },
            },
          },
        },
        {
          name: "laravel_config",
          description: "Get or set Laravel configuration values",
          inputSchema: {
            type: "object",
            properties: {
              action: {
                type: "string",
                enum: ["get", "cache", "clear"],
                description: "Configuration action",
              },
              key: {
                type: "string",
                description: "Configuration key (for get action)",
              },
            },
            required: ["action"],
          },
        },
        {
          name: "laravel_test",
          description: "Run Laravel tests",
          inputSchema: {
            type: "object",
            properties: {
              path: {
                type: "string",
                description: "Specific test file or directory",
                default: "",
              },
              options: {
                type: "array",
                items: { type: "string" },
                description: "Test options (e.g., --filter, --group)",
                default: [],
              },
            },
          },
        },
        {
          name: "laravel_serve",
          description: "Start Laravel development server",
          inputSchema: {
            type: "object",
            properties: {
              host: {
                type: "string",
                description: "Host to bind to",
                default: "127.0.0.1",
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
          name: "laravel_tinker",
          description: "Execute PHP code in Laravel Tinker environment",
          inputSchema: {
            type: "object",
            properties: {
              code: {
                type: "string",
                description: "PHP code to execute",
              },
            },
            required: ["code"],
          },
        },
        {
          name: "composer_install",
          description: "Install Composer dependencies",
          inputSchema: {
            type: "object",
            properties: {
              options: {
                type: "array",
                items: { type: "string" },
                description: "Composer install options",
                default: [],
              },
            },
          },
        },
        {
          name: "composer_require",
          description: "Add new Composer dependencies",
          inputSchema: {
            type: "object",
            properties: {
              packages: {
                type: "array",
                items: { type: "string" },
                description: "Package names to install",
              },
              dev: {
                type: "boolean",
                description: "Install as dev dependency",
                default: false,
              },
            },
            required: ["packages"],
          },
        },
      ],
    }));

    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        switch (name) {
          case "laravel_artisan":
            return await this.runArtisan(args.command, args.args || []);
          case "laravel_make":
            return await this.makeComponent(args.type, args.name, args.options || []);
          case "laravel_db":
            return await this.dbOperation(args.action, args.options || []);
          case "laravel_routes":
            return await this.showRoutes(args.filter || "");
          case "laravel_config":
            return await this.configOperation(args.action, args.key);
          case "laravel_test":
            return await this.runTests(args.path || "", args.options || []);
          case "laravel_serve":
            return await this.startServer(args.host || "127.0.0.1", args.port || "8000");
          case "laravel_tinker":
            return await this.runTinker(args.code);
          case "composer_install":
            return await this.composerInstall(args.options || []);
          case "composer_require":
            return await this.composerRequire(args.packages, args.dev || false);
          default:
            throw new McpError(ErrorCode.MethodNotFound, `Unknown tool: ${name}`);
        }
      } catch (error) {
        throw new McpError(ErrorCode.InternalError, `Tool execution failed: ${error.message}`);
      }
    });
  }

  async runArtisan(command, args = []) {
    const fullCommand = ["docker-compose", "exec", "-T", "app", "php", "artisan", command, ...args].join(" ");
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

  async makeComponent(type, name, options = []) {
    const command = `make:${type}`;
    return await this.runArtisan(command, [name, ...options]);
  }

  async dbOperation(action, options = []) {
    return await this.runArtisan(action, options);
  }

  async showRoutes(filter = "") {
    const args = filter ? ["--name=" + filter] : [];
    return await this.runArtisan("route:list", args);
  }

  async configOperation(action, key) {
    switch (action) {
      case "get":
        if (!key) throw new Error("Key required for get action");
        return await this.runArtisan("tinker", [`--execute=echo config('${key}');`]);
      case "cache":
        return await this.runArtisan("config:cache");
      case "clear":
        return await this.runArtisan("config:clear");
      default:
        throw new Error(`Unknown config action: ${action}`);
    }
  }

  async runTests(path = "", options = []) {
    const args = path ? [path, ...options] : options;
    return await this.runArtisan("test", args);
  }

  async startServer(host, port) {
    return await this.runArtisan("serve", [`--host=${host}`, `--port=${port}`]);
  }

  async runTinker(code) {
    const tempFile = `/tmp/tinker_${Date.now()}.php`;
    await fs.writeFile(tempFile, `<?php\n${code}`);
    
    try {
      const { stdout, stderr } = await execAsync(
        `docker-compose exec -T app php artisan tinker < ${tempFile}`
      );
      return {
        content: [
          {
            type: "text",
            text: `Tinker Output:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
          },
        ],
      };
    } finally {
      await fs.unlink(tempFile).catch(() => {});
    }
  }

  async composerInstall(options = []) {
    const command = ["docker-compose", "exec", "-T", "app", "composer", "install", ...options].join(" ");
    const { stdout, stderr } = await execAsync(command);
    return {
      content: [
        {
          type: "text",
          text: `Composer Install:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
        },
      ],
    };
  }

  async composerRequire(packages, dev = false) {
    const command = [
      "docker-compose", "exec", "-T", "app", "composer", "require",
      ...(dev ? ["--dev"] : []),
      ...packages
    ].join(" ");
    
    const { stdout, stderr } = await execAsync(command);
    return {
      content: [
        {
          type: "text",
          text: `Composer Require:\n${stdout}${stderr ? `\nErrors:\n${stderr}` : ""}`,
        },
      ],
    };
  }

  async run() {
    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error("Laravel MCP server running on stdio");
  }
}

const server = new LaravelMCPServer();
server.run().catch(console.error);