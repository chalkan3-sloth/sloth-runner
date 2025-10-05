# 🚀 Telemetry Deployment Guide

## Overview

This guide covers deploying telemetry in production, integrating with Prometheus, and setting up Grafana dashboards.

## Quick Start

### Enable Telemetry

Telemetry is **enabled by default** in Sloth Runner. Simply start your agent:

```bash
./sloth-runner agent start --name my-agent --master master-host:50053
```

### Custom Configuration

```bash
./sloth-runner agent start \
  --name my-agent \
  --master master-host:50053 \
  --metrics-port 9090 \          # Custom metrics port
  --telemetry                     # Explicitly enable
```

### Disable Telemetry

```bash
./sloth-runner agent start \
  --name my-agent \
  --master master-host:50053 \
  --telemetry=false
```

---

## Prometheus Integration

### Configure Prometheus Scraping

#### Static Configuration

Create or update `prometheus.yml`:

```yaml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

scrape_configs:
  - job_name: 'sloth-runner-agents'
    static_configs:
      - targets:
          - 'agent1.example.com:9090'
          - 'agent2.example.com:9090'
          - 'agent3.example.com:9090'
        labels:
          environment: 'production'
          cluster: 'main'
```

#### Service Discovery

=== "Kubernetes"
    ```yaml
    scrape_configs:
      - job_name: 'sloth-runner-k8s'
        kubernetes_sd_configs:
          - role: pod
        relabel_configs:
          - source_labels: [__meta_kubernetes_pod_label_app]
            action: keep
            regex: sloth-runner-agent
          - source_labels: [__meta_kubernetes_pod_ip]
            action: replace
            target_label: __address__
            replacement: '$1:9090'
    ```

=== "Consul"
    ```yaml
    scrape_configs:
      - job_name: 'sloth-runner-consul'
        consul_sd_configs:
          - server: 'consul.example.com:8500'
            services: ['sloth-runner-agent']
        relabel_configs:
          - source_labels: [__meta_consul_service]
            action: keep
            regex: sloth-runner-agent
    ```

=== "File SD"
    ```yaml
    scrape_configs:
      - job_name: 'sloth-runner-file'
        file_sd_configs:
          - files:
              - '/etc/prometheus/targets/sloth-runner-*.json'
            refresh_interval: 30s
    ```

    Create `/etc/prometheus/targets/sloth-runner-prod.json`:
    ```json
    [
      {
        "targets": [
          "agent1:9090",
          "agent2:9090"
        ],
        "labels": {
          "environment": "production",
          "datacenter": "us-east-1"
        }
      }
    ]
    ```

### Verify Scraping

Check Prometheus targets:

```
http://prometheus-host:9090/targets
```

Query metrics:

```promql
up{job="sloth-runner-agents"}
```

Expected output:
```
up{instance="agent1:9090",job="sloth-runner-agents"} 1
up{instance="agent2:9090",job="sloth-runner-agents"} 1
```

---

## Grafana Integration

### Import Dashboard

#### Option 1: From JSON

1. Download the dashboard JSON from GitHub:
   ```bash
   curl -O https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/grafana-dashboard.json
   ```

2. In Grafana UI:
   - Navigate to **Dashboards** → **Import**
   - Upload `grafana-dashboard.json`
   - Select Prometheus data source
   - Click **Import**

#### Option 2: Manual Creation

Create a new dashboard with these panels:

=== "Task Success Rate"
    ```promql
    sum(rate(sloth_tasks_total{status="success"}[5m]))
    /
    sum(rate(sloth_tasks_total[5m]))
    * 100
    ```

    - **Type**: Stat
    - **Unit**: Percent (0-100)
    - **Thresholds**: 95 (yellow), 98 (green)

=== "Task Execution Rate"
    ```promql
    sum(rate(sloth_tasks_total[5m])) by (status)
    ```

    - **Type**: Graph
    - **Legend**: `{{status}}`
    - **Stack**: Yes

