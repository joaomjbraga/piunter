import { isCommandAvailable } from './exec.js';
import { logger } from './logger.js';
import { writeFileSync, mkdirSync } from 'fs';
import { join } from 'path';
import { getHomeDir } from './os.js';

const BASH_COMPLETION = `#!/bin/bash

_piunter_completion() {
    local cur prev opts
    COMPREPLY=()
    cur="\${COMP_WORDS[COMP_CWORD]}"
    prev="\${COMP_WORDS[COMP_CWORD-1]}"

    opts="--all --cache --npm --yarn --pnpm --flatpak --snap --docker --logs --packages --large-files --appimage --thumbs --recent --analyze --dry-run --force --interactive --help --threshold="

    case "\${prev}" in
        --threshold)
            COMPREPLY=($(compgen -W "10 50 100 500 1000" -- \${cur}))
            return 0
            ;;
        --interactive|--help|--analyze|--all|--dry-run|--force)
            return 0
            ;;
        *)
            ;;
    esac

    COMPREPLY=($(compgen -W "\${opts}" -- \${cur}))
    return 0
}

complete -F _piunter_completion piunter
`;

const ZSH_COMPLETION = `#compdef piunter

_piunter() {
    local -a opts
    opts=(
        '--all[Selecionar todos os módulos]'
        '--cache[Limpar cache do usuário]'
        '--npm[Limpar cache do NPM]'
        '--yarn[Limpar cache do Yarn]'
        '--pnpm[Limpar cache do PNPM]'
        '--flatpak[Limpar Flatpak]'
        '--snap[Limpar Snap]'
        '--docker[Limpar Docker]'
        '--logs[Limpar logs do sistema]'
        '--packages[Limpar gerenciador de pacotes]'
        '--large-files[Detectar arquivos grandes]'
        '--appimage[Limpar AppImages]'
        '--thumbs[Limpar miniaturas]'
        '--recent[Arquivos recentes]'
        '--analyze[Apenas analisar]'
        '--dry-run[Simular limpeza]'
        '--force[Pular confirmação]'
        '--interactive[Modo interativo]'
        '--help[Mostrar ajuda]'
        '--version[Versão]'
        '--list[Listar módulos]'
        '--threshold[Threshold para arquivos grandes]:threshold (10 50 100 500 1000)'
    )

    _arguments -s "\${opts[@]}"
}

_piunter "$@"
`;

export async function installBashCompletion(): Promise<boolean> {
  if (!isCommandAvailable('bash')) {
    logger.error('Bash não está disponível');
    return false;
  }

  const bashCompletionDir = '/etc/bash_completion.d';
  const completionPath = join(bashCompletionDir, 'piunter');

  try {
    writeFileSync(completionPath, BASH_COMPLETION, { mode: 0o755 });
    logger.success('Completion para Bash instalada em ' + bashCompletionDir);
    return true;
  } catch {
    const userCompletionDir = join(getHomeDir(), '.bash_completion.d');
    try {
      mkdirSync(userCompletionDir, { recursive: true });
      writeFileSync(join(userCompletionDir, 'piunter'), BASH_COMPLETION, { mode: 0o644 });
      logger.info('Completion para Bash instalada em ' + userCompletionDir);
      logger.info('Adicione ao seu ~/.bashrc: source ' + userCompletionDir + '/piunter');
      return true;
    } catch (err) {
      logger.error('Falha ao instalar completion: ' + (err as Error).message);
      return false;
    }
  }
}

export async function installZshCompletion(): Promise<boolean> {
  if (!isCommandAvailable('zsh')) {
    logger.error('Zsh não está disponível');
    return false;
  }

  const zshCompletionDir = join(getHomeDir(), '.zsh', 'completion');

  try {
    mkdirSync(zshCompletionDir, { recursive: true });
    writeFileSync(join(zshCompletionDir, '_piunter'), ZSH_COMPLETION, { mode: 0o644 });
    logger.success('Completion para Zsh instalada');
    logger.info('Adicione ao seu ~/.zshrc: fpath+=(' + zshCompletionDir + ')');
    return true;
  } catch (err) {
    logger.error('Falha ao instalar completion: ' + (err as Error).message);
    return false;
  }
}

export function showCompletionHelp(): void {
  console.log(`
Shell Completion
================

Para instalar completion no Bash:
  sudo piunter --completion bash > /etc/bash_completion.d/piunter

Ou para o usuário:
  piunter --completion bash >> ~/.bashrc

Para Zsh:
  piunter --completion zsh >> ~/.zshrc

Após instalar, reinicie o terminal ou execute:
  source ~/.bashrc  # ou
  source ~/.zshrc
`);
}
