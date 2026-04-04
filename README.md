# BabyStories 📚✨

Um aplicativo encantador para leitura e gerenciamento de histórias infantis.
BabyStories usa uma estrutura *Monorepo* contendo uma API robusta em **Go** e uma interface de usuário super moderna em **React 19**.

---

## 🛠 Tecnologias Utilizadas

### Frontend (`apps/web/`)
- **React 19** & **TypeScript**
- **Vite** (Build tool e dev server super rápido)
- **Tailwind CSS** (Estilização *mobile-first*)
- **Shadcn/UI** (Componentes de interface geniais)
- **Jotai** (Gerenciamento atômico de estado)
- **React Router Dom** (Navegação SPA client-side)

### Backend (`cmd/api/` & `internal/`)
- **Go 1.22+**
- **SQLite3** (Banco de dados leve e eficiente)
- **Padrão Repository** (Arquitetura limpa para acesso a dados)

---

## 🚀 Executando o Projeto Localmente

Certifique-se de que você possui **Go 1.22+** e o **Node.js + pnpm** instalados na sua máquina.

### 1. Iniciar a API (Backend)
Na raiz do projeto (`/`), rode o servidor Go:
```bash
go mod tidy
go run ./cmd/api
```
O backend vai gerar um banco `stories.db` e rodar a API em **`http://localhost:8080`**.

### 2. Iniciar a Interface (Frontend)
Em um segundo terminal, instale as dependências com `pnpm` e inicie o Vite:
```bash
pnpm install
pnpm --filter web dev
```
O frontend mágico subirá em **`http://localhost:5173`**.

---

## 📑 Endpoints da API

A interface web se comunica com as seguintes rotas baseadas na URL `/api/v1/`:

| Método | Rota | Descrição |
|---|---|---|
| `GET` | `/stories` | Lista todas as histórias (com suporte a filtros como `?size=small\|large`, `?ai=true\|false`) |
| `POST` | `/stories` | Adiciona uma nova história na base de dados |
| `GET` | `/stories/{id}` | Carrega uma história completa e **incrementa o contador de leituras** |
| `PUT` | `/stories/{id}` | Edita os detalhes de uma história existente |
| `DELETE` | `/stories/{id}` | Exclui uma história permanentemente |
| `GET` | `/health` | Health Check do servidor |

## 📦 Estrutura do Monorepo

```plaintext
stories-go/
│
├── apps/
│   └── web/            # Frontend (React, Vite, Tailwind)
│       ├── src/        # Fontes (Páginas, Componentes, lib/api)
│       └── package.json
│
├── cmd/
│   └── api/            # Ponto de entrada do Backend (Go)
│       └── main.go
│
├── internal/           # Lógica do Backend
│   ├── handlers/       # Controladores HTTP
│   ├── models/         # Estruturas (Story)
│   └── repository/     # Conexão e queries com SQLite
│
├── go.mod              # Dependências Go
├── package.json        # Dependências Globais JS / Husky
├── pnpm-workspace.yaml # Configuração dos Workspaces do pnpm
└── README.md
```

## 🌙 Recursos de Acessibilidade da Interface
Na rota de leitura da história, integramos recursos como:
- **Tamanho Dinâmico**: Aumente ou diminua o texto em 5 escalas para ajuste de miopia ou conforto.
- **Modo Noturno Localizado**: Troque do agradável pergamninho Creme para uma interface Cinza/Escura especialmente voltada a economizar os olhos em leitura no escuro.

---
*Feito com magia e organização para armazenar os melhores sonhos.*
