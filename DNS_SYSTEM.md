# Atempo DNS System

This is a fully Dockerized local DNS and reverse proxy setup designed for zero-install developer environments. It enables you to route local applications and services using custom domains like:

- `my-app.local`
- `api.my-app.local`
- `admin.my-other-app.local`

All without installing anything except Docker.

---

## 🚀 What It Does

- Runs a DNS server (`dnsmasq`) and reverse proxy (`nginx`) inside a single Docker container.
- Automatically configures `.local` domain resolution to point to your local services.
- Creates a dedicated Docker network (`atempo-net`) and assigns a static IP (`172.21.0.53`) to the DNS container.
- Configures your system DNS resolver (macOS only) via `/etc/resolver/local` to forward `.local` queries to that static IP.

---

## 🧱 How It Works

- `dnsmasq` reads project-specific domain mappings from `~/.atempo/dns/projects/*.dns`.
- `nginx` proxies HTTP requests to the appropriate service on your machine using matching subdomains.
- DNS and HTTP services are exposed only within the custom Docker network or via configured resolver settings.

---

## 🛠 Installation & Setup

### 1. Build and Run

```bash
go build -o atempo cmd/atempo/main.go
./atempo dns setup
```

This will:
- Create a custom Docker network (`atempo-net`) with subnet `172.21.0.0/24`
- Start the DNS container with static IP `172.21.0.53`
- Configure macOS system resolver to use the containerized DNS
- Flush DNS cache to ensure immediate effect

### 2. Verify Setup

```bash
./atempo dns status
```

Expected output:
```
DNS Configuration
──────────────────────────────────────────────────
✓ DNS service: running
✓ Resolver: configured

Active Domains:
  my-app.local
  my-project.local
```

### 3. Test DNS Resolution

```bash
# Test direct DNS query
dig @172.21.0.53 my-app.local

# Test system resolver
nslookup my-app.local

# Test HTTP proxy
curl -I http://my-app.local
```

---

## 📁 File Structure

```
~/.atempo/dns/
├── dnsmasq.conf              # Main dnsmasq configuration
├── startup.sh                # Container startup script
└── projects/
    ├── my-app.dns            # DNS entries for my-app project
    ├── my-app.nginx          # Nginx proxy config for my-app
    ├── other-project.dns     # DNS entries for other-project
    └── other-project.nginx   # Nginx proxy config for other-project
```

### Example DNS Configuration (`my-app.dns`)

```
address=/my-app.local/127.0.0.1
address=/mysql.my-app.local/127.0.0.1
address=/redis.my-app.local/127.0.0.1
address=/mailhog.my-app.local/127.0.0.1
```

### Example Nginx Configuration (`my-app.nginx`)

