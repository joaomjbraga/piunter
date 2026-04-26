# piunter (v1.4.0)

<pre align="center">

██████╗ ██╗██╗   ██╗███╗   ██╗████████╗███████╗██████╗
██╔══██╗██║██║   ██║████╗  ██║╚══██╔══╝██╔════╝██╔══██╗
████╔╝██║██║   ██║██╔██╗ ██║   ██║   █████╗  ██████╔╝
██╔═══╝ ██║██║   ██║██║╚██╗██║   ██║   ██╔══╝  ██╔══██╗
██║     ██║╚██████╔╝██║ ╚████║   ██║   ███████╗██║  ██║
╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝

</pre>

CLI para limpeza e otimização de sistemas Linux.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-green.svg" alt="Go">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux-purple.svg" alt="Platform">
</p>

## Instalação

### Via binary release

```bash
# Baixe a versão mais recente
curl -L https://github.com/joaomjbraga/piunter/releases/latest/download/piunter-linux-amd64 -o piunter
chmod +x piunter
./piunter --all
```

### Via Go

```bash
go install github.com/joaomjbraga/piunter@latest
```

### Build local

```bash
git clone https://github.com/joaomjbraga/piunter.git
cd piunter/piunter-cli-go
go build -o piunter ./cmd
./piunter --help
```

## Uso

```bash
# Ver help
./piunter --help

# Lista módulos disponíveis
./piunter --list

# Limpar tudo
./piunter --all

# Limpar específicos
./piunter --npm --nvm --cache --trash

# Analisar sem limpar
./piunter --all --analyze

# Simular (dry-run)
./piunter --all --dry-run
```

## Módulos

| Módulo       | Flag            | Descrição                        |
| ------------ | --------------- | -------------------------------- |
| Pacotes      | `--packages`    | Remove pacotes órfãos           |
| NPM          | `--npm`         | Limpa cache do npm                |
| Yarn         | `--yarn`        | Limpa cache do Yarn              |
| PNPM         | `--pnpm`        | Limpa cache do pnpm             |
| NVM          | `--nvm`         | Limpa cache do NVM              |
| SDKMAN       | `--sdkman`      | Limpa cache do SDKMAN           |
| Cache        | `--cache`       | Limpa ~/.cache                  |
| Flatpak      | `--flatpak`     | Remove dados órfãos            |
| Snap         | `--snap`        | Remove revisões antigas          |
| Docker       | `--docker`      | Remove containers/imagens       |
| Logs         | `--logs`        | Limpa logs do sistema            |
| Large Files  | `--large-files` | Encontra arquivos grandes      |
| AppImage     | `--appimage`    | Remove AppImages               |
| Thumbs       | `--thumbs`      | Remove miniaturas              |
| Recent       | `--recent`      | Lista arquivos recentes        |
| Trash        | `--trash`       | Esvazia a lixeira             |

## Flags

| Flag             | Descrição                            |
| ---------------- | ------------------------------------ |
| `-a`, `--all`   | Executa todos os módulos             |
| `--analyze`     | Analisa sem limpar                  |
| `-n`, `--dry-run` | Simula execução                   |
| `-f`, `--force` | Pula confirmações                   |
| `--list`        | Lista módulos disponíveis           |
| `-h`, `--help`  | Mostra ajuda                        |
| `--threshold`   | Tamanho mínimo para arquivos grandes |

## Configuração

O piunter suporta um arquivo de configuração em `~/.config/piunter/config.yaml`:

```yaml
# Configuração do Piunter
version: 1.0
threshold_mb: 100
dry_run_default: false
debug_enabled: false
parallel: false

# Módulos desabilitados
disabled_modules:
  - npm

# Paths a excluir
exclude_paths:
  - /home/user/important

# Tamanhos estimados (MB)
orphan_package_mb: 10
flatpak_app_mb: 50
snap_revision_mb: 200
```

## Compatibilidade

- Debian/Ubuntu (APT)
- Arch/Manjaro (Pacman)
- Fedora/RHEL (DNF)

## Estrutura do Projeto

```
piunter-cli-go/
├── cmd/main.go           # Entry point + CLI
├── pkg/types/types.go    # Tipos compartilhados
└── internal/
    ├── core/
    │   ├── analyzer.go   # Análise de espaço
    │   └── cleaner.go    # Limpeza
    ├── modules/
    │   ├── index.go     # Registro de módulos
    │   ├── module.go     # Interface base
    │   ├── cache.go     # Cache usuário
    │   ├── npm.go       # NPM/Yarn/PNPM
    │   ├── nvm.go       # NVM
    │   ├── sdkman.go    # SDKMAN
    │   ├── packages.go  # Pacotes órfãos
    │   ├── docker.go    # Docker
    │   ├── system.go    # Logs/Flatpak/Snap
    │   ├── files.go     # Large files/AppImage/Thumbs/Recent
    │   ├── trash.go    # Lixeira
    │   ├── extract.go  # Extração de arquivos
    │   └── compress.go # Compressão de arquivos
    └── utils/
        ├── os.go        # Utils SO
        ├── logger.go   # Logging
        ├── config.go   # Configuração
        ├── errors.go   # Tratamento de erros
        ├── validator.go # Validação de paths
        ├── executor.go # Executor de comandos
        └── parallel.go # Execução paralela
```

## Segurança

- Nunca executa operações sem confirmação (exceto com `--force`)
- Dry-run disponível para testar antes
- Verifica comandos antes de executar
- Tratamento robusto de erros
- Validação de paths (proteção contra symlink attacks)
- Sistema de configuração persistente

## Testes

```bash
go test ./...
```

## Licença

MIT - João Braga

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md)