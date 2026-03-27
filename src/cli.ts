#!/usr/bin/env node

import chalk from 'chalk';
import inquirer from 'inquirer';
import { createAnalyzer, createCleaner } from './core/index.js';
import { getAvailableModules } from './modules/index.js';
import type { CleanOptions, CliFlags } from './types/index.js';
import { validateThreshold } from './utils/config.js';
import { logger } from './utils/logger.js';
import { getDistroInfo } from './utils/os.js';
import { requestSudo, hasSudoPassword } from './utils/exec.js';

const VERSION = '1.2.1';

const MODULES_REQUIRING_SUDO = ['packages', 'logs', 'flatpak'];

async function promptYesNo(message: string): Promise<boolean> {
  return new Promise(resolve => {
    const cleanup = () => {
      process.stdin.setRawMode(false);
      process.stdin.pause();
      process.stdin.removeAllListeners('data');
      process.stdin.removeAllListeners('error');
    };

    const handleInterrupt = () => {
      cleanup();
      console.log(chalk.dim('\nOperacao cancelada.'));
      process.exit(0);
    };

    const ask = () => {
      process.stdout.write(`${message} `);
      process.stdin.setRawMode(true);
      process.stdin.resume();

      const handler = (chunk: Buffer) => {
        cleanup();
        const char = chunk.toString().toLowerCase();
        process.stdout.write(char + '\n');

        if (char === 'y' || char === 's') {
          resolve(true);
        } else if (char === 'n') {
          resolve(false);
        } else {
          ask();
        }
      };

      process.stdin.once('data', handler);
    };

    process.stdin.once('error', handleInterrupt);
    process.on('SIGINT', handleInterrupt);
    ask();
  });
}

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

function requiresSudo(moduleIds: string[]): boolean {
  return moduleIds.some(id => MODULES_REQUIRING_SUDO.includes(id));
}

export function getModulesFromFlags(flags: CliFlags): string[] {
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

export function parseFlags(args: string[]): CliFlags {
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
      (() => {
        const val = args.find(a => a.startsWith('--threshold='))?.split('=')[1];
        const parsed = val ? parseInt(val) : NaN;
        return isNaN(parsed) ? 100 : parsed;
      })(),
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
  console.log();
  console.log(chalk.cyan.bold(`  piunter`) + chalk.dim(' · CLI para Linux'));
  console.log(chalk.dim(`  ${line()}`));
  console.log();
}

function printHelp(): void {
  printHeader();

  console.log(`  ${chalk.bold('USO')}`);
  console.log(`    ${chalk.dim('piunter')} ${chalk.cyan('[opcoes]')}`);
  console.log();

  console.log(`  ${chalk.bold('OPCOES DE EXECUCAO')}`);
  console.log(`    ${chalk.cyan('--all')}            ${chalk.dim('Executa todos os modulos')}`);
  console.log(`    ${chalk.cyan('--analyze')}        ${chalk.dim('Analisa sem limpar')}`);
  console.log(`    ${chalk.cyan('--dry-run')}         ${chalk.dim('Simula a execucao')}`);
  console.log(`    ${chalk.cyan('--force')}          ${chalk.dim('Pula confirmacoes')}`);
  console.log(`    ${chalk.cyan('--interactive')}    ${chalk.dim('Modo interativo')}`);
  console.log();

  console.log(`  ${chalk.bold('MODULOS')}`);
  console.log(`    ${chalk.cyan('--cache')}         ${chalk.dim('Cache do usuario')}`);
  console.log(`    ${chalk.cyan('--npm')}            ${chalk.dim('Cache do NPM')}`);
  console.log(`    ${chalk.cyan('--yarn')}           ${chalk.dim('Cache do Yarn')}`);
  console.log(`    ${chalk.cyan('--pnpm')}           ${chalk.dim('Cache do PNPM')}`);
  console.log(`    ${chalk.cyan('--packages')}       ${chalk.dim('Pacotes orfaos')}`);
  console.log(`    ${chalk.cyan('--docker')}         ${chalk.dim('Containers e imagens')}`);
  console.log(`    ${chalk.cyan('--logs')}           ${chalk.dim('Logs do sistema')}`);
  console.log(`    ${chalk.cyan('--flatpak')}        ${chalk.dim('Dados orfaos do Flatpak')}`);
  console.log(`    ${chalk.cyan('--snap')}           ${chalk.dim('Revisoes antigas do Snap')}`);
  console.log(`    ${chalk.cyan('--large-files')}    ${chalk.dim('Arquivos grandes')}`);
  console.log(`    ${chalk.cyan('--appimage')}       ${chalk.dim('AppImages')}`);
  console.log(`    ${chalk.cyan('--thumbs')}         ${chalk.dim('Miniaturas em cache')}`);
  console.log(`    ${chalk.cyan('--recent')}         ${chalk.dim('Arquivos recentes')}`);
  console.log();

  console.log(`  ${chalk.bold('OUTRAS OPCOES')}`);
  console.log(`    ${chalk.cyan('--threshold=100')}   ${chalk.dim('Tamanho minimo (MB)')}`);
  console.log(`    ${chalk.cyan('--help')}            ${chalk.dim('Mostra esta ajuda')}`);
  console.log(`    ${chalk.cyan('--version')}         ${chalk.dim('Mostra a versao')}`);
  console.log(`    ${chalk.cyan('--list')}            ${chalk.dim('Lista modulos disponiveis')}`);
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
    console.log(chalk.yellow('Nenhum modulo selecionado.'));
    return [];
  }

  const confirm = await promptYesNo(chalk.yellow('Continuar com a limpeza? (y/s/N)'));

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
    const proceed = await promptYesNo(chalk.red.bold('Confirmar limpeza? (y/s/N)'));

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

    if (requiresSudo(selectedModules) && !isRoot() && !hasSudoPassword()) {
      console.log();
      console.log(chalk.yellow('  Alguns modulos requerem privilegios de administrador.'));
      const sudoOk = await requestSudo();
      if (!sudoOk) {
        console.log(chalk.dim('  Modulos que requerem sudo serao pulados.'));
      }
      console.log();
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

  if (requiresSudo(selectedModules) && !isRoot() && !hasSudoPassword()) {
    console.log();
    console.log(chalk.yellow('  Alguns modulos requerem privilegios de administrador.'));
    const sudoOk = await requestSudo();
    if (!sudoOk) {
      console.log(chalk.dim('  Modulos que requerem sudo serao pulados.'));
    }
    console.log();
  }

  await cleanMode(selectedModules, {
    dryRun: flags.dryRun,
    force: flags.force,
    modules: selectedModules,
  });
}

main().catch(error => {
  logger.error(`Erro: ${error.message || error}`);
  process.exit(1);
});
