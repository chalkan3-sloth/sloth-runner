# Sloth-Runner Master-Agent Architecture

`sloth-runner` is designed with a master-agent architecture to facilitate distributed task execution. This allows you to orchestrate and run tasks across multiple remote machines from a central control point.

## Core Concepts

### Master Server

The Master Server is the central component of the `sloth-runner` ecosystem. Its primary responsibilities include:

*   **Agent Registry:** Maintains a registry of all connected and available agents.
*   **Task Orchestration:** Receives task execution requests and dispatches them to the appropriate agents.
*   **Communication Hub:** Acts as the communication hub between the user (via the CLI) and the agents.

### Agent

An Agent is a lightweight process that runs on a remote machine. Its main functions are:

*   **Registration:** Registers itself with the Master Server upon startup, providing its network address and name.
*   **Task Execution:** Receives commands and tasks from the Master Server and executes them locally.
*   **Status Reporting:** Reports the status and output of executed tasks back to the Master Server.

### Communication Protocol

Master and Agents communicate using **gRPC**, a high-performance, open-source universal RPC framework. This ensures efficient and reliable communication between the distributed components.

## Installation and Startup

### Master Server Installation

To set up the `sloth-runner` Master Server, you typically run it on your local machine or a designated control server. The master listens for agent connections on a specified port.

**Command:**

```bash
go run ./cmd/sloth-runner master -p <port> [--daemon]
```

*   `-p, --port <port>`: Specifies the port on which the master server will listen for agent connections. The default port is `50053`.
*   `--daemon`: (Optional) Runs the master server as a background daemon process. This is recommended for continuous operation.

**Example:**

To start the master server on port `50053` in daemon mode:

```bash
go run ./cmd/sloth-runner master -p 50053 --daemon
```

Upon successful startup, the master will log that it is listening for agent registrations.

### Agent Installation

Agents are deployed on the remote machines where you intend to execute tasks. Each agent needs to be configured with a unique name and the address of the Master Server.

**Command:**

```bash
sloth-runner agent start --name <agent_name> --master <master_ip>:<master_port> --port <agent_port> --bind-address <agent_ip> [--daemon]
```

*   `--name <agent_name>`: A unique name for this agent (e.g., `agent1`, `web-server-agent`). This name is used by the master to identify and address the agent.
*   `--master <master_ip>:<master_port>`: The IP address and port of the running Master Server. Agents will connect to this address to register and receive tasks.
*   `--port <agent_port>`: The port on which the agent itself will listen for direct communication from the master (e.g., for task execution requests). The default port is `50051`.
*   `--bind-address <agent_ip>`: **Crucial for remote agents.** This specifies the specific IPv4 address that the agent should bind to and report to the master. This ensures the master can correctly connect to the agent, especially in environments with multiple network interfaces or IPv6 preference. **Always set this to the remote machine's accessible IPv4 address.**
*   `--daemon`: (Optional) Runs the agent as a background daemon process.

**Example:**

To start an agent named `agent1` on a machine with IP `192.168.1.16`, connecting to a master at `192.168.1.21:50053`, and listening on port `50051`:

```bash
sloth-runner agent start --name agent1 --master 192.168.1.21:50053 --port 50051 --bind-address 192.168.1.16 --daemon
```

## Task Execution Workflow

1.  **Master Startup:** The `sloth-runner` master server starts and begins listening for agent registrations.
2.  **Agent Startup & Registration:** An agent starts on a remote machine, connects to the configured master, and registers itself, providing its unique name and accessible network address.
3.  **Agent Listing:** The user can list all registered agents using `sloth-runner agent list` from the master's machine.
4.  **Task Request:** The user initiates a task execution on a specific agent using `sloth-runner agent run <agent_name> <command>`.
5.  **Task Dispatch:** The master receives the request, looks up the agent's address in its registry, and dispatches the command to the target agent via gRPC.
6.  **Task Execution:** The agent receives the command, executes it locally (e.g., using `bash -c <command>`), and captures its standard output, standard error, and exit status.
7.  **Result Reporting:** The agent sends the execution results (stdout, stderr, success/failure) back to the master.
8.  **Output Presentation:** The master receives the results and presents them to the user in a clear, formatted, and colored output (as described in the [Enhanced `sloth-runner agent run` Output](enhanced-agent-output.md) documentation).

This architecture provides a flexible and scalable way to manage and execute tasks across your infrastructure.

## Special Configurations

### Agents in Incus/LXC Containers

When deploying agents inside Incus (or LXC) containers, you need to configure port forwarding and use the `--report-address` flag because the container's internal IP is not accessible from the master.

