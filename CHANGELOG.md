# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.6.0] - 2026-05-30

### Added

- **Docker nuclear:** `docker stop` antes do prune; contagem de containers, imagens, volumes e redes na análise
- **Filtro `isLogGzFile()`:** Só considera `.gz` com padrão de log (terminam em `.log` ou dígito) em `/var/log`

### Changed

- **Docker:** `docker system prune -a --volumes -f` — remove volumes e build cache (antes só `prune -a`)
- **Docker:** Resultado de `docker stop` verificado; erros acumulados no módulo
- **`shouldSkipDir`:** `strings.Contains` substituído por `==` — evita falsos positivos (`node_modules_backup`, `my.cache.data`)
- **LogsModule:** Espaço limpo calculado por operação (journal vs gz) — corrige sobrestimativa quando uma falha
- **Version cache:** `saveVersionCache` chamado uma única vez com objeto completo — corrige overwrite de dados frescos
- **`packages.go`:** Exit code 1 do `apt-get` só é ignorado se stderr não contiver `"E:"` ou `"error"`
- **Executor pattern:** `LogsModule` e `SnapModule` usam `GetExecutor()` — agora testáveis via mock
- **Pacotes:** Tamanhos reais via `dpkg-query` (apt), `pacman -Si` (pacman), `dnf info` (dnf) em vez de estimativa fixa `10 MB × N`
- **Erros ignorados:** `json.Unmarshal`, `os.MkdirAll`, `GetDirSizeAsync` agora logam com `Debug`
- **Constante renomeada:** `VERSION` → `Version` (convenção Go)

### Removed

- **Sistema de config YAML:** `config.go` (131 linhas), `config_test.go` (141 linhas), dependência `gopkg.in/yaml.v3`
- **Código paralelo morto:** `analyzeParallel()` + `cleanParallel()` — 117 linhas; campo `parallel`; `GetOptimalWorkers()`
- **`CleaningError.ItemPath`:** campo nunca preenchido

### Fixed

- **LogsModule:** Espaço reportado agora reflete apenas as operações que realmente sucederam
- **Version cache:** Notificação de nova versão já não é silenciada por overwrite do cache
- **`shouldSkipDir`:** `node_modules_backup` já não corresponde a `node_modules`

## [1.5.0] - 2026-05-10

### Added

- **Verificador de versão:** Aviso automático quando uma nova versão do piunter é lançada
  - Consulta a GitHub Releases API
  - Cache de 24h em `~/.config/piunter/version_cache.json`
  - Silencioso em erro de rede

### Changed

- **SnapModule:** `Analyze` identifica revisões desactivadas (`disabled` na coluna Notes) e mede tamanho real do ficheiro `.snap`; `Clean` usa `snap remove --revision` em vez de `snap refresh --list` (que não limpava nada)
- **LogsModule:** Remove `/tmp` da varredura; mede tamanho real do journald (`journalctl --disk-usage`) e ficheiros `.gz` antigos (>30 dias) em vez de subdiretorias inteiras
- **FlatpakModule:** Mede tamanho real dos diretórios `/var/lib/flatpak/runtime` e `.removed` em vez de estimativa fixa de 50MB × nº de apps; `Clean` reporta espaço real em vez de 100MB fixo
- **Config parser:** `disabled_modules` e `exclude_paths` agora são lidos corretamente (estavam a ser escritos mas ignorados na leitura)
- **ExtractModule:** `getArchiveSize` usa `os.Stat` em vez de `xtractr.ExtractFile` (que extraía o arquivo como efeito secundário durante a análise)
- **CacheModule:** `icon-cache` já não é ignorado — é analisado e limpo como qualquer outro diretório de cache

### Removed

- **Módulos removidos:** NVM, SDKMAN, Mise, NPM, Yarn, PNPM, Extract e Compress — focar no propósito principal de limpeza
- **Dependência removida:** `golift.io/xtractr` e ~25 dependências transitivas
- **Código morto removido:** `ErrorHandler` (6 métodos), `Warn`, `List`, `ParseThreshold`, `HasPrefixCI`, `contains`, `stringsJoin`, `ConfigManager` (6 métodos), `CliFlags`, variável `allErrors` não utilizada

