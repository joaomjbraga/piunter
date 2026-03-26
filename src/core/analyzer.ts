import chalk from 'chalk';
import type { AnalysisResult } from '../types/index.js';
import { getAvailableModules, getModuleByIds, type Module } from '../modules/index.js';
import { logger } from '../utils/logger.js';

export class Analyzer {
  private modules: Module[];

  constructor(moduleIds?: string[]) {
    this.modules = moduleIds ? getModuleByIds(moduleIds) : getModuleByIds(getAvailableModules().map(m => m.id));
  }

  async analyze(): Promise<AnalysisResult[]> {
    const results: AnalysisResult[] = [];

    for (const module of this.modules) {
      if (!module.isAvailable()) {
        logger.debug(`${module.name} nao disponivel`);
        continue;
      }

      try {
        const result = await module.analyze();
        results.push(result);
      } catch (error) {
        logger.debug(`Erro: ${(error as Error).message}`);
      }
    }

    return results;
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
    
    results.forEach((result) => {
      const size = logger.formatBytes(result.totalSize);
      const count = result.items.length > 0 ? `${result.items.length} itens` : '';
      
      if (result.items.length > 0) {
        console.log(`    ${chalk.dim('-')} ${result.module.padEnd(12)} ${chalk.cyan(size)} ${chalk.dim(`(${count})`)}`);
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

export function createAnalyzer(moduleIds?: string[]): Analyzer {
  return new Analyzer(moduleIds);
}
