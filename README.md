```text
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║ ██████╗ ██╗██╗   ██╗███╗   ██╗████████╗███████╗██████╗            ║
║ ██╔══██╗██║██║   ██║████╗  ██║╚══██╔══╝██╔════╝██╔══██╗           ║
║ ██████╔╝██║██║   ██║██╔██╗ ██║   ██║   █████╗  ██████╔╝           ║
║ ██╔═══╝ ██║██║   ██║██║╚██╗██║   ██║   ██╔══╝  ██╔══██╗           ║
║ ██║     ██║╚██████╔╝██║ ╚████║   ██║   ███████╗██║  ██║           ║
║ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝           ║
║                                                                   ║
║              Limpeza e Otimização para Linux                      ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
```

<p align="center">
  <img src="https://img.shields.io/badge/Node.js-18+-green.svg" alt="Node.js">
  <img src="https://img.shields.io/badge/TypeScript-5.3-blue.svg" alt="TypeScript">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux-purple.svg" alt="Platform">
  <img src="https://img.shields.io/badge/Tests-10-brightgreen.svg" alt="Tests">
</p>

<p align="center">
  <img src="https://img.shields.io/badge/APT-Debian-orange.svg" alt="APT">
  <img src="https://img.shields.io/badge/Pacman-Arch-green.svg" alt="Pacman">
  <img src="https://img.shields.io/badge/DNF-Fedora-blue.svg" alt="DNF">
</p>

> CLI profissional para limpeza e otimização de sistemas Linux. Desenvolvido com TypeScript, seguindo Clean Architecture e boas práticas de código.

<!-- DEMO GIF PLACEHOLDER - Adicione seu GIF aqui -->
<!-- ![Demo](docs/demo.gif) -->

## Recursos

- Detecção automática de distribuição Linux (Debian, Ubuntu, Arch, Fedora, etc.)
- Suporte a múltiplos gerenciadores de pacotes (APT, Pacman, DNF)
- Módulos de limpeza:
  - Cache do usuário (~/.cache)
  - NPM, Yarn, PNPM
  - Flatpak e Snap
  - Docker
  - Logs do sistema (journalctl)
  - Gerenciador de pacotes
  - Arquivos grandes e AppImages
  - Miniaturas do sistema
  - Arquivos recentes
- Sistema de plugins para extensibilidade
- Shell completion (bash/zsh)
- Modo interativo com seleção por checkbox
- Modo dry-run para simular limpeza
- Confirmação obrigatória para operações destrutivas
- Elevação automática de privilégios (sudo)
- Relatório detalhado de espaço liberado
- Config file personalizável (~/.piunter.json)
- Testes unitários com Vitest

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

### Arquivos grandes

```bash
# Detectar arquivos > 100MB
piunter --large-files

# Com threshold customizado (em MB)
piunter --large-files --threshold=500
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
| `--snap` | Limpar Snap |
| `--docker` | Limpar Docker |
| `--logs` | Limpar logs do sistema |
| `--packages` | Limpar gerenciador de pacotes |
| `--large-files` | Detectar arquivos grandes |
| `--appimage` | Limpar AppImages |
| `--thumbs` | Limpar miniaturas |
| `--recent` | Limpar arquivos recentes |
| `--threshold=MB` | Threshold para arquivos grandes |
| `--analyze` | Apenas analisar sem limpar |
| `--dry-run` | Simular limpeza |
| `--force` | Pular confirmação |
| `--interactive` | Modo interativo |

## Desenvolvimento

```bash
# Instalar dependências
npm install

# Executar em modo dev
npm run dev

# Build
npm run build

# Lint
npm run lint

# Formatar código
npm run format

# Rodar testes
npm test

# Rodar testes em watch mode
npm run test:watch
```

## Arquitetura

```
src/
├── cli.ts              # Interface CLI
├── core/               # Lógica principal
│   ├── analyzer.ts    # Análise de espaço
│   └── cleaner.ts     # Execução de limpeza
├── modules/           # Módulos de limpeza
│   ├── cache.ts       # Cache do usuário
│   ├── npm.ts         # NPM/Yarn/PNPM
│   ├── flatpak.ts     # Flatpak
│   ├── snap.ts        # Snap
│   ├── docker.ts      # Docker
│   ├── logs.ts        # Logs do sistema
│   ├── packages.ts    # Gerenciadores de pacotes
│   ├── disk.ts        # Arquivos grandes e uso de disco
│   ├── appimage.ts    # AppImages
│   ├── thumbs.ts      # Miniaturas
│   └── recent.ts      # Arquivos recentes
├── utils/              # Utilitários
│   ├── exec.ts         # Execução de comandos
│   ├── os.ts          # Informações do sistema
│   ├── logger.ts      # Logging
│   ├── config.ts      # Config file
│   ├── completion.ts  # Shell completion
│   ├── plugins.ts     # Plugin system
│   └── progress.ts     # Progress bars
└── types/              # TypeScript types
```

## Segurança

- Nunca executa operações destrutivas sem confirmação (exceto com `--force`)
- Modo dry-run disponível para testar antes de aplicar
- Detecção de comandos disponíveis antes de executar
- Elevação automática de privilégios para operações do sistema
- Tratamento robusto de erros com fallbacks

## Config File

Crie `~/.piunter.json` para personalizar configurações:

```json
{
  "version": "1.0.0",
  "defaults": {
    "dryRun": false,
    "force": false,
    "modules": ["packages", "cache", "npm"]
  },
  "thresholds": {
    "largeFilesMB": 100,
    "logDays": 30,
    "journalSizeMB": 500
  }
}
```

## Compatibilidade

- Debian/Ubuntu (APT)
- Arch Linux/Manjaro (Pacman)
- Fedora/RHEL (DNF)
- Pop!_OS
- Linux Mint
- E outras distribuições baseadas nestas

## Plugins

Crie plugins personalizados em `~/.piunter/plugins/`:

```typescript
export default {
  id: 'meu-plugin',
  name: 'Meu Plugin',
  description: 'Descrição do plugin',
  version: '1.0.0',

  isAvailable() {
    return true;
  },

  async analyze() {
    return { module: this.id, items: [], totalSize: 0 };
  },

  async clean(dryRun = false) {
    return {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: []
    };
  }
};
```

## Wiki

Documentação detalhada disponível na pasta [wiki/](wiki/):
- Home
- Installation
- Basic Usage
- Configuration
- Modules
- Plugins
- Troubleshooting

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md) para guidelines.

## Changelog

Veja [CHANGELOG.md](CHANGELOG.md) para histórico de mudanças.

## Licença

MIT - João Braga

## Autor

<a href="https://github.com/joaomjbraga">
  <img src="https://img.shields.io/badge/GitHub-joaomjbraga-blue?style=flat&logo=github" alt="GitHub">
</a>
<a href="https://www.npmjs.com/~bforgeio">
  <img src="https://img.shields.io/badge/npm-@bforgeio-red?style=flat&logo=npm" alt="npm">
</a>
