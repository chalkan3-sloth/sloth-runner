# Telemetry & Observability Documentation

This directory contains comprehensive documentation for Sloth Runner's telemetry and observability features.

## Documentation Structure

### 📊 [Overview](index.md)
**Size**: 7.6 KB | **Path**: `docs/en/telemetry/index.md`

Complete introduction to telemetry features including:
- Prometheus integration overview
- Terminal dashboard introduction
- Quick start guide
- Use cases and architecture
- Troubleshooting basics

**Start here** if you're new to Sloth Runner telemetry.

---

### 📈 [Prometheus Metrics Reference](prometheus-metrics.md)
**Size**: 11 KB | **Path**: `docs/en/telemetry/prometheus-metrics.md`

Detailed reference of all available metrics:
- Complete metric catalog with types (Counter, Gauge, Histogram)
- Label descriptions
- Example outputs
- PromQL query examples
- Best practices for querying
- Alert rule examples
- Recording rules

**Reference this** when building dashboards or writing alerts.

---

### 💻 [Grafana-Style Terminal Dashboard](grafana-dashboard.md)
**Size**: 19 KB | **Path**: `docs/en/telemetry/grafana-dashboard.md`

Complete guide to the terminal dashboard:
- All dashboard sections explained
- Color coding reference
- Use cases and workflows
- Watch mode features
- Advanced usage patterns
- Troubleshooting dashboard issues

**Use this** to learn how to use the `agent metrics grafana` command.

---

### 🚀 [Deployment Guide](deployment.md)
**Size**: 16 KB | **Path**: `docs/en/telemetry/deployment.md`

Production deployment and integration:
- Prometheus configuration (static, service discovery)
- Grafana dashboard setup
- Docker Compose examples
- Kubernetes deployment (DaemonSet, ServiceMonitor)
- Network configuration (firewall, reverse proxy)
- Security (TLS, authentication)
- Performance tuning
- Monitoring best practices

**Follow this** for production deployments.

---

## Quick Links

### Common Tasks

| Task | Documentation | Command Example |
|------|---------------|----------------|
| Enable telemetry | [Overview](index.md#enable-telemetry-on-agent) | `--telemetry --metrics-port 9090` |
| Get metrics URL | [Overview](index.md#access-metrics) | `agent metrics prom <name>` |
| View dashboard | [Dashboard Guide](grafana-dashboard.md#quick-start) | `agent metrics grafana <name>` |
| Watch mode | [Dashboard Guide](grafana-dashboard.md#watch-mode) | `agent metrics grafana <name> --watch` |
| Prometheus setup | [Deployment](deployment.md#prometheus-integration) | See prometheus.yml examples |
| Grafana setup | [Deployment](deployment.md#grafana-integration) | Import dashboard JSON |

### By Role

| Role | Recommended Reading |
|------|-------------------|
| **Developer** | [Overview](index.md) → [Dashboard Guide](grafana-dashboard.md) |
| **DevOps/SRE** | [Overview](index.md) → [Deployment Guide](deployment.md) → [Metrics Reference](prometheus-metrics.md) |
| **Data Analyst** | [Metrics Reference](prometheus-metrics.md) → [Deployment Guide](deployment.md#grafana-integration) |
| **Security** | [Deployment Guide](deployment.md#security) |

## Language Versions

- 🇺🇸 **English**: `docs/en/telemetry/`
- 🇧🇷 **Português**: `docs/pt/telemetria/` (Overview only)

## Related Documentation

- [Agent Setup Guide](../master-agent-architecture.md) - How to set up agents
- [Distributed Agents](../distributed.md) - Multi-agent deployments
- [CLI Reference](../CLI.md) - Command-line interface
- [Enterprise Features](../enterprise-features.md) - Other enterprise capabilities

## External Resources

- [Prometheus Official Docs](https://prometheus.io/docs/)
- [Grafana Official Docs](https://grafana.com/docs/)
- [pterm Library](https://github.com/pterm/pterm) - Terminal UI library used
- [Prometheus Operator](https://prometheus-operator.dev/) - For Kubernetes

## Contributing

To improve this documentation:

1. Edit the Markdown files directly
2. Follow the existing structure and style
3. Test examples before documitting
4. Update the navigation in `mkdocs.yml`
5. Submit a pull request

## Document Statistics

| File | Size | Lines | Topics Covered |
|------|------|-------|----------------|
| index.md | 7.6 KB | ~200 | Overview, Quick Start, Architecture |
| prometheus-metrics.md | 11 KB | ~350 | All metrics, PromQL, Alerts |
| grafana-dashboard.md | 19 KB | ~600 | Dashboard sections, Use cases |
| deployment.md | 16 KB | ~500 | Production, Security, K8s |
| **Total** | **53.6 KB** | **~1,650** | **Comprehensive telemetry** |

Last updated: 2025-10-05
