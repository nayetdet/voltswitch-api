# voltswitch-api

API HTTP simples para desligar o host onde o container está rodando.

O serviço expõe dois endpoints:

- `GET /` para health check
- `POST /shutdown` para executar o comando de desligamento definido em `SHUTDOWN_COMMAND`

O código da API está em [`main.go`](main.go) e escuta em `:8000`.

## Requisitos

- Go `1.26.2` para execução local
- Docker e Docker Compose para execução em container
- Host Linux com o comando de desligamento disponível
- Permissão para rodar o container em modo privilegiado, quando usar a imagem com `pid: host`

## Configuração

A aplicação depende da variável de ambiente `SHUTDOWN_COMMAND`.

Exemplo:

```bash
SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff"
```

O valor é executado diretamente pelo processo, então passe o comando completo e evite depender de parsing avançado de shell.

## Como executar

### Localmente

```bash
export SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff"
go mod download
go run .
```

A API fica disponível em:

```text
http://localhost:8000
```

### Com Docker

Build da imagem:

```bash
docker build -t voltswitch-api .
```

Execução manual:

```bash
docker run -d \
  --name voltswitch-api \
  --pid host \
  --privileged \
  -p 8000:8000 \
  -e SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff" \
  voltswitch-api
```

Esse exemplo usa `nsenter`, então precisa de `--pid host` e `--privileged` para alcançar o PID 1 do host.

### Com Docker Compose

O repositório inclui dois arquivos:

- [`docker-compose.passthrough.yml`](docker-compose.passthrough.yml): usa `pid: host` e `privileged: true`
- [`docker-compose.ssh.yml`](docker-compose.ssh.yml): monta `~/.ssh` e parte do pressuposto de acesso por SSH

Em ambos os casos, defina `SHUTDOWN_COMMAND` no ambiente antes de subir o serviço:

```bash
export SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff"
docker compose -f docker-compose.passthrough.yml up -d
```

ou:

```bash
export SHUTDOWN_COMMAND="ssh user@host poweroff"
docker compose -f docker-compose.ssh.yml up -d
```

Os Compose expõem a porta `8000` do container como `3939` no host.

```text
http://localhost:3939
```

## Endpoints

### `GET /`

Health check simples.

Resposta:

```http
204 No Content
```

### `POST /shutdown`

Executa o comando configurado em `SHUTDOWN_COMMAND` e tenta desligar o host.

Exemplo:

```bash
curl -X POST http://localhost:3939/shutdown
```

Resposta de sucesso:

```http
204 No Content
```

Resposta de erro:

```json
{
  "error": "mensagem do erro"
}
```

## Observações

- `SHUTDOWN_COMMAND` é obrigatório; sem ele a API responde `500`.
- O comando é dividido por espaços antes de ser executado.
- A imagem runtime instala `util-linux` e `openssh-client`, o que cobre os cenários documentados aqui.

## Segurança

Esta API permite desligar a máquina host. Não exponha esse serviço diretamente na internet.

Recomendações:

- Restrinja o acesso por firewall, rede privada ou proxy autenticado.
- Publique a porta apenas em interfaces confiáveis quando possível.
- Use com cuidado em ambientes compartilhados ou de produção.
