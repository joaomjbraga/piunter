# Configuração

O piunter pode ser configurado através do arquivo `~/.piunter.json`.

## Arquivo de Configuração

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
  },
  "sudo": {
    "autoPrompt": true
  }
}
```

## Configurações Disponíveis

### defaults

Configurações padrão para comandos.

### thresholds

Limites para detecção de arquivos.

### sudo

Configurações de elevação de privilégios.