#### Quick Start

For a fast setup in an Incus container:

```bash
# 1. On the HOST - Configure port forwarding
sudo incus config device add main sloth-proxy proxy \
  listen=tcp:0.0.0.0:50052 \
  connect=tcp:127.0.0.1:50051

# 2. In the CONTAINER - Install with bootstrap script
sudo incus exec main -- bash -c "curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh | bash -s -- --name main --master 192.168.1.29:50053 --incus 192.168.1.17:50052"

# Done! The agent is now running and configured.
```

#### Setup Steps

1. **Configure Port Forwarding on the Host**

   Add a proxy device to forward a host port to the container's agent port:

   ```bash
   # On the host machine running Incus
   sudo incus config device add <container_name> sloth-proxy proxy \
     listen=tcp:0.0.0.0:<host_port> \
     connect=tcp:127.0.0.1:<agent_port>
   ```

   **Example:**
   ```bash
   sudo incus config device add main sloth-proxy proxy \
     listen=tcp:0.0.0.0:50052 \
     connect=tcp:127.0.0.1:50051
   ```

2. **Start Agent with Report Address**

   Inside the container, start the agent with:

   **Option A: Using Bootstrap Script (Recommended)**

   ```bash
   # Inside the container
   bash <(curl -fsSL https://raw.githubusercontent.com/chalkan3-sloth/sloth-runner/master/bootstrap.sh) \
     --name <agent_name> \
     --master <master_ip>:<master_port> \
     --incus <host_ip>:<host_port>
   ```

   The `--incus` flag automatically sets:
   - `--bind-address 0.0.0.0` (listen on all interfaces)
   - `--report-address <host_ip>:<host_port>` (master connects via host)
   - Creates and enables systemd service

   **Option B: Manual Configuration**

   - `--bind-address 0.0.0.0` to listen on all interfaces
   - `--report-address <host_ip>:<host_port>` to tell the master how to reach this agent

   ```bash
   # Inside the container
   sloth-runner agent start \
     --name <agent_name> \
     --master <master_ip>:<master_port> \
     --port <agent_port> \
     --bind-address 0.0.0.0 \
     --report-address <host_ip>:<host_port> \
     --daemon
   ```

   **Example:**
   ```bash
   # Inside container "main" on host 192.168.1.17
   sloth-runner agent start \
     --name main \
     --master 192.168.1.29:50053 \
     --port 50051 \
     --bind-address 0.0.0.0 \
     --report-address 192.168.1.17:50052 \
     --daemon
   ```

3. **Systemd Service Configuration (Recommended)**

   Create a systemd service file at `/etc/systemd/system/sloth-runner-agent.service`:

   ```ini
   [Unit]
   Description=Sloth Runner Agent - <agent_name>
   Documentation=https://chalkan3.github.io/sloth-runner/
   After=network-online.target
   Wants=network-online.target

   [Service]
   Type=simple
   User=root
   WorkingDirectory=/var/lib/sloth-runner
   Restart=always
   RestartSec=5s
   StartLimitInterval=60s
   StartLimitBurst=5

   # Agent Configuration
   ExecStart=/usr/local/bin/sloth-runner agent start \
     --name <agent_name> \
     --master <master_ip>:<master_port> \
     --port <agent_port> \
     --bind-address 0.0.0.0 \
     --report-address <host_ip>:<host_port>

   # Logging
   StandardOutput=journal
   StandardError=journal
   SyslogIdentifier=sloth-runner-agent

   # Performance
   LimitNOFILE=65536

   # Security
   NoNewPrivileges=true
   PrivateTmp=true

   [Install]
   WantedBy=multi-user.target
   ```

   Then enable and start the service:

   ```bash
   systemctl daemon-reload
   systemctl enable sloth-runner-agent
   systemctl start sloth-runner-agent
   ```

#### Port Mapping Summary

| Component | Internal IP:Port | Exposed Host IP:Port | Master Sees |
|-----------|------------------|---------------------|-------------|
| Container Agent | 10.x.x.x:50051 | host_ip:50052 | host_ip:50052 |
| Host Agent | host_ip:50051 | host_ip:50051 | host_ip:50051 |

#### Troubleshooting

**Agent shows as "Active" but commands timeout:**
- Verify port forwarding is configured: `incus config device list <container_name>`
- Check the agent is using `--report-address` with the host's IP and forwarded port
- Test connectivity: `nc -zv <host_ip> <host_port>` from the master machine

**Multiple containers on the same host:**
- Use different host ports for each container (e.g., 50052, 50053, 50054)
- Update each agent's `--report-address` accordingly 
