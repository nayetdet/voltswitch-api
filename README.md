# voltswitch-api

API HTTP simples para desligar o host onde o container esta rodando.

O endpoint `POST /shutdown` executa o comando definido em `SHUTDOWN_COMMAND`. O valor esperado e o comando completo, por exemplo:

```bash
nsenter --target 1 --mount --uts --ipc --net --pid poweroff
```

O arquivo [`.env.example`](.env.example) traz o valor padrao usado hoje. Crie um `.env` a partir dele quando quiser sobrescrever a configuracao local.

## Requisitos

- Docker e Docker Compose
- Host Linux com `poweroff`
- Permissao para rodar container privilegiado

## Configuracao

Variavel de ambiente usada pela aplicacao:

- `SHUTDOWN_COMMAND`: comando completo executado no shutdown

Exemplo de `.env`:

```bash
SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff"
```

## Como executar

### Com Docker Compose

Defina `SHUTDOWN_COMMAND` no seu ambiente ou em `.env` antes de subir o container.

```bash
docker compose up -d
```

O serviço fica disponivel em:

```text
http://localhost:3939
```

O `docker-compose.yml` publica a porta `8000` do container como `3939` no host.

### Sem Docker

```bash
go mod download
go run .
```

Nesse modo, a API escuta em:

```text
http://localhost:8000
```

## Endpoints

### `GET /`

Health check simples.

Resposta:

```http
204 No Content
```

### `POST /shutdown`

Desliga o host.

Exemplo:

```bash
curl -X POST http://localhost:3939/shutdown
```

Resposta em caso de sucesso:

```http
204 No Content
```

Resposta em caso de erro:

```json
{
  "error": "mensagem do erro"
}
```

## Docker

Build da imagem:

```bash
docker build -t voltswitch-api .
```

Execucao manual equivalente ao Compose:

```bash
docker run -d \
  --name voltswitch-api \
  --pid host \
  --privileged \
  -p 3939:8000 \
  -e SHUTDOWN_COMMAND="nsenter --target 1 --mount --uts --ipc --net --pid poweroff" \
  voltswitch-api
```

## Seguranca

Esta API permite desligar a maquina host. Nao exponha esse servico diretamente na internet.

Recomendacoes:

- Restrinja o acesso por firewall, rede privada ou proxy autenticado.
- Publique a porta apenas em interfaces confiaveis quando possivel.
- Use com cuidado em ambientes compartilhados ou de producao.
