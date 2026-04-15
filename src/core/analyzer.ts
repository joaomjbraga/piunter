import chalk from 'chalk';
import type { AnalysisResult } from '../types/index.js';
import { getAvailableModules, getModuleByIds, type Module } from '../modules/index.js';
import { logger } from '../utils/logger.js';

export interface AnalyzerOptions {
  threshold?: number;
}

export class Analyzer {
  private modules: Module[];
  private threshold?: number;

  constructor(moduleIds?: string[], options?: AnalyzerOptions) {
    this.modules = moduleIds
      ? getModuleByIds(moduleIds)
      : getModuleByIds(getAvailableModules().map(m => m.id));
    this.threshold = options?.threshold;
  }

  async analyze(): Promise<AnalysisResult[]> {
    const availableModules = this.modules.filter(m => m.isAvailable());

    if (availableModules.length === 0) {
      return [];
    }

    const results = await Promise.all(
      availableModules.map(m =>
        m.analyze(this.threshold).catch(error => {
          logger.debug(`${m.name}: ${(error as Error).message}`);
          return null;
        })
      )
    );

    return results.filter((r): r is AnalysisResult => r !== null);
  }

  getSummary(results: AnalysisResult[]): {
    totalSize: number;
    totalItems: number;
    byModule: Record<string, { size: number; items: number }>;
  } {
    const byModule: Record<string, { size: number; items: number }> = {};
    let totalSize = 0;
    let totalItems = 0;

    for (const result of results) {
      byModule[result.module] = {
        size: result.totalSize,
        items: result.items.length,
      };
      totalSize += result.totalSize;
      totalItems += result.items.length;
    }

    return { totalSize, totalItems, byModule };
  }

  printAnalysis(results: AnalysisResult[]): void {
    const summary = this.getSummary(results);

    console.log(`  ${chalk.bold('Analise de espaco recuperavel')}`);
    console.log();

    results.forEach(result => {
      const size = logger.formatBytes(result.totalSize);
      const count = result.items.length > 0 ? `${result.items.length} itens` : '';

      if (result.items.length > 0) {
        console.log(
          `    ${chalk.dim('-')} ${result.module.padEnd(12)} ${chalk.cyan(size)} ${chalk.dim(`(${count})`)}`
        );
      } else {
        console.log(`    ${chalk.dim('-')} ${result.module.padEnd(12)} ${chalk.dim('0 B')}`);
      }
    });

    logger.space();
    console.log(`  ${chalk.dim('─'.repeat(Math.min(process.stdout.columns || 60, 40) - 4))}`);

    const totalSize = logger.formatBytes(summary.totalSize);

    console.log();
    console.log(`  ${chalk.bold('Total')}`);
    console.log(`    ${chalk.dim('-')} ${chalk.white('Espaco:')} ${chalk.green.bold(totalSize)}`);
    console.log(`    ${chalk.dim('-')} ${chalk.white('Itens:')} ${chalk.cyan(summary.totalItems)}`);
    console.log();
  }
}

export function createAnalyzer(moduleIds?: string[], options?: AnalyzerOptions): Analyzer {
  return new Analyzer(moduleIds, options);
}
