#!/usr/bin/env node

import chalk from 'chalk';
import inquirer from 'inquirer';
import { createAnalyzer, createCleaner } from './core/index.js';
import { getAvailableModules } from './modules/index.js';
import type { CleanOptions, CliFlags } from './types/index.js';
import { validateThreshold } from './utils/config.js';
import { logger } from './utils/logger.js';
import { getDistroInfo } from './utils/os.js';

const VERSION = '1.0.0';

function getTerminalWidth(): number {
  return process.stdout.columns || 80;
}

function isRoot(): boolean {
  return process.getuid?.() === 0 || process.env.USER === 'root';
}

function padEnd(str: string, len: number): string {
  const width = getTerminalWidth();
  const maxLen = Math.min(len, width - 10);
  return str.length >= maxLen ? str.substring(0, maxLen - 3) + '...' : str.padEnd(maxLen);
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
  if (flags.snap) modules.push('snap');
  if (flags.docker) modules.push('docker');
  if (flags.logs) modules.push('logs');
  if (flags.packages) modules.push('packages');
  if (flags.largeFiles) modules.push('large-files');
  if (flags.appimage) modules.push('appimage');
  if (flags.thumbs) modules.push('thumbs');
  if (flags.recent) modules.push('recent');

  return modules;
}

function parseFlags(args: string[]): CliFlags {
  return {
    all: args.includes('--all') || args.includes('-a'),
    cache: args.includes('--cache'),
    npm: args.includes('--npm'),
    yarn: args.includes('--yarn'),
    pnpm: args.includes('--pnpm'),
    flatpak: args.includes('--flatpak'),
    snap: args.includes('--snap'),
    docker: args.includes('--docker'),
    logs: args.includes('--logs'),
    packages: args.includes('--packages'),
    analyze: args.includes('--analyze'),
    dryRun: args.includes('--dry-run') || args.includes('-n'),
    force: args.includes('--force') || args.includes('-f'),
    interactive: args.includes('--interactive') || args.includes('-i'),
    largeFiles: args.includes('--large-files'),
    largeFilesThreshold: validateThreshold(
      parseInt(args.find(a => a.startsWith('--threshold='))?.split('=')[1] || '100'),
      1,
      10000
    ),
    appimage: args.includes('--appimage'),
    thumbs: args.includes('--thumbs'),
    recent: args.includes('--recent'),
  };
}

function line(char: string = '─', len?: number): string {
  const width = len || getTerminalWidth() - 4;
  return char.repeat(Math.max(width, 10));
}

function printHeader(): void {
  const w = getTerminalWidth();
  console.log();
  console.log(chalk.cyan.bold(`  piunter`) + chalk.dim(' · CLI para Linux'));
  console.log(chalk.dim(`  ${line()}`));
  console.log();
}

function printHelp(): void {
  const w = getTerminalWidth();
  printHeader();
  
  console.log(`  ${chalk.bold('USO')}`);
  console.log(`    ${chalk.dim('piunter')} ${chalk.cyan('[flags]')}`);
  console.log();
  
  console.log(`  ${chalk.bold('FLAGS')}`);
  console.log(`    ${chalk.cyan('--all')}`);
  console.log(`    ${chalk.cyan('--analyze')}`);
  console.log(`    ${chalk.cyan('--dry-run')}  ${chalk.dim('- Simula execução')}`);
  console.log(`    ${chalk.cyan('--force')}`);
  console.log(`    ${chalk.cyan('--interactive')}`);
  console.log();
  
  console.log(`  ${chalk.bold('MÓDULOS')}`);
  console.log(`    ${chalk.cyan('--cache')}         Cache do usuário`);
  console.log(`    ${chalk.cyan('--npm')}            Cache do NPM`);
  console.log(`    ${chalk.cyan('--yarn')}           Cache do Yarn`);
  console.log(`    ${chalk.cyan('--pnpm')}           Cache do PNPM`);
  console.log(`    ${chalk.cyan('--packages')}       Pacotes órfãos`);
  console.log(`    ${chalk.cyan('--docker')}         Containers`);
  console.log(`    ${chalk.cyan('--logs')}           Logs`);
  console.log(`    ${chalk.cyan('--flatpak')}        Flatpak`);
  console.log(`    ${chalk.cyan('--snap')}           Snap`);
  console.log(`    ${chalk.cyan('--large-files')}    Arquivos grandes`);
  console.log(`    ${chalk.cyan('--appimage')}       AppImages`);
  console.log(`    ${chalk.cyan('--thumbs')}         Miniaturas`);
  console.log(`    ${chalk.cyan('--recent')}         Arquivos recentes`);
  console.log();
  
  console.log(`  ${chalk.bold('OPÇÕES')}`);
  console.log(`    ${chalk.cyan('--threshold=100')}  ${chalk.dim('Tamanho mínimo em MB')}`);
  console.log();
  
  console.log(`  ${chalk.bold('EXEMPLOS')}`);
  console.log(`    ${chalk.dim('piunter --all')}`);
  console.log(`    ${chalk.dim('piunter --npm --cache')}`);
  console.log(`    ${chalk.dim('piunter --all --dry-run')}`);
  console.log();
}

