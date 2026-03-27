import { existsSync, readdirSync, statSync } from 'fs';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getHomeDir } from '../utils/os.js';
import { logger } from '../utils/logger.js';

export class AppImageModule {
  readonly id = 'appimage';
  readonly name = 'AppImage';
  readonly description = 'Detecta e limpa arquivos AppImage antigos';

  private getAppImageDirs(): string[] {
    return [
      join(getHomeDir(), 'Applications'),
      join(getHomeDir(), 'Downloads'),
      join(getHomeDir(), '.local', 'bin'),
      '/usr/local/bin',
    ];
  }

  isAvailable(): boolean {
    return true;
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    for (const dir of this.getAppImageDirs()) {
      if (!existsSync(dir)) continue;

      try {
        const files = readdirSync(dir);
        for (const file of files) {
          if (file.endsWith('.AppImage') || file.endsWith('.appimage')) {
            const fullPath = join(dir, file);
            try {
              const stat = statSync(fullPath);
              items.push({
                path: fullPath,
                size: stat.size,
                type: 'appimage',
                description: `AppImage: ${file}`,
              });
              totalSize += stat.size;
            } catch {
              // Skip inaccessible files
            }
          }
        }
      } catch {
        // Skip inaccessible directories
      }
    }

    return { module: this.id, items, totalSize };
  }

  async clean(dryRun: boolean = false): Promise<CleaningResult> {
    const analysis = await this.analyze();
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    if (analysis.items.length === 0) {
      logger.info('Nenhum AppImage encontrado');
      return result;
    }

    if (dryRun) {
      logger.info(
        `[DRY-RUN] AppImage: limparía ${analysis.items.length} arquivos (${logger.formatBytes(analysis.totalSize)})`
      );
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    logger.warn('AppImages encontrados. Remoção manual recomendada:');
    for (const item of analysis.items) {
      logger.item(item.path, logger.formatBytes(item.size));
    }

    return result;
  }
}

export const appimageModule = new AppImageModule();
