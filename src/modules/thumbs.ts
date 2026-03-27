import { existsSync, readdirSync, rmSync } from 'fs';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getHomeDir } from '../utils/os.js';
import { logger } from '../utils/logger.js';
import { getDirSize } from '../utils/fs.js';

export class ThumbsModule {
  readonly id = 'thumbs';
  readonly name = 'Miniaturas';
  readonly description = 'Limpa miniaturas e thumbnails do sistema';

  private getThumbsDirs(): string[] {
    const home = getHomeDir();
    return [
      join(home, '.cache', 'thumbnails'),
      join(home, '.thumbnails'),
      join(home, '.local', 'share', 'thumbnails'),
    ];
  }

  isAvailable(): boolean {
    return true;
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    for (const dir of this.getThumbsDirs()) {
      if (!existsSync(dir)) continue;

      try {
        const size = getDirSize(dir);
        items.push({
          path: dir,
          size,
          type: 'thumbs-cache',
          description: `Miniaturas: ${dir.split('/').pop()}`,
        });
        totalSize += size;
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

    if (dryRun) {
      logger.info(`[DRY-RUN] Miniaturas: limparía ${logger.formatBytes(analysis.totalSize)}`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    for (const item of analysis.items) {
      if (!existsSync(item.path)) continue;

      let freedFromThis = 0;

      try {
        const entries = readdirSync(item.path);
        for (const entry of entries) {
          if (entry !== 'large' && entry !== 'normal') {
            const fullPath = join(item.path, entry);
            try {
              const entrySize = getDirSize(fullPath);
              rmSync(fullPath, { recursive: true, force: true });
              freedFromThis += entrySize;
              result.itemsRemoved++;
            } catch {
              // Skip
            }
          }
        }

        for (const subdir of ['large', 'normal']) {
          const subPath = join(item.path, subdir);
          if (existsSync(subPath)) {
            try {
              const subSize = getDirSize(subPath);
              rmSync(subPath, { recursive: true, force: true });
              freedFromThis += subSize;
              result.itemsRemoved++;
            } catch {
              // Skip
            }
          }
        }

        result.spaceFreed += freedFromThis;
        if (freedFromThis > 0) {
          logger.item(
            `${this.name}: ${item.path.split('/').pop()} (${logger.formatBytes(freedFromThis)})`
          );
        }
      } catch {
        result.errors.push(`Falha ao limpar ${item.path}`);
      }
    }

    return result;
  }
}

export const thumbsModule = new ThumbsModule();