async function interactiveMode(): Promise<string[]> {
  const availableModules = getAvailableModules();
  
  const choices = availableModules.map(m => ({
    name: `${m.name} - ${m.description}`,
    value: m.id,
    checked: m.available && ['packages', 'cache', 'npm'].includes(m.id),
    disabled: !m.available,
  }));

  const { modules } = await inquirer.prompt([
    {
      type: 'checkbox',
      name: 'modules',
      message: chalk.cyan('Selecione os modulos:'),
      choices,
      pageSize: 10,
    },
  ]);

  if (!modules || modules.length === 0) {
    return [];
  }

  const { confirm } = await inquirer.prompt([
    {
      type: 'confirm',
      name: 'confirm',
      message: chalk.yellow('Continuar com a limpeza?'),
      default: false,
    },
  ]);

  if (!confirm) {
    console.log(chalk.dim('Operacao cancelada.'));
    process.exit(0);
  }

  return modules;
}

async function analyzeMode(moduleIds?: string[]): Promise<void> {
  const analyzer = createAnalyzer(moduleIds);
  const results = await analyzer.analyze();
  analyzer.printAnalysis(results);
}

async function cleanMode(moduleIds: string[], options: CleanOptions): Promise<void> {
  if (!options.force && !options.dryRun) {
    console.log();
    const { proceed } = await inquirer.prompt([
      {
        type: 'confirm',
        name: 'proceed',
        message: chalk.red.bold('Confirmar limpeza?'),
        default: false,
      },
    ]);

    if (!proceed) {
      console.log(chalk.dim('Operacao cancelada.'));
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
    printHelp();
    return;
  }

  if (args.includes('--version') || args.includes('-V')) {
    console.log(chalk.cyan(`piunter v${VERSION}`));
    return;
  }

  if (args.includes('--list')) {
    const modules = getAvailableModules();
    console.log();
    console.log(chalk.bold('  Modulos disponiveis:'));
    console.log();
    for (const m of modules) {
      const status = m.available ? chalk.green('*') : chalk.red('-');
      const name = padEnd(m.name, 12);
      console.log(`  ${status} ${chalk.white(name)} ${chalk.dim(m.description)}`);
    }
    console.log();
    return;
  }

  const distro = getDistroInfo();
  
  if (flags.analyze) {
    const modules = getModulesFromFlags(flags);
    printHeader();
    console.log(chalk.dim(`  Sistema: ${distro.name}`));
    console.log(chalk.dim(`  Gerenciador: ${distro.packageManager}`));
    console.log();
    await analyzeMode(modules.length > 0 ? modules : undefined);
    return;
  }

  if (flags.interactive || args.length === 0) {
    printHeader();
    console.log(chalk.dim(`  Sistema: ${distro.name}`));
    console.log(chalk.dim(`  Gerenciador: ${distro.packageManager}`));
    console.log();
    
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
    console.log(chalk.red('Nenhum modulo especificado.'));
    console.log(chalk.dim('Use ') + chalk.cyan('--help') + chalk.dim(' para ver as opcoes.'));
    console.log();
    process.exit(1);
  }

  printHeader();
  console.log(chalk.dim(`  Sistema: ${distro.name}`));
  console.log(chalk.dim(`  Gerenciador: ${distro.packageManager}`));
  console.log();

  if (flags.dryRun) {
    console.log(chalk.yellow('  Modo dry-run ativo'));
    console.log();
  }

  if (!isRoot() && (flags.packages || flags.logs)) {
    console.log(chalk.dim('  Alguns modulos requerem sudo'));
    console.log();
  }

  await cleanMode(selectedModules, {
    dryRun: flags.dryRun,
    force: flags.force,
    modules: selectedModules,
  });
}

main().catch((error) => {
  logger.error(`Erro: ${error.message}`);
  process.exit(1);
});
