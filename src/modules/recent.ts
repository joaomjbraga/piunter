import { existsSync, rmSync } from 'fs';
import { stat } from 'fs/promises';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getHomeDir } from '../utils/os.js';
import { logger } from '../utils/logger.js';

export class RecentFilesModule {
  readonly id = 'recent';
  readonly name = 'Arquivos Recentes';
  readonly description = 'Limpa registros de arquivos recentes';

  private getRecentFilesDirs(): string[] {
    const home = getHomeDir();
    return [
      join(home, '.local', 'share', 'recently-used.xbel'),
      join(home, '.local', 'share', 'recently-used.xbel.bak'),
      join(home, '.gtk-bookmarks'),
      join(home, '.local', 'share', 'zeitgeist'),
    ];
  }

  isAvailable(): boolean {
    return true;
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];

    const results = await Promise.all(
      this.getRecentFilesDirs().map(async filePath => {
        if (!existsSync(filePath)) return null;
        try {
          const statInfo = await stat(filePath);
          return { filePath, size: statInfo.size };
        } catch {
          return null;
        }
      })
    );

    let totalSize = 0;
    for (const result of results) {
      if (result) {
        items.push({
          path: result.filePath,
          size: result.size,
          type: 'recent-files',
          description: `Arquivos recentes: ${result.filePath.split('/').pop()}`,
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
      logger.info(
        `[DRY-RUN] Arquivos recentes: limparía ${logger.formatBytes(analysis.totalSize)}`
      );
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    for (const item of analysis.items) {
      if (!existsSync(item.path)) continue;

      try {
        if (item.path.endsWith('.xbel') || item.path.endsWith('.gtk-bookmarks')) {
          rmSync(item.path, { force: true });
        }

        result.spaceFreed += item.size;
        result.itemsRemoved++;
        logger.item(`${this.name}: ${item.path.split('/').pop()}`);
      } catch {
        result.errors.push(`Falha ao limpar ${item.path}`);
      }
    }

    if (result.itemsRemoved === 0) {
      logger.info('Nenhum registro de arquivos recentes encontrado');
    }

    return result;
  }
}

export const recentFilesModule = new RecentFilesModule();
