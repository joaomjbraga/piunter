# Uso Básico

## Modo Interativo

```bash
piunter
```

## Via Flags

### Limpeza completa
```bash
piunter --all
```

### Limpeza seletiva
```bash
piunter --npm --cache --yarn
```

### Análise
```bash
piunter --analyze
```

### Dry-run
```bash
piunter --all --dry-run
```

## Opções Úteis

| Flag | Descrição |
|------|-----------|
| `--all` | Todos os módulos |
| `--dry-run` | Simular sem executar |
| `--force` | Pular confirmação |
| `--interactive` | Modo interativo |
