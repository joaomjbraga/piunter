```text
╔═════════════════════════════════════════════════════════════╗
║                                                             ║
║  ██████╗ ██╗     ██╗ ██████╗ ██████╗ ███████╗███████╗   ║
║ ██╔════╝ ██║     ██║██╔════╝██╔═══██╗██╔════╝██╔════╝   ║
║ ██║  ███╗██║     ██║██║     ██║   ██║█████╗  ███████╗   ║
║ ██║   ██║██║     ██║██║     ██║   ██║██╔══╝  ╚════██║   ║
║ ╚██████╔╝╚██████╗██║╚██████╗╚██████╔╝███████╗███████║   ║
║  ╚═════╝  ╚═════╝╚═╝ ╚═════╝ ╚═════╝ ╚══════╝╚══════╝   ║
║                                                             ║
║            Limpeza e Otimização para Linux                 ║
║                                                             ║
╚═════════════════════════════════════════════════════════════╝
```

<p align="center">
  <img src="https://img.shields.io/badge/Node.js-18+-green.svg" alt="Node.js">
  <img src="https://img.shields.io/badge/TypeScript-5.3-blue.svg" alt="TypeScript">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux-purple.svg" alt="Platform">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/APT-Debian-orange.svg" alt="APT">
  <img src="https://img.shields.io/badge/Pacman-Arch-green.svg" alt="Pacman">
  <img src="https://img.shields.io/badge/DNF-Fedora-blue.svg" alt="DNF">
</p>

## Recursos

- Detecção automática de distribuição Linux (Debian, Ubuntu, Arch, Fedora, etc.)
- Suporte a múltiplos gerenciadores de pacotes (APT, Pacman, DNF)
- Módulos de limpeza:
  - Cache do usuário
  - NPM, Yarn, PNPM
  - Flatpak
  - Docker
  - Logs do sistema
  - Gerenciador de pacotes
  - Arquivos grandes
- Modo interativo com seleção por checkbox
- Modo dry-run para simular limpeza
- Confirmação obrigatória para operações destrutivas
- Relatório detalhado de espaço liberado

## Instalação

### Via npm (global)

```bash
npm install -g @bforgeio/piunter
```

### Via npx

```bash
npx piunter
```

### Do código fonte

```bash
git clone https://github.com/joaomjbraga/piunter.git
cd piunter
npm install
npm run build
npm link
```

## Uso

### Modo interativo

```bash
piunter
# ou
piunter --interactive
```

### Análise

```bash
piunter --analyze
```

### Limpeza completa

```bash
piunter --all
```

### Limpeza seletiva

```bash
# Limpar cache npm
piunter --npm

# Limpar cache e logs
piunter --cache --logs

# Limpar Docker
piunter --docker

# Limpar gerenciador de pacotes
piunter --packages
```

### Dry-run (simulação)

```bash
piunter --all --dry-run
```

### Forçar execução

```bash
piunter --all --force
```

## Opções

| Flag | Descrição |
|------|-----------|
| `--all` | Selecionar todos os módulos disponíveis |
| `--cache` | Limpar cache do usuário |
| `--npm` | Limpar cache do NPM |
| `--yarn` | Limpar cache do Yarn |
| `--pnpm` | Limpar cache do PNPM |
| `--flatpak` | Limpar Flatpak |
| `--docker` | Limpar Docker |
| `--logs` | Limpar logs do sistema |
| `--packages` | Limpar gerenciador de pacotes |
| `--large-files` | Detectar arquivos grandes |
| `--analyze` | Apenas analisar sem limpar |
| `--dry-run` | Simular limpeza |
| `--force` | Pular confirmação |
| `--interactive` | Modo interativo |

## Arquitetura

```
src/
├── core/           # Lógica principal
│   ├── analyzer.ts # Análise de espaço
│   └── cleaner.ts  # Execução de limpeza
├── modules/        # Módulos de limpeza
│   ├── cache.ts    # Cache do usuário
│   ├── npm.ts      # NPM/Yarn/PNPM
│   ├── flatpak.ts  # Flatpak
│   ├── docker.ts   # Docker
│   ├── logs.ts     # Logs do sistema
│   └── packages.ts # Gerenciadores de pacotes
├── utils/          # Utilitários
│   ├── exec.ts     # Execução de comandos
│   ├── os.ts       # Informações do sistema
│   └── logger.ts   # Logging
└── cli.ts          # Interface CLI
```

## Segurança

- Nunca executa operações destrutivas sem confirmação (exceto com `--force`)
- Modo dry-run disponível para testar antes de aplicar
- Detecção de comandos disponíveis antes de executar
- Tratamento robusto de erros

## Compatibilidade

- Debian/Ubuntu (APT)
- Arch Linux/Manjaro (Pacman)
- Fedora/RHEL (DNF)
- Pop!_OS
- Linux Mint
- E outras distribuições baseadas nestas

## Licença

MIT

## Autor

João Braga - [GitHub](https://github.com/joaomjbraga)
