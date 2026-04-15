import { stat } from 'fs/promises';
import { exec } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';
import { loadConfig } from '../utils/config.js';

export class LargeFilesModule {
  readonly id = 'large-files';
  readonly name = 'Arquivos Grandes';
  readonly description = 'Detecta arquivos grandes (remocao manual recomendada)';

  isAvailable(): boolean {
    return true;
  }

  async analyze(thresholdMB?: number): Promise<AnalysisResult> {
    const config = loadConfig();
    const threshold = thresholdMB ?? config.thresholds.largeFilesMB;
    const items: AnalysisResult['items'] = [];

    const findResult = await exec('find', [
      '/home',
      '-type',
      'f',
      '-size',
      `+${threshold}M`,
      '-not',
      '-path',
      '*/proc/*',
      '-not',
      '-path',
      '*/sys/*',
    ]);

    let totalSize = 0;

    if (findResult.success && findResult.stdout) {
      const files = findResult.stdout.split('\n').filter(l => l.trim());

      const fileStats = await Promise.all(
        files.map(async file => {
          try {
            const statInfo = await stat(file);
            return { file, size: statInfo.size };
          } catch {
            return null;
          }
        })
      );

      for (const fileStat of fileStats) {
        if (fileStat) {
          items.push({
            path: fileStat.file,
            size: fileStat.size,
            type: 'large-file',
            description: `Arquivo grande: ${fileStat.file.split('/').pop()}`,
          });
          totalSize += fileStat.size;
        }
      }
    }

    return { module: this.id, items, totalSize };
  }

  async clean(dryRun: boolean = false): Promise<CleaningResult> {
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    const analysis = await this.analyze();

    if (analysis.items.length === 0) {
      logger.info('Nenhum arquivo grande encontrado');
      return result;
    }

    if (dryRun) {
      logger.info(
        `[DRY-RUN] Removería ${analysis.items.length} arquivos (${logger.formatBytes(analysis.totalSize)})`
      );
      return result;
    }

    logger.info(
      `${analysis.items.length} arquivo(s) grande(s) encontrado(s). Remocao manual necessaria.`
    );
    logger.space();
    for (const item of analysis.items.slice(0, 5)) {
      logger.item(item.path, logger.formatBytes(item.size));
    }
    if (analysis.items.length > 5) {
      logger.info(`  ... e mais ${analysis.items.length - 5} arquivo(s)`);
    }
    logger.space();
    logger.info(`Total: ${logger.formatBytes(analysis.totalSize)}`);

    return result;
  }
}

export const largeFilesModule = new LargeFilesModule();