=== "Task P99 Latency"
    ```promql
    histogram_quantile(0.99,
      sum(rate(sloth_task_duration_seconds_bucket[5m])) by (task, le)
    )
    ```

    - **Type**: Graph
    - **Legend**: `{{task}}`
    - **Unit**: seconds (s)

=== "Memory Usage"
    ```promql
    sloth_memory_allocated_bytes / 1024 / 1024
    ```

    - **Type**: Graph
    - **Unit**: MiB
    - **Thresholds**: 400 (yellow), 500 (red)

=== "Active Agents"
    ```promql
    count(up{job="sloth-runner-agents"} == 1)
    ```

    - **Type**: Stat
    - **Color**: Value-based

### Dashboard Template

Full dashboard configuration:

```json
{
  "dashboard": {
    "title": "Sloth Runner - Agent Fleet",
    "tags": ["sloth-runner", "automation"],
    "timezone": "browser",
    "panels": [
      {
        "title": "Active Agents",
        "targets": [
          {
            "expr": "count(up{job=\"sloth-runner-agents\"} == 1)"
          }
        ],
        "type": "stat"
      },
      {
        "title": "Task Success Rate",
        "targets": [
          {
            "expr": "sum(rate(sloth_tasks_total{status=\"success\"}[5m])) / sum(rate(sloth_tasks_total[5m])) * 100"
          }
        ],
        "type": "gauge",
        "fieldConfig": {
          "defaults": {
            "unit": "percent",
            "thresholds": {
              "steps": [
                { "value": 0, "color": "red" },
                { "value": 95, "color": "yellow" },
                { "value": 98, "color": "green" }
              ]
            }
          }
        }
      }
    ]
  }
}
```

---

## Docker Deployment

### Docker Compose

Complete monitoring stack with Sloth Runner:

```yaml
version: '3.8'

services:
  # Sloth Runner Agent
  sloth-agent:
    image: slothrunner/agent:latest
    container_name: sloth-agent-1
    command:
      - agent
      - start
      - --name=agent-1
      - --master=sloth-master:50053
      - --telemetry
      - --metrics-port=9090
    ports:
      - "9090:9090"  # Metrics port
    networks:
      - monitoring
    restart: unless-stopped

  # Prometheus
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--storage.tsdb.retention.time=30d'
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus-data:/prometheus
    ports:
      - "9091:9090"
    networks:
      - monitoring
    restart: unless-stopped

  # Grafana
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    volumes:
      - grafana-data:/var/lib/grafana
      - ./grafana-dashboards:/etc/grafana/provisioning/dashboards
      - ./grafana-datasources.yml:/etc/grafana/provisioning/datasources/datasources.yml
    ports:
      - "3000:3000"
    networks:
      - monitoring
    restart: unless-stopped

volumes:
  prometheus-data:
  grafana-data:

networks:
  monitoring:
```

`prometheus.yml` for Docker Compose:

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'sloth-runner'
    static_configs:
      - targets: ['sloth-agent:9090']
```

`grafana-datasources.yml`:

```yaml
apiVersion: 1

datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: http://prometheus:9090
    isDefault: true
```

Start the stack:

```bash
docker-compose up -d
```

Access:
- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9091
- Metrics: http://localhost:9090/metrics

---

## Kubernetes Deployment

### Agent DaemonSet

Deploy agents as DaemonSet:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: sloth-runner-agent
  namespace: automation
spec:
  selector:
    matchLabels:
      app: sloth-runner-agent
  template:
    metadata:
      labels:
        app: sloth-runner-agent
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "9090"
        prometheus.io/path: "/metrics"
    spec:
      containers:
        - name: agent
          image: slothrunner/agent:v1.2.3
          args:
            - agent
            - start
            - --name=$(NODE_NAME)
            - --master=sloth-master.automation.svc.cluster.local:50053
            - --telemetry
            - --metrics-port=9090
          env:
            - name: NODE_NAME
              valueFrom:
                fieldRef:
                  fieldPath: spec.nodeName
          ports:
            - name: metrics
              containerPort: 9090
              protocol: TCP
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi
          livenessProbe:
            httpGet:
              path: /health
              port: 9090
            initialDelaySeconds: 10
            periodSeconds: 30
          readinessProbe:
            httpGet:
              path: /health
              port: 9090
            initialDelaySeconds: 5
            periodSeconds: 10
```

