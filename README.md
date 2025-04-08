# ContainerHub - A Professional Container Management Tool

**ContainerHub** is a lightweight, container management system designed for simplicity and efficiency. Built with Go and leveraging `runc` as its runtime, it provides essential container lifecycle management features with a client-server architecture.

## Features
- **Container Lifecycle**: Create, start, stop, and delete containers.
- **Command Execution**: Run commands inside running containers (`exec`).
- **Container Logs**: Retrieve logs from containers.
- **Listing**: View all containers and their statuses.
- **gRPC-based IPC**: Robust client-server communication over Unix sockets.
- **Minimal Footprint**: Depends only on `runc` and standard Go libraries.

## Prerequisites
- Go 1.18 or later
- `runc` installed (`sudo yum install runc` or equivalent)
- `protoc` for gRPC code generation (`sudo yum install protobuf-compiler` or equivalent)
- Root privileges for daemon and container operations
- Write permissions to `/var/lib/containerhub` and `/tmp`

## Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/shashanksingh24/containerhub.git
   cd containerhub
   ```
2. Build the daemon and client:
   ```bash
   make build
   ```
3. Ensure `/var/lib/containerhub` exists and is writable, if we use socketPath `/tmp/containerhub` then we don't require this:
   
   ```bash
   sudo mkdir -p /var/lib/containerhub
   sudo chmod 755 /var/lib/containerhub
   ```

## Usage
### Step 1: Start the Daemon
The daemon must be running before using the client. Open a terminal and run:
```bash
sudo bin/containerhubd
```
This starts the daemon, listening on `/tmp/containerhub.sock`. Keep this terminal open.

### Step 2: Use the Client
In a separate terminal, interact with the daemon using the `containerhub` CLI:

#### Create a Container
```bash
sudo bin/containerhub create ./rootfs mycontainer "sleep 1000"
```

#### Start a Container
```bash
sudo bin/containerhub start hub_12345678
```

#### Stop a Container
```bash
sudo bin/containerhub stop hub_12345678
```

#### Delete a Container
```bash
sudo bin/containerhub delete hub_12345678
```

#### Execute a Command
```bash
sudo bin/containerhub exec hub_12345678 "echo Hello"
```

#### List Containers
```bash
sudo bin/containerhub list
```

#### View Logs
```bash
sudo bin/containerhub logs hub_12345678
```

### Example Workflow
1. Prepare a root filesystem (e.g., Alpine):
   ```bash
   mkdir rootfs
   curl -O https://dl-cdn.alpinelinux.org/alpine/v3.19/releases/x86_64/alpine-minirootfs-3.19.1-x86_64.tar.gz
   tar -C rootfs -xzf alpine-minirootfs-3.19.1-x86_64.tar.gz
   ```
2. Start the daemon in one terminal:
   ```bash
   sudo bin/containerhubd
   ```
3. In another terminal, manage a container:
   ```bash
   sudo bin/containerhub create ./rootfs testcont "sleep 1000"
   sudo bin/containerhub start hub_12345678
   sudo bin/containerhub exec hub_12345678 "echo Running"
   sudo bin/containerhub logs hub_12345678
   sudo bin/containerhub stop hub_12345678
   sudo bin/containerhub delete hub_12345678
   ```

## Troubleshooting
- **"dial unix /tmp/containerhub.sock: connect: no such file or directory"**: Ensure the daemon (`containerhubd`) is running in another terminal with `sudo bin/containerhubd`.
- **Permission denied**: Run commands with `sudo` and ensure `/tmp` and `/var/lib/containerhub` are writable by the user.
- **runc errors**: Verify `runc` is installed and the rootfs is a valid OCI filesystem.

## Design
- **Daemon (`containerhubd`)**: gRPC server managing container state and `runc` interactions.
- **Client (`containerhub`)**: CLI tool communicating over `/tmp/containerhub.sock`.
- **Storage**: Containers stored in `/var/lib/containerhub/<id>` with OCI `config.json` and logs.

## Limitations
- Manual rootfs preparation (no image pulling).
- Basic networking and volume support not included.
- Logs stored in JSON format via `runc` (requires parsing for human-readable output).

## Development
- Generate protobuf: `make proto`
- Build: `make build`
- Clean: `make clean`

## Contributing
Submit issues or PRs to add features like networking, image management, or enhanced logging.

## License
MIT License - see [LICENSE](LICENSE) for details.