## [1.4.1] - 2026-04-27

### Fixed

- **Output simplificado:** Removidos erros de permissão inline no Analyze()
- **Thumbnails fix:** Correção para diretórios não vazios

### Refactored

- **Código morto:** Removido validator.go (177 linhas)
- **npm/nvm/sdkman:** Removidos errorHandlers não usados

## [1.4.0] - 2026-04-26

### Added

- **Script de instalação:** Script shell para instalação automática (`install/install.sh`) — removido em 1.5.1
- **GitHub Actions:** Workflow para build e release automático em `amd64` e `arm64`
- **Módulo NVM:** Suporte para limpar cache do Node Version Manager
- **Módulo SDKMAN:** Suporte para limpar cache do SDKMAN
- **Módulo Extract:** Extração de arquivos (removido em 1.4.2)
- **Módulo Compress:** Compressão de arquivos (removido em 1.4.2)

### Changed

- **Removido modo interativo:** CLI agora é puramente flag-based com help integrado
- **Performance otimizada:** Workers dinâmicos baseados em `runtime.NumCPU()`
- **Config cacheado:** Config carregado uma única vez por execução (`sync.Once`)
- **Large files otimizado:** Pula diretórios comuns (`.cache`, `node_modules`, `.git`, etc.)

### Fixed

- **Security:** Corrigido uso de `fmt.Errorf` com string não-constante (go vet)

## [1.3.0] - 2026-04-26

### Changed

- **Reescrita em Go:** O projeto foi completamente reescrito em Go para melhor performance e distribuição
- **Estrutura simplificada:** Projeto agora em um único diretório raiz

### Added

- **Novo módulo:** Suporte para esvaziar a lixeira (`--trash`)
- **Build em Go:** CLI leve e rápida com dependências mínimas
- **Binary standalone:** Não requer Node.js ou runtime adicional

### Migrando para Go

A versão Go é agora a recomendada. Para migrar:

```bash
# Clone o repositório
git clone https://github.com/joaomjbraga/piunter.git
cd piunter

# Build
go build -o piunter ./cmd

# Execute
./piunter --all
```

### Mantendo a Versão Node.js

A versão TypeScript continua disponível em `piunter-cli-npm/`:

```bash
cd piunter-cli-npm
npm install
npm run build
npx piunter --all
```

## [1.2.3] - 2026-04-15

### Security

- **Command Injection Fix:** Rewrote shell completion installation to use fs module directly instead of shell commands with string interpolation, eliminating potential command injection vulnerability
- **Type Safety:** Improved error handling in exec.ts to properly handle unknown error types

### Fixed

- **Critical Bug:** Fixed command injection vulnerability in bash/zsh completion installation (CVE equivalent fix)
- **logs.ts:** Fixed incorrect variable usage (`days` → `logDays`) in `cleanOldLogs()` that caused inconsistent behavior
- **cli.ts:** Fixed memory leak with signal listeners not being removed after prompt completion
- **cli.ts:** Fixed empty catch blocks violating ESLint no-empty rule
- **exec.ts:** Fixed type safety for error handling with unknown error types
- **os.ts:** Fixed empty catch block linting issue

### Changed

- **tsconfig.json:** Updated moduleResolution from "node" to "bundler" for better ESM compatibility with TypeScript 5+
- **cli.ts:** Refactored duplicate sudo check code into `checkSudoForModules()` function
- **completion.ts:** Shell completion now uses fs module directly for safer file operations
- **cli.ts:** Improved stdin operations with proper try-catch error handling and removed unused AbortController
- **Code Quality:** Removed all unused imports resulting in 0 lint warnings and 0 errors

### Added

- **completion.ts:** Added missing flags to shell completion (--appimage, --thumbs, --recent, --snap, --version, --list)