### ServiceMonitor

For Prometheus Operator:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: sloth-runner-agents
  namespace: automation
  labels:
    app: sloth-runner
spec:
  selector:
    matchLabels:
      app: sloth-runner-agent
  endpoints:
    - port: metrics
      interval: 15s
      path: /metrics
```

### Grafana Dashboard ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: sloth-runner-dashboard
  namespace: monitoring
  labels:
    grafana_dashboard: "1"
data:
  sloth-runner.json: |
    {
      "dashboard": {
        "title": "Sloth Runner - Kubernetes Fleet",
        ...
      }
    }
```

---

## Network Configuration

### Firewall Rules

#### iptables

Allow metrics port:

```bash
# Allow from Prometheus server
sudo iptables -A INPUT -p tcp -s prometheus-ip --dport 9090 -j ACCEPT

# Allow from monitoring subnet
sudo iptables -A INPUT -p tcp -s 10.0.0.0/24 --dport 9090 -j ACCEPT

# Save rules
sudo iptables-save > /etc/iptables/rules.v4
```

#### firewalld

```bash
# Add metrics port
sudo firewall-cmd --permanent --add-port=9090/tcp

# Or create service
sudo firewall-cmd --permanent --new-service=sloth-metrics
sudo firewall-cmd --permanent --service=sloth-metrics --add-port=9090/tcp
sudo firewall-cmd --permanent --add-service=sloth-metrics

# Reload
sudo firewall-cmd --reload
```

#### ufw

```bash
# Allow from specific IP
sudo ufw allow from prometheus-ip to any port 9090

# Allow from subnet
sudo ufw allow from 10.0.0.0/24 to any port 9090
```

### Reverse Proxy

For auth and TLS termination:

=== "Nginx"
    ```nginx
    server {
        listen 443 ssl;
        server_name metrics.example.com;

        ssl_certificate /etc/ssl/certs/metrics.crt;
        ssl_certificate_key /etc/ssl/private/metrics.key;

        location /metrics {
            auth_basic "Metrics";
            auth_basic_user_file /etc/nginx/.htpasswd;

            proxy_pass http://localhost:9090/metrics;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }

        location /health {
            proxy_pass http://localhost:9090/health;
            allow 10.0.0.0/24;
            deny all;
        }
    }
    ```

=== "Caddy"
    ```caddyfile
    metrics.example.com {
        reverse_proxy /metrics localhost:9090 {
            basic_auth {
                prometheus $2a$14$...
            }
        }

        reverse_proxy /health localhost:9090 {
            @denied not remote_ip 10.0.0.0/24
            respond @denied 403
        }
    }
    ```

---

## Security

### Authentication

Prometheus doesn't support native auth. Use reverse proxy:

```bash
# Create htpasswd file
htpasswd -c /etc/nginx/.htpasswd prometheus

# Configure Nginx (see above)
```

Update Prometheus config:

```yaml
scrape_configs:
  - job_name: 'sloth-runner-secure'
    basic_auth:
      username: prometheus
      password: your-password
    static_configs:
      - targets: ['agent:443']
    scheme: https
```

### TLS

Generate self-signed cert:

```bash
openssl req -x509 -newkey rsa:4096 \
  -keyout key.pem -out cert.pem \
  -days 365 -nodes \
  -subj "/CN=agent.example.com"
```

Configure reverse proxy with TLS (see Nginx example above).

### Network Isolation

Best practices:

1. **Private Network**: Deploy agents in private subnet
2. **VPN**: Access metrics through VPN
3. **SSH Tunnel**: For ad-hoc access:
   ```bash
   ssh -L 9090:localhost:9090 agent-host
   # Access at http://localhost:9090/metrics
   ```

---

## Monitoring the Monitors

### Prometheus Self-Monitoring

Alert on scrape failures:

```yaml
groups:
  - name: monitoring
    rules:
      - alert: SlothAgentDown
        expr: up{job="sloth-runner-agents"} == 0
        for: 1m
        annotations:
          summary: "Sloth agent {{ $labels.instance }} is down"

      - alert: SlothAgentScrapeFailed
        expr: up{job="sloth-runner-agents"} == 0
        for: 5m
        annotations:
          summary: "Cannot scrape {{ $labels.instance }}"
```

### Health Checks

Monitor telemetry health:

```bash
# Simple health check script
#!/bin/bash
AGENT_HOST="agent.example.com"
METRICS_PORT="9090"

# Check health endpoint
if curl -sf http://$AGENT_HOST:$METRICS_PORT/health > /dev/null; then
  echo "✓ Telemetry is healthy"
  exit 0
else
  echo "✗ Telemetry is down"
  exit 1
fi
```

Add to cron or monitoring system:

```cron
*/5 * * * * /usr/local/bin/check-telemetry.sh || /usr/local/bin/alert-ops.sh
```

---

## Performance Tuning

### Metrics Cardinality

Monitor label cardinality:

```promql
# Count unique label combinations
count(sloth_tasks_total) by (__name__)
```

Best practices:

- ✅ Use `group` label for task groups
- ✅ Use `task` label for individual tasks
- ❌ Don't use high-cardinality labels (user IDs, timestamps, etc.)
- ❌ Don't create metrics for every unique value

### Scrape Interval

Recommendations:

| Environment | Scrape Interval | Retention |
|-------------|----------------|-----------|
| Development | 5s | 7 days |
| Staging | 15s | 15 days |
| Production | 15-30s | 30-90 days |

### Resource Limits

Telemetry resource usage:

```yaml
# Kubernetes resources
resources:
  requests:
    cpu: 50m
    memory: 64Mi
  limits:
    cpu: 200m
    memory: 256Mi
```

---

## Troubleshooting

### Metrics Not Appearing

1. **Check agent logs**:
   ```bash
   tail -f agent.log | grep telemetry
   ```

2. **Verify endpoint**:
   ```bash
   curl http://localhost:9090/metrics
   ```

3. **Check Prometheus targets**:
   Navigate to `http://prometheus:9090/targets`

4. **Validate config**:
   ```bash
   promtool check config prometheus.yml
   ```

### High Memory Usage

If telemetry uses too much memory:

1. **Reduce scrape interval**: Change from 15s to 30s or 60s
2. **Limit metric labels**: Remove unnecessary labels
3. **Increase retention**: Allow Prometheus to aggregate older data

### Connection Issues

Test connectivity:

```bash
# From Prometheus host
telnet agent-host 9090

# Test scrape
curl -v http://agent-host:9090/metrics

# Check firewall
nmap -p 9090 agent-host
```

---

## Best Practices

### Production Checklist

- [ ] Telemetry enabled on all agents
- [ ] Prometheus scraping configured
- [ ] Grafana dashboards imported
- [ ] Alerts configured
- [ ] Firewall rules applied
- [ ] TLS/auth configured (if needed)
- [ ] Backup Prometheus data
- [ ] Document runbooks

### Monitoring Strategy

1. **Real-time**: Terminal dashboard for immediate feedback
2. **Short-term**: Prometheus for recent trends (1-7 days)
3. **Long-term**: Export to long-term storage (S3, BigQuery)

### Alert Guidelines

| Metric | Threshold | Action |
|--------|-----------|--------|
| Task failure rate | > 10% | Investigate failing tasks |
| gRPC latency P99 | > 1s | Check network/master |
| Memory usage | > 80% | Scale or optimize |
| Agent down | > 1m | Restart agent |

---

## Next Steps

- [Prometheus Metrics Reference](prometheus-metrics.md) - Learn about available metrics
- [Grafana Dashboard Guide](grafana-dashboard.md) - Use the terminal dashboard
- [Telemetry Overview](index.md) - Back to overview

## External Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Operator](https://prometheus-operator.dev/)
- [Best Practices for Monitoring](https://prometheus.io/docs/practices/)
