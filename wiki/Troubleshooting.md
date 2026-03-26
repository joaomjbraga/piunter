# Troubleshooting

## Problemas Comuns

### Permissão negada

 بعض الأوامر تتطلب صلاحيات sudo. سيطلب النظام كلمة المرور تلقائياً.

### Docker não encontrado

Ensure Docker is installed and running:
```bash
docker --version
sudo systemctl status docker
```

### Flatpak não responde

```bash
flatpak repair --system
```

## Logs

Para debug, execute com `DEBUG=1`:
```bash
DEBUG=1 piunter --all --dry-run
```

## Obter Ajuda

```bash
piunter --help
```
