import chalk from 'chalk';
import type { CleaningResult, Report, CleanOptions } from '../types/index.js';
import { getModuleByIds, type Module } from '../modules/index.js';
import { logger } from '../utils/logger.js';

export class Cleaner {
  private modules: Module[];
  private dryRun: boolean;
  private force: boolean;

  constructor(moduleIds: string[], options: CleanOptions) {
    this.modules = getModuleByIds(moduleIds);
    this.dryRun = options.dryRun;
    this.force = options.force;
  }

  async clean(): Promise<Report> {
    const startTime = new Date();
    const results: CleaningResult[] = [];
    const errors: string[] = [];

    for (const module of this.modules) {
      if (!module.isAvailable()) {
        logger.debug(`${module.name} nao disponivel`);
        continue;
      }

      try {
        const result = await module.clean(this.dryRun, this.force);
        results.push(result);

        if (result.errors.length > 0) {
          errors.push(...result.errors);
        }
      } catch (error) {
        errors.push(`${module.name}: ${(error as Error).message}`);
        results.push({
          module: module.id,
          success: false,
          spaceFreed: 0,
          itemsRemoved: 0,
          errors: [(error as Error).message],
        });
      }
    }

    const endTime = new Date();

    return {
      startTime,
      endTime,
      modules: results,
      totalSpaceFreed: results.reduce((sum, r) => sum + r.spaceFreed, 0),
      totalItemsRemoved: results.reduce((sum, r) => sum + r.itemsRemoved, 0),
      errors,
    };
  }

  printReport(report: Report): void {
    const duration = report.endTime.getTime() - report.startTime.getTime();
    const durationStr = duration < 60000 
      ? `${(duration / 1000).toFixed(1)}s`
      : `${(duration / 60000).toFixed(1)}min`;

    const items = report.modules.map(r => ({
      name: r.module,
      value: r.success ? logger.formatBytes(r.spaceFreed) : 'erro',
      success: r.success,
    }));

    if (items.length > 0) {
      logger.list(items);
    }

    logger.space();
    console.log(`  ${chalk.dim('─'.repeat(Math.min(process.stdout.columns || 60, 40) - 4))}`);

    const totalSize = logger.formatBytes(report.totalSpaceFreed);
    const totalItems = report.totalItemsRemoved.toString();
    const totalErrors = report.errors.length.toString();

    console.log();
    console.log(`  ${chalk.bold('Resumo')}`);
    console.log(`    ${chalk.dim('-')} ${chalk.white('Espaco liberado:')} ${chalk.green(totalSize)}`);
    console.log(`    ${chalk.dim('-')} ${chalk.white('Itens removidos:')} ${chalk.cyan(totalItems)}`);
    console.log(`    ${chalk.dim('-')} ${chalk.white('Erros:')} ${totalErrors === '0' ? chalk.green(totalErrors) : chalk.red(totalErrors)}`);

    if (report.errors.length > 0) {
      logger.space();
      console.log(`  ${chalk.bold.red('Erros:')}`);
      report.errors.forEach((error) => {
        console.log(`    ${chalk.dim('-')} ${chalk.red(error)}`);
      });
    }

    logger.space();

    if (this.dryRun) {
      console.log(`  ${chalk.yellow('!')} Dry-run concluido`);
      console.log(chalk.dim(`    Execute sem --dry-run para aplicar`));
    } else {
      console.log(`  ${chalk.green('*')} Limpeza concluida em ${durationStr}`);
    }
    console.log();
  }
}

export function createCleaner(moduleIds: string[], options: CleanOptions): Cleaner {
  return new Cleaner(moduleIds, options);
}
