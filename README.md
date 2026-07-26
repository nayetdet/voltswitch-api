# Voltswitch API

API for the Voltswitch project that allows turning a computer on and off remotely over the local network.

The service exposes two endpoints:

- `GET /` for a health check
- `POST /shutdown` to execute the shutdown command defined by `SHUTDOWN_COMMAND`

The API code is in [`main.go`](main.go) and listens on `:3939`.

## Requirements

- Go for local execution
- Docker and Docker Compose for container execution
- A Linux host with an available shutdown command
- Permission to run the container in privileged mode when using the `pid: host` image configuration

## Configuration

The application requires the `SHUTDOWN_COMMAND` environment variable.

Example:

```bash
SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff"
```

The value is executed through `sh -c`, so shell quotes and redirections are supported.

## Running

### Locally

```bash
export SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff"
go mod download
go run .
```

The API is available at:

```text
http://localhost:3939
```

### With Docker

Build the image:

```bash
docker build -t voltswitch-api .
```

Run it manually:

```bash
docker run -d \
  --name voltswitch-api \
  --pid host \
  --privileged \
  -p 3939:3939 \
  -e SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff" \
  voltswitch-api
```

This example uses `nsenter`, so it requires `--pid host` and `--privileged` to reach the host's PID 1. The API is available at `http://localhost:3939`.

### With Docker Compose

The repository includes two Compose files:

- [`docker-compose.passthrough.yml`](docker-compose.passthrough.yml): uses `pid: host` and `privileged: true`
- [`docker-compose.ssh.yml`](docker-compose.ssh.yml): mounts `~/.ssh` and assumes SSH access to the target host

In both cases, set `SHUTDOWN_COMMAND` in the environment before starting the service:

```bash
export SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff"
docker compose -f docker-compose.passthrough.yml up -d
```

Or:

```bash
export SHUTDOWN_COMMAND="ssh user@host poweroff"
docker compose -f docker-compose.ssh.yml up -d
```

The Compose files expose the container's port `3939` as port `3939` on the host.

```text
http://localhost:3939
```

### Helm

The chart in [`k8s/voltswitch-api`](k8s/voltswitch-api) reproduces the SSH Compose setup:

- `Service` on port `3939`, forwarding to the API's port `3939`
- Optional `Ingress` for exposing the API through a local hostname
- Configurable `SHUTDOWN_COMMAND`
- `SSH_PRIVATE_KEY` mounted at `/root/.ssh/id_ed25519` as read-only
- `SSH_KNOWN_HOSTS` mounted at `/root/.ssh/known_hosts` as read-only

By default, the chart synchronizes the SSH credentials through an `ExternalSecret`. The remote Secret must expose `SSH_PRIVATE_KEY` and `SSH_KNOWN_HOSTS`, which are mounted at `/root/.ssh/id_ed25519` and `/root/.ssh/known_hosts`.

Install it with:

```bash
helm upgrade --install voltswitch-api ./k8s/voltswitch-api \
  --set shutdownCommand="ssh user@host poweroff" \
  --set externalSecret.storeName="cluster-secret-store"
```

To use an existing Kubernetes Secret instead, create it with:

```bash
kubectl create secret generic voltswitch-api \
  --from-file=SSH_PRIVATE_KEY="$HOME/.ssh/id_ed25519" \
  --from-file=SSH_KNOWN_HOSTS="$HOME/.ssh/known_hosts"
```

Then install the chart with ExternalSecret disabled:

```bash
helm upgrade --install voltswitch-api ./k8s/voltswitch-api \
  --set shutdownCommand="ssh user@host poweroff" \
  --set secret.enabled=true \
  --set externalSecret.enabled=false
```

The `SSH_PRIVATE_KEY` and `SSH_KNOWN_HOSTS` Secret entries are mounted as `/root/.ssh/id_ed25519` and `/root/.ssh/known_hosts` with read-only permissions.

To enable the Ingress, configure a hostname and an Ingress Controller in your cluster:

```bash
helm upgrade --install voltswitch-api ./k8s/voltswitch-api \
  --set shutdownCommand="ssh user@host poweroff" \
  --set ingress.enabled=true \
  --set 'ingress.hosts[0].host=voltswitch-api.local' \
  --set 'ingress.hosts[0].paths[0].path=/'
```

## Endpoints

### `GET /`

Simple health check.

Response:

```http
204 No Content
```

### `POST /shutdown`

Executes the command configured in `SHUTDOWN_COMMAND` and attempts to shut down the host.

Example:

```bash
curl -X POST http://localhost:3939/shutdown
```

Successful response:

```http
204 No Content
```

Error response:

```json
{
  "error": "error message"
}
```

## Notes

- `SHUTDOWN_COMMAND` is required; without it, the API returns `500`.
- The command is split by the shell before execution.
- The runtime image installs `util-linux` and `openssh-client`, covering the scenarios documented here.

## Security

This API can shut down the host machine. Do not expose it directly to the internet.

Recommendations:

- Restrict access using a firewall, private network, or authenticated proxy.
- Publish the port only on trusted interfaces whenever possible.
- Use it carefully in shared or production environments.
