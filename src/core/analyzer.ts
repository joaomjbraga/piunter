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
        logger.debug(`${module.name} não disponível`);
        continue;
      }

      logger.startSpinner(`Analisando ${module.name}...`);
      try {
        const result = await module.analyze();
        results.push(result);
        logger.stopSpinner(true, `${module.name}: ${logger.formatBytes(result.totalSize)}`);
      } catch (error) {
        logger.stopSpinner(false, `Erro ao analisar ${module.name}`);
        logger.debug(`Erro: ${(error as Error).message}`);
      }
    }

    return results;
  }

  async analyzeSingle(moduleId: string): Promise<AnalysisResult | null> {
    const module = getModuleByIds([moduleId])[0];
    
    if (!module) {
      logger.error(`Módulo não encontrado: ${moduleId}`);
      return null;
    }

    if (!module.isAvailable()) {
      logger.warn(`${module.name} não está disponível no sistema`);
      return null;
    }

    logger.startSpinner(`Analisando ${module.name}...`);
    try {
      const result = await module.analyze();
      logger.stopSpinner(true, `${module.name}: ${logger.formatBytes(result.totalSize)}`);
      return result;
    } catch (error) {
      logger.stopSpinner(false, `Erro ao analisar ${module.name}`);
      logger.debug(`Erro: ${(error as Error).message}`);
      return null;
    }
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
    logger.title('Análise de Espaço');

    for (const result of results) {
      if (result.items.length > 0) {
        logger.subtitle(`${result.module}: ${logger.formatBytes(result.totalSize)}`);
        
        const topItems = result.items.slice(0, 5);
        for (const item of topItems) {
          logger.item(item.description, logger.formatBytes(item.size));
        }

        if (result.items.length > 5) {
          logger.item(`... e mais ${result.items.length - 5} itens`);
        }
      }
    }

    const summary = this.getSummary(results);
    logger.space();
    logger.info(`Total recuperável: ${logger.formatBytes(summary.totalSize)}`);
    logger.info(`Total de itens: ${summary.totalItems}`);
  }
}

export function createAnalyzer(moduleIds?: string[]): Analyzer {
  return new Analyzer(moduleIds);
}
