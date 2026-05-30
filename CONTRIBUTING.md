# Contributing to piunter

Obrigado pelo seu interesse em contribuir para o piunter!

## Começando

### Setup

```bash
# Clone o repositório
git clone https://github.com/joaomjbraga/piunter.git
cd piunter

# Instale dependências
go mod download

# Build
go build -o piunter ./cmd

# Execute
./piunter --help
```

### Requisitos

- Go 1.26+
- Linux (amd64 ou arm64)

### Verificações antes de commit

```bash
# Formatar código
go fmt ./...

# Verificar erros
go vet ./...

# Executar testes
go test ./...

# Verificar cobertura
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | tail -1
rm -f coverage.out

# Build final
go build -o piunter ./cmd
```

## Estrutura do Projeto

```
piunter/
├── .github/workflows/
│   └── release.yml          # Release automático (amd64 + arm64)
├── cmd/main.go              # Entry point + CLI (Cobra)
├── pkg/types/types.go        # Tipos compartilhados
└── internal/
    ├── core/
    │   ├── analyzer.go      # Análise de espaço
    │   └── cleaner.go       # Limpeza
    ├── modules/
    │   ├── index.go         # Registro de módulos
    │   ├── module.go        # Interface base
    │   └── *.go            # Módulos de limpeza
    └── utils/
        ├── os.go            # Utils SO
        ├── config.go        # Configuração + validação
        ├── executor.go      # Executor de comandos (mock para testes)
        ├── logger.go        # Formatação de output
        ├── update.go        # Verificador de versão
        ├── parallel.go      # Execução paralela
        └── errors.go        # Tipos de erro AnalysisError/CleaningError
```

## Adicionando um Novo Módulo

1. Crie o arquivo em `internal/modules/`:

```go
func (m *MyModule) IsAvailable() bool {
    // Verifica se o módulo pode ser usado
    return utils.IsCommandAvailable("comando-necessario")
}

func (m *MyModule) Analyze(threshold int) (*types.AnalysisResult, error) {
    // Analisa o que pode ser limpo
    return result, nil
}

func (m *MyModule) Clean(dryRun bool) (*types.CleaningResult, error) {
    // Executa a limpeza
    return result, nil
}
```

2. Registre no `internal/modules/index.go`.

3. Adicione a flag em `cmd/main.go` e no array `allModuleFlags`.

4. Adicione ao changelog e à tabela de módulos no README.

## Formato de Mensagens de Commit

Seguimos [Conventional Commits](https://www.conventionalcommits.org/):

| Tipo     | Descrição                        |
|----------| -------------------------------- |
| `feat`   | Nova funcionalidade              |
| `fix`    | Correção de bug                 |
| `docs`   | Documentação                    |
| `refactor` | Refatoração                   |
| `perf`   | Melhoria de performance         |
| `test`   | Adicionar/editar testes         |
| `chore`  | Manutenção, deps, CI/CD         |
| `ci`     | Configuração de CI/CD           |

Exemplos:
```bash
feat: add module for cleaning cargo cache
fix: correct error handling in docker module
docs: update installation instructions
perf: optimize parallel execution
```

## Enviando Alterações

1. Fork o repositório

2. Clone seu fork:
```bash
git clone https://github.com/SEU-USUARIO/piunter.git
cd piunter
```

3. Crie um branch:
```bash
git checkout -b feat/nova-funcionalidade
```

4. Faça suas alterações e commite:
```bash
git commit -m "feat: add new feature"
```

5. Push para seu fork:
```bash
git push origin feat/nova-funcionalidade
```

6. Abra um Pull Request no repositório original

## Diretrizes de Código

- Código deve passar em `go vet ./...` sem erros
- Código deve ser testado quando possível
- Preferir `fmt.Errorf("%s", err)` em vez de `fmt.Errorf(err)` (segurança)
- Usar `GetOptimalWorkers()` para workers paralelos
- Config deve ser cacheado com `sync.Once`

## Perguntas?

Abra uma issue no GitHub para dúvidas ou sugestões!