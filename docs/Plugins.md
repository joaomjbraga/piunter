# Plugins

O piunter suporta um sistema de plugins para extensibilidade.

## Criando um Plugin

```typescript
// ~/.piunter/plugins/meu-plugin.js
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

## Instalando Plugins

Coloque seus plugins em `~/.piunter/plugins/`.

## Gerenciando Plugins

```bash
# Listar plugins
piunter --plugins list

# Habilitar plugin
piunter --plugins enable meu-plugin

# Desabilitar plugin
piunter --plugins disable meu-plugin
```
