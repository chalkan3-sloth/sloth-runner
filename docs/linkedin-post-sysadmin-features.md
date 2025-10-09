🚀 Stop SSHing Into Dozens of Servers!

BEFORE:
ssh web-01 "tail -f /var/log/syslog"
ssh web-02 "journalctl -f"
ssh web-03 "tail /var/log/nginx/error.log"
... repeat 50x 😫

NOW:
sloth-runner sysadmin logs remote --agent web-01 --system syslog

⚡ THE IMPACT:

Before → Now
⏰ 5 minutes → 30 seconds
🔐 10 SSH keys → Zero config
😫 High friction → Zero friction

💎 HOW IT WORKS:

Remote logs via gRPC:
• No interactive SSH
• Works through firewalls
• Supports: syslog, journalctl, auth, kern, custom paths

REAL EXAMPLE - Incident Response:

1. Check health (3s)
sloth-runner sysadmin health agent web-01

2. View logs (5s)
sloth-runner sysadmin logs remote --agent web-01 --system syslog --lines 100

3. Check auth (5s)
sloth-runner sysadmin logs remote --agent web-01 --system auth | grep "failed"

4. Monitor real-time (17s)
sloth-runner sysadmin logs remote --agent web-01 --follow

Total: 30 seconds. SSH required: 0. ✅

🎁 BONUS FEATURES:

Health checks:
sloth-runner sysadmin health check
sloth-runner sysadmin health agent --all

Log management:
sloth-runner sysadmin logs search --query "error" --since 1h
sloth-runner sysadmin logs export --format json

🚀 GETTING STARTED:

Install:
go install github.com/chalkan3-sloth/sloth-runner@latest

Use:
sloth-runner sysadmin logs remote --agent my-server --system syslog

💬 FOR YOU:

Managing multiple servers?
Spending time SSHing to check logs?

You need to know sloth-runner! 🎯

⭐ GitHub: github.com/chalkan3-sloth/sloth-runner
📖 Complete docs in PT-BR

---

Tag that SRE/DevOps who lives SSHing into servers! 👇

#DevOps #SRE #Infrastructure #Automation #LogManagement #Golang #OpenSource #Observability #CloudNative #DistributedSystems #RemoteManagement #InfrastructureAsCode #SiteReliabilityEngineering

Built with ❤️ for the SRE/DevOps community
