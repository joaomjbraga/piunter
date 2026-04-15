import { existsSync } from 'fs';
import { readdir, stat } from 'fs/promises';
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
    ];
  }

  isAvailable(): boolean {
    return true;
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];

    const results = await Promise.all(
      this.getAppImageDirs().map(async dir => {
        if (!existsSync(dir)) return [];

        try {
          const files = await readdir(dir);
          const appimageFiles = files.filter(
            f => f.endsWith('.AppImage') || f.endsWith('.appimage')
          );

          const fileStats = await Promise.all(
            appimageFiles.map(async file => {
              const fullPath = join(dir, file);
              try {
                const statInfo = await stat(fullPath);
                return { fullPath, size: statInfo.size, file };
              } catch {
                return null;
              }
            })
          );

          return fileStats.filter(
            (f): f is { fullPath: string; size: number; file: string } => f !== null
          );
        } catch {
          return [];
        }
      })
    );

    let totalSize = 0;
    for (const dirFiles of results) {
      for (const file of dirFiles) {
        items.push({
          path: file.fullPath,
          size: file.size,
          type: 'appimage',
          description: `AppImage: ${file.file}`,
        });
        totalSize += file.size;
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

    logger.info(`${analysis.items.length} AppImage(s) encontrados. Remocao manual necessaria:`);
    for (const item of analysis.items) {
      logger.item(item.path, logger.formatBytes(item.size));
    }
    logger.space();
    logger.info(`Para remover: rm "${analysis.items[0]?.path || '<arquivo>'}..."`);

    return result;
  }
}

export const appimageModule = new AppImageModule();
