# voltswitch-api

API HTTP simples para desligar o host onde o container esta rodando.

O projeto expoe um endpoint que executa `poweroff` no namespace do host usando `nsenter`. Por isso, ele foi pensado para rodar em container com `pid: host` e modo privilegiado.

## Requisitos

- Docker e Docker Compose
- Host Linux com `poweroff`
- Permissao para executar containers privilegiados

Para desenvolvimento local sem Docker:

- Go 1.26.2

## Como executar

Com Docker Compose:

```bash
docker compose up -d
```

O servico fica disponivel em:

```text
http://localhost:3939
```

O container escuta na porta `8000`, e o `docker-compose.yml` publica essa porta como `3939` no host.

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

## Desenvolvimento

Instale as dependencias e rode a API localmente:

```bash
go mod download
go run .
```

A API local escuta em:

```text
http://localhost:8000
```

Build local:

```bash
go build -o voltswitch-api .
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
  voltswitch-api
```

## Seguranca

Esta API permite desligar a maquina host. Nao exponha esse servico diretamente na internet.

Recomendacoes:

- Restrinja o acesso por firewall, rede privada ou proxy autenticado.
- Publique a porta apenas em interfaces confiaveis quando possivel.
- Use com cuidado em ambientes compartilhados ou de producao.
