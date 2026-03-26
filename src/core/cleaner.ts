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

    if (this.dryRun) {
      logger.warn('Modo DRY-RUN ativo - nenhuma alteração será feita');
    }

    for (const module of this.modules) {
      if (!module.isAvailable()) {
        logger.debug(`${module.name} não disponível, pulando`);
        continue;
      }

      logger.startSpinner(`Limpando ${module.name}...`);
      try {
        const result = await module.clean(this.dryRun, this.force);
        results.push(result);

        if (result.success) {
          logger.stopSpinner(true, `${module.name}: ${logger.formatBytes(result.spaceFreed)} liberada`);
        } else {
          logger.stopSpinner(false, `${module.name}: Erro`);
        }

        if (result.errors.length > 0) {
          errors.push(...result.errors);
        }
      } catch (error) {
        logger.stopSpinner(false, `Erro ao limpar ${module.name}`);
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
    logger.title('Relatório de Limpeza');

    const duration = report.endTime.getTime() - report.startTime.getTime();
    const durationStr = duration < 60000 
      ? `${Math.round(duration / 1000)}s`
      : `${Math.round(duration / 60000)}min`;

    logger.info(`Tempo de execução: ${durationStr}`);
    logger.space();

    if (report.modules.length > 0) {
      logger.subtitle('Módulos processados:');
      for (const result of report.modules) {
        if (result.success) {
          logger.item(`${result.module}: ${logger.formatBytes(result.spaceFreed)}`, `${result.itemsRemoved} itens`);
        } else {
          logger.item(`${result.module}: ERRO`, 'Falhou');
        }
      }
    }

    logger.space();
    logger.subtitle('Resumo:');
    logger.info(`Espaço liberado: ${logger.formatBytes(report.totalSpaceFreed)}`);
    logger.info(`Itens removidos: ${report.totalItemsRemoved}`);
    logger.info(`Módulos com erro: ${report.errors.length}`);

    if (report.errors.length > 0) {
      logger.space();
      logger.warn('Erros encontrados:');
      for (const error of report.errors) {
        logger.item(error);
      }
    }

    logger.space();
    if (this.dryRun) {
      logger.success('Dry-run concluído - use sem --dry-run para aplicar');
    } else {
      logger.success('Limpeza concluída com sucesso!');
    }
  }
}

export function createCleaner(moduleIds: string[], options: CleanOptions): Cleaner {
  return new Cleaner(moduleIds, options);
}
