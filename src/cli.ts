#!/usr/bin/env node

import chalk from 'chalk';
import inquirer from 'inquirer';
import { createAnalyzer, createCleaner } from './core/index.js';
import { getAvailableModules } from './modules/index.js';
import type { CleanOptions, CliFlags } from './types/index.js';
import { logger } from './utils/logger.js';
import { getDistroInfo } from './utils/os.js';

const VERSION = '1.0.0';

function isRoot(): boolean {
  return process.getuid?.() === 0 || process.env.USER === 'root';
}

const MODULE_MAP: Record<string, string[]> = {
  'cache': ['cache'],
  'npm': ['npm'],
  'yarn': ['yarn'],
  'pnpm': ['pnpm'],
  'flatpak': ['flatpak'],
  'docker': ['docker'],
  'logs': ['logs'],
  'packages': ['packages'],
  'large-files': ['large-files'],
  'disk': ['disk-usage'],
};

function parseFlags(args: string[]): CliFlags {
  return {
    all: args.includes('--all') || args.includes('-a'),
    cache: args.includes('--cache'),
    npm: args.includes('--npm'),
    yarn: args.includes('--yarn'),
    pnpm: args.includes('--pnpm'),
    flatpak: args.includes('--flatpak'),
    docker: args.includes('--docker'),
    logs: args.includes('--logs'),
    packages: args.includes('--packages'),
    analyze: args.includes('--analyze'),
    dryRun: args.includes('--dry-run') || args.includes('-n'),
    force: args.includes('--force') || args.includes('-f'),
    interactive: args.includes('--interactive') || args.includes('-i'),
    largeFiles: args.includes('--large-files'),
    largeFilesThreshold: parseInt(args.find(a => a.startsWith('--threshold='))?.split('=')[1] || '100'),
  };
}

function getModulesFromFlags(flags: CliFlags): string[] {
  const modules: string[] = [];

  if (flags.all) {
    return getAvailableModules()
      .filter(m => m.available)
      .map(m => m.id);
  }

  if (flags.cache) modules.push('cache');
  if (flags.npm) modules.push('npm');
  if (flags.yarn) modules.push('yarn');
  if (flags.pnpm) modules.push('pnpm');
  if (flags.flatpak) modules.push('flatpak');
  if (flags.docker) modules.push('docker');
  if (flags.logs) modules.push('logs');
  if (flags.packages) modules.push('packages');
  if (flags.largeFiles) modules.push('large-files');

  return modules;
}

async function showBanner(): Promise<void> {
  console.log(chalk.cyan(`
╔═══════════════════════════════════════════════════════════════════╗
║                                                                   ║
║ ██████╗ ██╗██╗   ██╗███╗   ██╗████████╗███████╗██████╗          ║
║ ██╔══██╗██║██║   ██║████╗  ██║╚══██╔══╝██╔════╝██╔══██╗         ║
║ ██████╔╝██║██║   ██║██╔██╗ ██║   ██║   █████╗  ██████╔╝         ║
║ ██╔═══╝ ██║██║   ██║██║╚██╗██║   ██║   ██╔══╝  ██╔══██╗         ║
║ ██║     ██║╚██████╔╝██║ ╚████║   ██║   ███████╗██║  ██║         ║
║ ╚═╝     ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚══════╝╚═╝  ╚═╝         ║
║                                                                   ║
║              Limpeza e Otimização para Linux                     ║
║                                                                   ║
╚═══════════════════════════════════════════════════════════════════╝
`));
}

async function showSystemInfo(): Promise<void> {
  const distro = getDistroInfo();
  console.log(chalk.dim(`  Sistema: ${distro.name} (${distro.packageManager})`));
  console.log();
}

async function interactiveMode(): Promise<string[]> {
  const availableModules = getAvailableModules();
  const choices = availableModules.map(m => ({
    name: `${m.available ? '○' : '✗'} ${m.name} - ${m.description}${!m.available ? ' (indisponivel)' : ''}`,
    value: m.id,
    disabled: !m.available,
    checked: m.available && ['packages', 'cache', 'npm'].includes(m.id),
  }));

  const answers = await inquirer.prompt([
    {
      type: 'checkbox',
      name: 'modules',
      message: chalk.cyan('Selecione os modulos para limpar:'),
      choices,
      default: ['packages', 'cache', 'npm'],
    },
    {
      type: 'confirm',
      name: 'confirm',
      message: chalk.yellow('Continuar com a limpeza?'),
      default: false,
    },
  ]);

  if (!answers.confirm) {
    console.log(chalk.dim('Operacao cancelada.'));
    process.exit(0);
  }

  return answers.modules;
}

