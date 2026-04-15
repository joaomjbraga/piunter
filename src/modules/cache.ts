import { existsSync, rmSync } from 'fs';
import { readdir, stat } from 'fs/promises';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getCacheDir } from '../utils/os.js';
import { logger } from '../utils/logger.js';
import { getDirSizeAsync } from '../utils/fs.js';

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
      const entries = await readdir(cacheDir);
      const stats = await Promise.all(
        entries.map(async entry => {
          const fullPath = join(cacheDir, entry);
          try {
            const statInfo = await stat(fullPath);
            return { entry, fullPath, statInfo, error: null };
          } catch {
            return { entry, fullPath, statInfo: null, error: true };
          }
        })
      );

      const dirPromises = stats
        .filter(s => !s.error && s.statInfo?.isDirectory())
        .map(async s => {
          const size = await getDirSizeAsync(s.fullPath);
          return { ...s, size, type: 'directory' as const };
        });

      const fileStats = stats.filter(s => !s.error && s.statInfo && !s.statInfo.isDirectory());

      const dirsWithSize = await Promise.all(dirPromises);

      for (const dir of dirsWithSize) {
        items.push({
          path: dir.fullPath,
          size: dir.size,
          type: 'directory',
          description: `Diretório de cache: ${dir.entry}`,
        });
        totalSize += dir.size;
      }

      for (const file of fileStats) {
        if (file.statInfo) {
          items.push({
            path: file.fullPath,
            size: file.statInfo.size,
            type: 'file',
            description: `Arquivo de cache: ${file.entry}`,
          });
          totalSize += file.statInfo.size;
        }
      }
    } catch (e) {
      logger.debug(`Cache dir not readable: ${(e as Error).message}`);
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

    const dirsToSkip = ['thumbnails', 'thumbnail', 'icon-cache'];

    for (const item of analysis.items) {
      const itemBasename = item.path.split('/').pop() || '';
      const shouldSkip = dirsToSkip.some(d => itemBasename === d);
      if (shouldSkip) continue;

      if (!dryRun) {
        try {
          if (item.type === 'directory') {
            rmSync(item.path, { recursive: true, force: true });
          } else {
            rmSync(item.path, { force: true });
          }
          result.spaceFreed += item.size;
          result.itemsRemoved++;
        } catch (e: unknown) {
          result.errors.push(`Falha ao remover ${item.path}: ${(e as Error).message}`);
        }
      }
    }

    if (dryRun) {
      const cleanableItems = analysis.items.filter(i => !dirsToSkip.some(d => i.path.includes(d)));
      result.spaceFreed = cleanableItems.reduce((sum, i) => sum + i.size, 0);
      result.itemsRemoved = cleanableItems.length;
    }

    if (result.spaceFreed > 0) {
      logger.item(`${this.name}: Cache limpo`, logger.formatBytes(result.spaceFreed));
    }

    return result;
  }
}

export const cacheModule = new CacheModule();
