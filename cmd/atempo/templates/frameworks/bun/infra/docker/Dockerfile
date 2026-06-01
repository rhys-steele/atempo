FROM oven/bun:1

WORKDIR /app

COPY src/package.json src/bun.lock* ./
RUN bun install --frozen-lockfile 2>/dev/null || bun install

COPY src/ .

EXPOSE 3000
CMD ["bun", "run", "index.ts"]
