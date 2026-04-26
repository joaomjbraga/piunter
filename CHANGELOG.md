# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.3.0] - 2026-04-26

### Changed

- **Reescrita em Go:** O projeto foi completamente reescrito em Go para melhor performance e distribuição
- **Nova estrutura:** Projeto agora possui duas versões:
  - `piunter-cli-npm/` - Versão original em TypeScript/Node.js
  - `piunter-cli-go/` - Nova versão em Go (recomendada)

### Added

- **Novo módulo:** Suporte para esvaziar a lixeira (`--trash`)
- **Build em Go:** CLI leve e rápida com dependências mínimas
- **Binary standalone:** Não requer Node.js ou runtime adicional

### Migrando para Go

A versão Go é agora a recomendada. Para migrar:

```bash
# Clone o repositório
git clone https://github.com/joaomjbraga/piunter.git
cd piunter/piunter-cli-go

# Build
go build -o piunter ./cmd/main.go

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