FROM node:18-alpine

# Set working directory
WORKDIR /app

# Install system dependencies for node-gyp and native modules
RUN apk add --no-cache \
    python3 \
    make \
    g++ \
    git \
    curl \
    bash \
    vim \
    nano

# Create app directory and set ownership
RUN mkdir -p /app && chown -R node:node /app

# Copy application code first
COPY --chown=node:node . .

# Switch to node user
USER node

# Install dependencies if package.json exists
RUN if [ -f package.json ]; then npm install; fi

# Create logs directory
RUN mkdir -p logs

# Expose port 3000
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

# Start the application
# Use a flexible command that works whether we have package.json or not
CMD ["sh", "-c", "if [ -f package.json ]; then npm run dev 2>/dev/null || npm start 2>/dev/null || node server.js 2>/dev/null || node index.js 2>/dev/null; else echo 'No package.json found, creating basic server...'; node -e \"const http = require('http'); const server = http.createServer((req, res) => { res.writeHead(200, {'Content-Type': 'application/json'}); res.end(JSON.stringify({status: 'OK', message: 'Express template ready for development'})); }); server.listen(3000, '0.0.0.0', () => console.log('Basic server running on port 3000'));\"; fi"]