async function analyzeMode(moduleIds?: string[]): Promise<void> {
  const analyzer = createAnalyzer(moduleIds);
  const results = await analyzer.analyze();
  analyzer.printAnalysis(results);
}

async function cleanMode(moduleIds: string[], options: CleanOptions): Promise<void> {
  if (!options.force && !options.dryRun) {
    const answer = await inquirer.prompt([
      {
        type: 'confirm',
        name: 'proceed',
        message: chalk.red('Confirma que deseja limpar estes modulos? Esta acao pode ser irreversivel.'),
        default: false,
      },
    ]);

    if (!answer.proceed) {
      console.log(chalk.dim('Operacao cancelada pelo usuario.'));
      process.exit(0);
    }
  }

  const cleaner = createCleaner(moduleIds, options);
  const report = await cleaner.clean();
  cleaner.printReport(report);
}

export async function main(): Promise<void> {
  const args = process.argv.slice(2);
  const flags = parseFlags(args);

  if (args.includes('--help') || args.includes('-h') || args.includes('help')) {
    console.log(chalk.cyan(`
╔════════════════════════════════════════════════════════════╗
║                    piunter - Ajuda                       ║
╠════════════════════════════════════════════════════════════╣
║                                                            ║
║  Modo interativo:                                          ║
║    $ piunter                                              ║
║    $ piunter --interactive                                ║
║                                                            ║
║  Analise:                                                  ║
║    $ piunter --analyze                                    ║
║                                                            ║
║  Modulos de limpeza:                                       ║
║    --all         Limpar todos os modulos                   ║
║    --cache       Cache do usuario (~/.cache)              ║
║    --npm         Cache do NPM                             ║
║    --yarn        Cache do Yarn                            ║
║    --pnpm        Cache do PNPM                            ║
║    --flatpak     Flatpak                                  ║
║    --docker      Docker                                   ║
║    --logs        Logs do sistema                          ║
║    --packages    Gerenciador de pacotes                   ║
║    --large-files Arquivos grandes                         ║
║                                                            ║
║  Opcoes:                                                  ║
║    --dry-run     Simular sem executar                     ║
║    --force       Pular confirmacao                        ║
║                                                            ║
║  Exemplos:                                                 ║
║    $ piunter --all                                        ║
║    $ piunter --npm --cache --dry-run                      ║
║    $ piunter --analyze                                    ║
║                                                            ║
╚════════════════════════════════════════════════════════════╝
    `));
    return;
  }

  if (flags.analyze) {
    const modules = getModulesFromFlags(flags);
    await showBanner();
    await showSystemInfo();
    await analyzeMode(modules.length > 0 ? modules : undefined);
    return;
  }

  if (flags.interactive || args.length === 0) {
    await showBanner();
    await showSystemInfo();
    const selectedModules = await interactiveMode();

    if (selectedModules.length === 0) {
      console.log(chalk.yellow('Nenhum modulo selecionado.'));
      return;
    }

    await cleanMode(selectedModules, {
      dryRun: flags.dryRun,
      force: flags.force,
      modules: selectedModules,
    });
    return;
  }

  const selectedModules = getModulesFromFlags(flags);

  if (selectedModules.length === 0) {
    console.log(chalk.yellow('Nenhum modulo especificado. Use --help para ver as opcoes.'));
    process.exit(1);
  }

  await showBanner();
  await showSystemInfo();

  if (!isRoot() && (flags.packages || flags.logs)) {
    logger.warn('Alguns modulos requerem privilegios sudo - o sistema solicitara sua senha quando necessario');
    console.log();
  }

  await cleanMode(selectedModules, {
    dryRun: flags.dryRun,
    force: flags.force,
    modules: selectedModules,
  });
}

main().catch((error) => {
  logger.error(`Erro fatal: ${error.message}`);
  console.error(error);
  process.exit(1);
});
