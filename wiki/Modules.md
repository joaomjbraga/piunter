# Módulos

O piunter possui módulos independentes para diferentes tipos de limpeza.

## Lista de Módulos

| Módulo | Descrição |
|--------|-----------|
| packages | Gerenciador de pacotes (APT, Pacman, DNF) |
| cache | Cache do usuário (~/.cache) |
| npm | Cache do NPM |
| yarn | Cache do Yarn |
| pnpm | Cache do PNPM |
| flatpak | Apps Flatpak |
| snap | Apps Snap |
| docker | Docker containers e imagens |
| logs | Logs do sistema |
| large-files | Arquivos grandes |
| appimage | Arquivos AppImage |
| thumbs | Miniaturas do sistema |
| recent | Arquivos recentes |

## Uso de Módulos

```bash
# Limpar módulo específico
piunter --npm

# Limpar múltiplos módulos
piunter --cache --npm --yarn

# Limpar todos disponíveis
piunter --all
```
