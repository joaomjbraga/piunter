# piunter (v1.5.0)

<div>
  <img src=".github/preview.gif">
</div>

CLI para limpeza e otimização de sistemas Linux, escrita em Go.

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26+-green.svg" alt="Go">
  <img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License">
  <img src="https://img.shields.io/badge/Platform-Linux-purple.svg" alt="Platform">
  <img src="https://img.shields.io/badge/Architecture-amd64%20%7C%20arm64-cyan.svg" alt="Architecture">
</p>

## Instalação

### Binary release (recomendado)

```bash
# amd64
curl -L https://github.com/joaomjbraga/piunter/releases/latest/download/piunter-linux-amd64 -o piunter
chmod +x piunter
sudo mv piunter /usr/local/bin/
```

```bash
# arm64
curl -L https://github.com/joaomjbraga/piunter/releases/latest/download/piunter-linux-arm64 -o piunter
chmod +x piunter
sudo mv piunter /usr/local/bin/
```

### Via Go

```bash
go install github.com/joaomjbraga/piunter/cmd@latest
```

## Uso

```bash
# Ver versão
piunter --version

# Ver help
piunter --help

# Limpar tudo
piunter --all

# Limpar específicos
piunter --docker --cache --trash

# Analisar sem limpar (ver quanto pode recuperar)
piunter --all --analyze

# Simular execução (não remove nada)
piunter --all --dry-run

# Pular confirmações
piunter --all --force

# Limpar arquivos grandes (threshold customizado)
piunter --large-files --threshold=500
```

## Módulos

| Módulo           | Flag            | Descrição                        |
| ---------------- | --------------- | -------------------------------- |
| Cache            | `--cache`       | Limpa cache do usuário (~/.cache)|
| Pacotes          | `--packages`    | Remove pacotes órfãos (APT/Pacman/DNF) |
| Flatpak          | `--flatpak`     | Remove dados órfãos do Flatpak   |
| Snap             | `--snap`        | Remove revisões desativadas do Snap |
| Docker           | `--docker`      | Remove containers/imagens não utilizados |
| Logs             | `--logs`        | Limpa logs antigos do sistema (journald + .gz) |
| Large Files      | `--large-files` | Encontra arquivos grandes (> threshold) |
| AppImage         | `--appimage`    | Remove AppImages do diretório Downloads |
| Thumbs           | `--thumbs`      | Remove miniaturas em cache (~/.cache/thumbnails) |
| Recent           | `--recent`      | Lista arquivos modificados nos últimos 7 dias |
| Trash            | `--trash`       | Esvazia a lixeira do usuário     |

Flags dos módulos também podem ser combinadas com `--all` para execução completa.

## Flags

| Flag              | Descrição                            |
| ----------------- | ------------------------------------ |
| `-a`, `--all`     | Executa todos os módulos             |
| `--analyze`       | Analisa sem limpar (preview)         |
| `-n`, `--dry-run` | Simula execução                      |
| `-f`, `--force`   | Pula todas as confirmações           |
| `--list`          | Lista módulos disponíveis            |
| `--version`       | Mostra a versão do piunter           |
| `--threshold=MB`  | Tamanho mínimo para arquivos grandes (default: 100) |
| `-h`, `--help`    | Mostra ajuda                         |

**Auto-update**: O piunter verifica automaticamente no GitHub se há uma nova versão (cache de 24h em `~/.config/piunter/version_cache.json`). A notificação aparece no cabeçalho ao executar o comando.

> Para desativar a verificação: `export PIUNTER_SKIP_UPDATE_CHECK=1`

## Configuração

O piunter lê configurações de `~/.config/piunter/config.yaml`:

```yaml
version: 1.0
threshold_mb: 100
dry_run_default: false
debug_enabled: false
parallel: false

disabled_modules:
  - flatpak
  - snap

exclude_paths:
  - /home/user/documents

package_sizes:
  orphan_package_mb: 10
  flatpak_app_mb: 50
  snap_revision_mb: 200
```

## Compatibilidade

| Distribuição  | Gerenciador |
| ------------- | ----------- |
| Debian/Ubuntu | APT         |
| Arch/Manjaro  | Pacman      |
| Fedora/RHEL   | DNF         |

### Requisitos

- Linux (amd64 ou arm64)
- curl (para instalação via script)

### Ferramentas opcionais (por módulo)

- `flatpak` - para módulo flatpak
- `snap` - para módulo snap
- `docker` - para módulo docker

## Segurança

- Confirmação antes de limpar (exceto com `--force`)
- Dry-run disponível para simulação
- Execução paralela opcional (configurável em `config.yaml`)
- Módulos que requerem sudo solicitam elevação de privilégio automaticamente

## Desenvolvimento

```bash
# Build
go build -o piunter ./cmd

# Testes
go test ./...

# Vet
go vet ./...
```

## Licença

MIT - João Braga

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md)
