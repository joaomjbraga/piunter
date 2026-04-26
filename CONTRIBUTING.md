# Contributing to piunter

Obrigado pelo seu interesse em contribuir para o piunter!

## Desenvolvimento - Versão Go (Recomendada)

### Setup

```bash
# Clone o repositório
git clone https://github.com/joaomjbraga/piunter.git
cd piunter/piunter-cli-go

# Instale dependências
go mod download

# Build
go build -o piunter ./cmd/main.go

# Execute
./piunter --help
```

### Padrões de Código

- Use Go para todo novo código
- Execute `go fmt` antes de commitar
- Execute `go vet` para verificar erros
- Escreva testes para novas funcionalidades

### Estrutura do Projeto

```
piunter-cli-go/
├── cmd/main.go           # Entry point + CLI (cobra)
├── pkg/types/types.go    # Tipos compartilhados
└── internal/
    ├── core/
    │   ├── analyzer.go   # Análise de espaço
    │   └── cleaner.go    # Limpeza
    ├── modules/
    │   ├── index.go      # Registro de módulos
    │   ├── module.go     # Interface base
    │   └── *.go          # Módulos de limpeza
    └── utils/
        ├── os.go         # Utils SO
        └── logger.go     # Logging
```

## Desenvolvimento - Versão Node.js

A versão TypeScript está disponível em `piunter-cli-npm/`:

```bash
cd piunter-cli-npm
npm install
npm run build
npm test
```

## Enviando Alterações

1. Crie um branch de feature:

```bash
git checkout -b feature/nova-funcionalidade
```

2. Faça suas alterações e commite:

```bash
git commit -m "feat: add new feature"
```

3. Push para seu fork:

```bash
git push origin feature/nova-funcionalidade
```

4. Abra um Pull Request

## Formato de Mensagens de Commit

- `feat:` Nova funcionalidade
- `fix:` Correção de bug
- `docs:` Documentação
- `refactor:` Refatoração
- `test:` Adicionar testes
- `chore:` Manutenção

## Perguntas?

Abra uma issue no GitHub para perguntas sobre contribuições.