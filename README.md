## FitGuide Exercises Dataset 💪

Open Exercise Dataset in `JSON` format (870+ exercises) with complete support for English (EN) and Brazilian Portuguese (PT-BR).

### Estrutura do Repositório

```
fitguide-exercises/
├── database/                       <-- Camada de Dados e Conteúdo
│   ├── exercises/                  <-- Exercícios individuais
│   │   ├── pt/                     <-- JSONs e fotos em PT-BR (Fonte da Verdade)
│   │   └── en/                     <-- JSONs e fotos em EN
│   ├── metadata/                   <-- Dicionários e taxonomias
│   │   ├── categories/             <-- pt.json, en.json
│   │   ├── equipments/             <-- pt.json, en.json
│   │   ├── forces/                 <-- pt.json, en.json
│   │   ├── levels/                 <-- pt.json, en.json
│   │   ├── mechanics/              <-- pt.json, en.json
│   │   └── muscles/                <-- pt.json, en.json, anatomy.json
│   ├── workouts/                   <-- Fichas e rotinas de treino
│   │   ├── seeds/                  <-- Treinos padrão iniciais
│   │   └── templates/              <-- Treinos pré-configurados
│   ├── schemas/                    <-- Schemas de validação JSON
│   │   └── exercise.schema.json
│   └── dist/                       <-- Datasets compilados consumidos pelo mobile
│       ├── exercises_pt.json
│       ├── exercises_en.json
│       ├── workouts_pt.json
│       └── workouts_en.json
│
├── admin/                          <-- Painel Administrativo & Servidor
│   ├── web/                        <-- Frontend SPA (Vue 3 + Vite)
│   ├── server/                     <-- Backend REST & Static Server (Go)
│   ├── Dockerfile                  <-- Build Multi-stage
│   └── docker-compose.yml
│
├── scripts/                        <-- Pipelines de Compilação e Automação
│   ├── compile_pt.py
│   ├── compile_en.py
│   ├── extract_metadata.py
│   ├── compile_workouts.py
│   ├── split_pt.py
│   ├── split_en.py
│   └── translate_gym_names.py
│
├── docker-compose.yml              <-- Orquestrador Docker na raiz
├── Makefile                        <-- Comandos de compilação e execução
└── README.md
```

---

### Painel Administrativo (Web & API)

O projeto conta com um painel visual (Vue 3 + Go) para cadastro e edição de exercícios em `database/exercises/pt/`, gerenciamento de treinos/fichas em `database/workouts/`, upload de imagens e disparo manual de compilação dos datasets.

#### 🐳 Executando com Docker (Recomendado)

Certifique-se de que o **Docker Desktop** está em execução em sua máquina.

1. **Subir o painel em segundo plano:**
   ```sh
   docker compose up -d
   # ou via Makefile:
   make docker-up
   ```

2. **Acessar no navegador:**
   Abra [http://localhost:8080](http://localhost:8080).

3. **Parar o container:**
   ```sh
   docker compose down
   # ou via Makefile:
   make docker-down
   ```

> **Nota:** A pasta `database/` é montada como volume, portanto qualquer alteração feita na interface persiste diretamente nos arquivos locais do seu repositório.

---

#### 💻 Executando Localmente (Sem Docker)

Requer **Go 1.22+**, **Node.js 20+** e **Python 3**.

```sh
# Build do frontend e binário Go:
make admin-build

# Iniciar o servidor local:
make admin-run
```
O painel estará disponível em [http://localhost:8080](http://localhost:8080).

---

### Comandos de Terminal & Automação (Makefile)

#### Compilação Completa

```sh
make all
```

#### Compilação Individual

```sh
make database/dist/exercises_en.json
make database/dist/exercises_pt.json
make metadata
make workouts
```

#### Validação Estrutural

```sh
make lint
make check_dupes
```
