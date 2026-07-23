# piunter (v1.6.0)

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

### Opção 1: binário (recomendada)

Use esta opção se você quiser instalar rapidamente sem depender do ambiente Go.

```bash
# amd64
curl -L https://github.com/joaomjbraga/piunter/releases/latest/download/piunter-linux-amd64 -o /tmp/piunter
chmod +x /tmp/piunter
sudo mv /tmp/piunter /usr/local/bin/piunter
```

```bash
# arm64
curl -L https://github.com/joaomjbraga/piunter/releases/latest/download/piunter-linux-arm64 -o /tmp/piunter
chmod +x /tmp/piunter
sudo mv /tmp/piunter /usr/local/bin/piunter
```

Verifique a instalação:

```bash
piunter --version
```

### Opção 2: via Go

Use esta opção se você já tiver Go instalado e preferir compilar localmente.

```bash
go install github.com/joaomjbraga/piunter/cmd@latest
```

Se o binário não estiver no PATH, adicione o diretório do Go:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

### Opção 3: build local

Para testar a versão atual do repositório localmente:

```bash
go build -o /tmp/piunter ./cmd
sudo install /tmp/piunter /usr/local/bin/piunter
```

### Teste rápido após a instalação

```bash
piunter --help
piunter --list
```

## Uso

### Autocomplete para shell

O CLI agora suporta autocompletar para Bash e Zsh.

Bash:

```bash
piunter completion bash > /etc/bash_completion.d/piunter
source /etc/bash_completion.d/piunter
```

Zsh:

```bash
piunter completion zsh > "${fpath[1]}/_piunter"
chmod +x "${fpath[1]}/_piunter"
source ~/.zshrc
```

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

# Mostrar detalhes adicionais da execução
piunter --all --verbose

# Entrar no modo interativo para escolher módulos
piunter --interactive

# Limpar arquivos grandes (threshold customizado)
piunter --large-files --threshold=500

# Limpar arquivos antigos em Downloads
piunter --downloads-old

# Limpar caches de navegadores
piunter --browser-cache

# Limpar caches de editores
piunter --editor-cache

# Limpar caches de mídia e aplicativos gráficos
piunter --media-cache

# Limpar caches de jogos e gaming
piunter --game-cache

# Limpar caches de containers e VMs
piunter --container-cache

# Limpar caches de ferramentas de build e IA
piunter --build-cache

# Limpar temporários de IDEs e resíduos de desinstalação
piunter --ides-cache

# Limpar resíduos de plugins de navegadores
piunter --browser-plugins

# Remover instaladores antigos
piunter --old-installers

# Remover arquivos swap e temporários de editores
piunter --swap-files

# Limpar logs de aplicativos e shells
piunter --app-logs
```

## Módulos

| Módulo      | Flag            | Descrição                                                                          |
| ----------- | --------------- | ---------------------------------------------------------------------------------- |
| Cache       | `--cache`       | Limpa cache do usuário (~/.cache)                                                  |
| Pacotes     | `--packages`    | Remove pacotes órfãos (APT/Pacman/DNF)                                             |
| Flatpak     | `--flatpak`     | Remove dados órfãos do Flatpak                                                     |
| Snap        | `--snap`        | Remove revisões desativadas do Snap                                                |
| Docker      | `--docker`      | Remove todos os recursos Docker (containers, imagens, volumes, redes, build cache) |
| Logs        | `--logs`        | Limpa logs antigos do sistema (journald + .gz)                                     |
| Large Files | `--large-files` | Encontra arquivos grandes (> threshold)                                            |
| AppImage    | `--appimage`    | Remove AppImages do diretório Downloads                                            |
| Downloads   | `--downloads-old` | Limpa arquivos antigos em Downloads e pastas semelhantes                        |
| Browser     | `--browser-cache` | Limpa caches de navegadores como Firefox e Chrome                                |
| Editor      | `--editor-cache`  | Limpa caches de editores como VS Code, IntelliJ e Vim                             |
| Media       | `--media-cache`   | Limpa caches de mídia e aplicativos gráficos                                      |
| Games       | `--game-cache`    | Limpa caches de jogos e ferramentas de gaming                                     |
| Containers  | `--container-cache` | Limpa caches e artefatos de containers e VMs                                     |
| Build       | `--build-cache`     | Limpa caches de ferramentas de build e IA                                        |
| IDEs        | `--ides-cache`      | Limpa temporários de IDEs e resíduos de desinstalação                            |
| Browser     | `--browser-plugins` | Limpa resíduos de extensões e plugins de navegadores                             |
| Installers  | `--old-installers`  | Remove instaladores antigos em Downloads e pastas semelhantes                    |
| Swap        | `--swap-files`      | Remove arquivos swap e temporários de editores                                   |
| Logs        | `--app-logs`        | Limpa logs de aplicativos e shells                                               |
| Thumbs      | `--thumbs`      | Remove miniaturas em cache (~/.cache/thumbnails)                                   |
| Recent      | `--recent`      | Lista arquivos modificados nos últimos 7 dias                                      |
| Trash       | `--trash`       | Esvazia a lixeira do usuário                                                       |

Flags dos módulos também podem ser combinadas com `--all` para execução completa.

## Flags

| Flag              | Descrição                                           |
| ----------------- | --------------------------------------------------- |
| `-a`, `--all`     | Executa todos os módulos                            |
| `--analyze`       | Analisa sem limpar (preview)                        |
| `-n`, `--dry-run` | Simula execução                                     |
| `-f`, `--force`   | Pula todas as confirmações                          |
| `--verbose`       | Mostra detalhes adicionais da execução              |
| `--interactive`   | Ativa o modo interativo para escolher módulos       |
| `--list`          | Lista módulos disponíveis                           |
| `--version`       | Mostra a versão do piunter                          |
| `--threshold=MB`  | Tamanho mínimo para arquivos grandes (default: 100) |
| `-h`, `--help`    | Mostra ajuda                                        |

**Auto-update**: O piunter verifica automaticamente no GitHub se há uma nova versão (cache de 24h em `~/.config/piunter/version_cache.json`). A notificação aparece no cabeçalho ao executar o comando.

> Para desativar a verificação: `export PIUNTER_SKIP_UPDATE_CHECK=1`

## Configuração

O piunter também aceita configuração via arquivo JSON ou variáveis de ambiente.

### Arquivo de configuração

Por padrão, ele procura em `~/.config/piunter/piunter.json`.

Exemplo:

```json
{
  "threshold_mb": 150,
  "dry_run": true,
  "force": false,
  "verbose": false,
  "modules": ["cache", "trash"],
  "skip_update_check": false
}
```

Você também pode apontar outro arquivo:

```bash
PIUNTER_CONFIG_FILE=/caminho/para/piunter.json piunter --all
```

### Variáveis de ambiente

```bash
PIUNTER_THRESHOLD=200
PIUNTER_DRY_RUN=1
PIUNTER_FORCE=1
PIUNTER_VERBOSE=1
PIUNTER_MODULES=cache,trash
PIUNTER_SKIP_UPDATE_CHECK=1
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

- Confirmação detalhada antes de limpar (exceto com `--force`)
- Dry-run disponível para simulação
- Modo `--verbose` para exibir detalhes da execução
- Execução sequencial por padrão
- Módulos que requerem sudo solicitam elevação de privilégio automaticamente
- Erros de permissão são tratados de forma mais resiliente durante a limpeza

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

MIT - João M J Braga

## Contribuindo

[![João M J Braga](https://github.com/joaomjbraga.png?size=100)](https://github.com/joaomjbraga)

Veja [CONTRIBUTING.md](CONTRIBUTING.md)
