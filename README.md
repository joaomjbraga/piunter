# piunter

CLI para limpeza e otimização de sistemas Linux.

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

## Instalação

```bash
npm install -g @bforgeio/piunter
```

## Uso Rápido

```bash
# Modo interativo
piunter

# Limpar tudo
piunter --all

# Limpar específicos
piunter --npm --cache --logs

# Simular (dry-run)
piunter --all --dry-run
```

## Recursos

- Detecção automática de distribuição (Debian, Ubuntu, Arch, Fedora)
- Suporte a APT, Pacman e DNF
- 14 módulos de limpeza
- Modo interativo com seleção por checkbox
- Dry-run para simular antes de executar
- Confirmação obrigatória para operações destrutivas
- Sistema de plugins
- Shell completion (bash/zsh)
- Config file personalizável

## Módulos

| Módulo | Flag | Descrição |
|--------|------|-----------|
| Pacotes | `--packages` | Remove pacotes órfãos |
| NPM | `--npm` | Limpa cache do npm |
| Yarn | `--yarn` | Limpa cache do Yarn |
| PNPM | `--pnpm` | Limpa cache do pnpm |
| Cache | `--cache` | Limpa ~/.cache |
| Flatpak | `--flatpak` | Remove dados órfãos |
| Snap | `--snap` | Remove revisões antigas |
| Docker | `--docker` | Remove containers/imagens |
| Logs | `--logs` | Limpa logs do sistema |
| Large Files | `--large-files` | Encontra arquivos grandes |
| AppImage | `--appimage` | Gerencia AppImages |
| Thumbs | `--thumbs` | Remove miniaturas |
| Recent | `--recent` | Limpa arquivos recentes |
| Analyze | `--analyze` | Analisa uso de disco |

## Flags

| Flag | Descrição |
|------|-----------|
| `--all` | Executa todos os módulos |
| `--analyze` | Analisa sem limpar |
| `--dry-run` | Simula execução |
| `--force` | Pula confirmações |
| `--threshold=MB` | Tamanho mínimo para arquivos grandes |
| `--config` | Arquivo de configuração customizado |

## Documentação

Documentação completa disponível em: **[https://joaomjbraga.github.io/piunter/docs.html](https://joaomjbraga.github.io/piunter/docs.html)**

## Desenvolvimento

```bash
# Instalar dependências
npm install

# Build
npm run build

# Testes
npm test

# Lint
npm run lint
```

## Segurança

- Nunca executa operações sem confirmação (exceto com `--force`)
- Dry-run disponível para testar antes
- Verifica comandos antes de executar
- Tratamento robusto de erros
- Cálculo preciso de espaço liberado (antes/depois)
- Timeout de 30s para prompts de senha sudo
- Sem command injection (pacotes processados de forma segura)

## Config File

Crie `~/.piunter.json`:

```json
{
  "threshold": 100,
  "modules": {
    "npm": true,
    "cache": true,
    "logs": true
  }
}
```

## Compatibilidade

- Debian/Ubuntu (APT)
- Arch/Manjaro (Pacman)
- Fedora/RHEL (DNF)

## Licença

MIT - João Braga

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md)
