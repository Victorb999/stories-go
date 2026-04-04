# Stories Go

API REST em Go para gerenciar histórias de bebês, com banco de dados SQLite.

## Pré-requisitos

- Go 1.22+

## Executar localmente

```bash
go mod tidy        # baixa dependências
go run ./cmd/api   # inicia o servidor em :8080
```

## Endpoints

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/api/v1/stories` | Lista histórias (`?size=small\|large`, `?ai=true\|false`) |
| `POST` | `/api/v1/stories` | Cria uma história |
| `GET` | `/api/v1/stories/{id}` | Busca por ID (incrementa views) |
| `PUT` | `/api/v1/stories/{id}` | Atualiza uma história |
| `DELETE` | `/api/v1/stories/{id}` | Remove uma história |
| `GET` | `/health` | Health check |

## Exemplo de payload

```json
{
  "title": "O Leão Dorminhoco",
  "cover_image": "https://example.com/lion.jpg",
  "author": "Maria",
  "content": "Era uma vez um leão que adorava dormir...",
  "ai_generated": true,
  "size": "small"
}
```

## Variáveis de ambiente

| Variável | Padrão | Descrição |
|---|---|---|
| `PORT` | `8080` | Porta HTTP |
| `DATA_DIR` | `.` | Diretório onde `stories.db` é criado |

## Deploy (Fly.io)

```bash
fly auth login
fly launch --name stories-go --region gru
fly volumes create stories_data --size 1 --region gru
fly deploy
```
