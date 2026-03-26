import { readdirSync, statSync, existsSync, rmSync } from 'fs';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getCacheDir } from '../utils/os.js';
import { logger } from '../utils/logger.js';

export class CacheModule {
  readonly id = 'cache';
  readonly name = 'Cache do Usuário';
  readonly description = 'Limpa cache geral do usuário (~/.cache)';

  isAvailable(): boolean {
    return existsSync(getCacheDir());
  }

  async analyze(): Promise<AnalysisResult> {
    const cacheDir = getCacheDir();
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    if (!existsSync(cacheDir)) {
      return { module: this.id, items, totalSize: 0 };
    }

    try {
      const entries = readdirSync(cacheDir);
      for (const entry of entries) {
        const fullPath = join(cacheDir, entry);
        try {
          const stat = statSync(fullPath);
          if (stat.isDirectory()) {
            const size = this.getDirSize(fullPath);
            items.push({
              path: fullPath,
              size,
              type: 'directory',
              description: `Diretório de cache: ${entry}`,
            });
            totalSize += size;
          } else {
            items.push({
              path: fullPath,
              size: stat.size,
              type: 'file',
              description: `Arquivo de cache: ${entry}`,
            });
            totalSize += stat.size;
          }
        } catch {
          // Skip inaccessible entries
        }
      }
    } catch {
      // Cache dir might not be readable
    }

    return { module: this.id, items, totalSize };
  }

  private getDirSize(dirPath: string): number {
    let size = 0;
    try {
      const entries = readdirSync(dirPath);
      for (const entry of entries) {
        const fullPath = join(dirPath, entry);
        try {
          const stat = statSync(fullPath);
          if (stat.isDirectory()) {
            size += this.getDirSize(fullPath);
          } else {
            size += stat.size;
          }
        } catch {
          // Skip
        }
      }
    } catch {
      // Skip
    }
    return size;
  }

  async clean(dryRun: boolean = false, _force: boolean = false): Promise<CleaningResult> {
    const analysis = await this.analyze();
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    const safeDirs = ['.cache'];

    for (const item of analysis.items) {
      if (safeDirs.some(d => item.path.includes(d)) && !item.path.includes('thumbnails')) {
        continue;
      }
      if (item.path.includes('thumbnail') || item.path.includes('thumbnails')) {
        result.itemsRemoved++;
        result.spaceFreed += item.size;
        if (!dryRun) {
          try {
            if (item.type === 'directory') {
              rmSync(item.path, { recursive: true, force: true });
            } else {
              rmSync(item.path, { force: true });
            }
          } catch (e: unknown) {
            result.errors.push(`Falha ao remover ${item.path}: ${(e as Error).message}`);
          }
        }
      }
    }

    if (result.spaceFreed > 0) {
      logger.item(`${this.name}: Thumbnails e outros caches`, logger.formatBytes(result.spaceFreed));
    }

    return result;
  }
}

export const cacheModule = new CacheModule();
