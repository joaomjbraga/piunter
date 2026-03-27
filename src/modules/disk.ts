import { statSync } from 'fs';
import { exec } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';

export class LargeFilesModule {
  readonly id = 'large-files';
  readonly name = 'Arquivos Grandes';
  readonly description = 'Detecta arquivos grandes (remocao manual recomendada)';

  isAvailable(): boolean {
    return true;
  }

  async analyze(directory: string = '/home', thresholdMB: number = 100): Promise<AnalysisResult> {
    const threshold = thresholdMB * 1024 * 1024;
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const findResult = await exec('find', [
      directory,
      '-type',
      'f',
      '-size',
      `+${Math.floor(threshold / (1024 * 1024))}M`,
      '-not',
      '-path',
      '*/proc/*',
      '-not',
      '-path',
      '*/sys/*',
    ]);

    if (findResult.success && findResult.stdout) {
      const files = findResult.stdout.split('\n').filter(l => l.trim());

      for (const file of files) {
        try {
          const stat = statSync(file);
          items.push({
            path: file,
            size: stat.size,
            type: 'large-file',
            description: `Arquivo grande: ${file.split('/').pop()}`,
          });
          totalSize += stat.size;
        } catch {
          // Skip inaccessible files
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

    logger.warn('Use o modo interativo para selecionar arquivos específicos para remoção');

    return result;
  }
}

export const largeFilesModule = new LargeFilesModule();