```nginx
# Main project domain
server {
    listen 80;
    server_name my-app.local;
    
    location / {
        proxy_pass http://host.docker.internal:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}

# Service subdomains
server {
    listen 80;
    server_name mysql.my-app.local;
    
    location / {
        proxy_pass http://host.docker.internal:3306;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 🔧 Architecture Details

### Container Configuration

- **Image**: `nginx:alpine` (lightweight base with nginx pre-installed)
- **Network**: Custom Docker network `atempo-net` (172.21.0.0/24)
- **Static IP**: `172.21.0.53` (easy to remember DNS IP)
- **Ports**: None exposed to host (uses static IP instead)
- **Volumes**: `~/.atempo/dns:/etc/atempo` (configuration files)

### DNS Server (dnsmasq)

- **Listen Address**: `0.0.0.0:53` (all interfaces, standard DNS port)
- **Configuration**: `/etc/atempo/dnsmasq.conf`
- **Project Configs**: `/etc/atempo/projects/*.dns`
- **Logging**: `/var/log/dnsmasq.log` (query logging enabled)
- **Features**: 
  - `bind-interfaces` for reliable interface binding
  - `log-queries` for debugging
  - Dynamic configuration reload via `SIGHUP`

### HTTP Proxy (nginx)

- **Listen Address**: `0.0.0.0:80`
- **Configuration**: Auto-generated `/etc/nginx/nginx.conf`
- **Project Configs**: `/etc/atempo/projects/*.nginx`
- **Backend**: `host.docker.internal` for host machine access
- **Features**:
  - Graceful configuration reload
  - Fallback 404 server for unmatched domains

### macOS Integration

- **Resolver File**: `/etc/resolver/local`
- **Configuration**: `nameserver 172.21.0.53`
- **Scope**: Only `.local` domain queries are routed to container
- **Cache Management**: Automatic DNS cache flushing after setup

---

## 📋 Commands

### Setup & Management

```bash
# Initial setup (one-time)
./atempo dns setup

# Start DNS service
./atempo dns start

# Stop DNS service  
./atempo dns stop

# Check service status
./atempo dns status

# Test DNS resolution
./atempo dns test
```

### Project Management

```bash
# Add project (done automatically during `atempo create`)
./atempo dns add-project my-app

# Remove project
./atempo dns remove-project my-app

# List active projects
./atempo dns projects
```

---

## 🐛 Troubleshooting

### 1. DNS Not Resolving

**Problem**: `nslookup my-app.local` returns `NXDOMAIN`

**Solutions**:
```bash
# Check if container is running
docker ps | grep atempo-dns

# Verify container has correct IP
docker inspect atempo-dns --format '{{.NetworkSettings.Networks.atempo-net.IPAddress}}'

# Test direct DNS query
dig @172.21.0.53 my-app.local

# Check resolver configuration
cat /etc/resolver/local

# Flush DNS cache
sudo dscacheutil -flushcache
sudo killall -HUP mDNSResponder
```

### 2. Container Won't Start

**Problem**: `./atempo dns start` fails

**Solutions**:
```bash
# Check Docker daemon
docker info

# Check for conflicting containers
docker ps -a | grep atempo

# Remove conflicting containers
docker rm -f atempo-dns

# Check network conflicts
docker network ls | grep atempo
docker network rm atempo-net  # if needed

# View container logs
docker logs atempo-dns
```

### 3. HTTP Proxy Not Working

**Problem**: `curl http://my-app.local` fails

**Solutions**:
```bash
# Verify nginx is running in container
docker exec atempo-dns ps aux | grep nginx

# Test direct connection to container
curl -I http://172.21.0.53

# Check nginx configuration
docker exec atempo-dns nginx -t

# Reload nginx configuration
docker exec atempo-dns nginx -s reload

# Check if backend service is running
curl -I http://localhost:8080  # or actual service port
```

### 4. Permission Issues

**Problem**: `sudo` commands fail during setup

**Solutions**:
```bash
# Ensure you can run sudo commands
sudo whoami

# Manually create resolver file
echo "nameserver 172.21.0.53" | sudo tee /etc/resolver/local

# Fix permissions
sudo chown root:wheel /etc/resolver/local
sudo chmod 644 /etc/resolver/local
```

---

## 🔍 Debugging

### Enable Verbose Logging

```bash
# View DNS query logs
docker exec atempo-dns tail -f /var/log/dnsmasq.log

# View nginx access logs
docker exec atempo-dns tail -f /var/log/nginx/access.log

# View container startup logs
docker logs atempo-dns
```

### Manual DNS Testing

```bash
# Test with different query types
dig @172.21.0.53 my-app.local A
dig @172.21.0.53 my-app.local ANY

# Test TCP DNS (fallback)
dig @172.21.0.53 +tcp my-app.local

# Test from inside container
docker exec atempo-dns nslookup my-app.local 127.0.0.1
```

### Network Diagnostics

```bash
# Check container network configuration
docker inspect atempo-dns | grep -A 20 "Networks"

# Test container connectivity
ping 172.21.0.53

# Check port accessibility
nc -zv 172.21.0.53 53
nc -zv 172.21.0.53 80

# List Docker networks
docker network ls
docker network inspect atempo-net
```

---

## 🚨 Known Issues

### macOS-Specific Issues

1. **UDP Port Forwarding**: Docker Desktop on macOS ARM64 has unreliable UDP port forwarding, which is why we use a custom network with static IP instead of port mapping.

2. **DNS Cache**: macOS aggressively caches DNS results. Always flush cache after configuration changes:
   ```bash
   sudo dscacheutil -flushcache
   sudo killall -HUP mDNSResponder
   ```

3. **`.local` Domain Warning**: You may see warnings about `.local` being reserved for mDNS. This is expected and can be ignored for development.

### Docker-Specific Issues

1. **Network Isolation**: The custom network isolates DNS traffic but requires `host.docker.internal` for proxy backends.

2. **Container Restart**: Configuration changes require either graceful reload or container restart. Graceful reload is preferred and automatic.

3. **IP Address Changes**: If Docker reassigns the network subnet, the static IP may change. Run `./atempo dns setup` to reconfigure.

---

## 🎯 Performance Optimizations

### Graceful Reloads

The system uses graceful reloads instead of container restarts:

```bash
# DNS configuration reload
docker exec atempo-dns pkill -HUP dnsmasq

# Nginx configuration reload  
docker exec atempo-dns nginx -s reload
```

### Efficient Configuration

- **Single Container**: Combines DNS and HTTP proxy for efficiency
- **Static IP**: Eliminates port mapping overhead
- **Bind Interfaces**: Reduces resource usage
- **Query Logging**: Optional for debugging, can be disabled for performance

---

## 🔮 Future Enhancements

1. **Cross-Platform Support**: Extend beyond macOS to Linux and Windows
2. **SSL/TLS Support**: Add automatic certificate generation for HTTPS
3. **Health Checks**: Built-in monitoring and self-healing
4. **Web Interface**: Management UI for configuration
5. **Integration**: Deeper integration with docker-compose projects
6. **Performance**: DNS query caching and optimization

---

## 📚 Technical References

- [dnsmasq Documentation](http://www.thekelleys.org.uk/dnsmasq/docs/dnsmasq-man.html)
- [nginx Documentation](https://nginx.org/en/docs/)
- [Docker Networks](https://docs.docker.com/network/)
- [macOS DNS Resolution](https://developer.apple.com/library/archive/documentation/Darwin/Reference/ManPages/man5/resolver.5.html)