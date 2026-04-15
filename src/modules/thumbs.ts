import { existsSync, rmSync, readdirSync, statSync } from 'fs';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getHomeDir } from '../utils/os.js';
import { logger } from '../utils/logger.js';
import { getDirSizeAsync, getDirSize } from '../utils/fs.js';

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

    const results = await Promise.all(
      this.getThumbsDirs().map(async dir => {
        if (!existsSync(dir)) return null;
        try {
          const size = await getDirSizeAsync(dir);
          return { dir, size };
        } catch {
          return null;
        }
      })
    );

    let totalSize = 0;
    for (const result of results) {
      if (result) {
        items.push({
          path: result.dir,
          size: result.size,
          type: 'thumbs-cache',
          description: `Miniaturas: ${result.dir.split('/').pop()}`,
        });
        totalSize += result.size;
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

    const sevenDaysAgo = Date.now() - 7 * 24 * 60 * 60 * 1000;

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
              const subEntries = readdirSync(subPath);
              for (const entry of subEntries) {
                const entryPath = join(subPath, entry);
                const stat = statSync(entryPath);
                if (stat.mtimeMs < sevenDaysAgo) {
                  const entrySize = stat.size;
                  rmSync(entryPath, { force: true });
                  freedFromThis += entrySize;
                  result.itemsRemoved++;
                }
              }
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
