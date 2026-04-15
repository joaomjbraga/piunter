import { exec, isCommandAvailable } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';

export class SnapModule {
  readonly id = 'snap';
  readonly name = 'Snap';
  readonly description = 'Remove snaps não utilizados e limpa cache do Snap';

  isAvailable(): boolean {
    return isCommandAvailable('snap');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    try {
      const listResult = await exec('snap', ['list']);
      if (listResult.success) {
        const lines = listResult.stdout
          .split('\n')
          .filter(l => l.trim() && !l.startsWith('Name'))
          .slice(0, 10);
        for (const line of lines) {
          const parts = line.trim().split(/\s+/);
          if (parts.length >= 2) {
            items.push({
              path: parts[0],
              size: 0,
              type: 'snap-app',
              description: `Snap: ${parts[0]}`,
            });
          }
        }
      }
    } catch (e) {
      logger.debug(`Snap command failed: ${(e as Error).message}`);
    }

    try {
      const duResult = await exec('du', ['-sb', '/var/lib/snapd/snaps']);
      if (duResult.success) {
        const match = duResult.stdout.match(/^(\d+)/);
        if (match) {
          const size = parseInt(match[1], 10);
          items.push({
            path: '/var/lib/snapd/snaps',
            size,
            type: 'snap-cache',
            description: 'Cache do Snap',
          });
          totalSize += size;
        }
      }
    } catch {
      // du command failed
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

    if (!this.isAvailable()) {
      result.errors.push('Snap não está instalado no sistema');
      result.success = false;
      return result;
    }

    const beforeAnalysis = await this.analyze();
    const beforeSize = beforeAnalysis.totalSize;

    if (dryRun) {
      logger.info(`[DRY-RUN] Snap: espaço em ${logger.formatBytes(beforeSize)}`);
      result.spaceFreed = beforeSize;
      logger.warn('Remoção manual necessária: snap remove <name>');
      return result;
    }

    try {
      const cacheCleanResult = await exec('snap', ['clean-cache']);
      if (cacheCleanResult.success) {
        logger.item(`${this.name}: Cache limpo`);
        result.itemsRemoved++;
      }
    } catch {
      logger.debug('Cache do snap não pode ser limpo');
    }

    try {
      const refreshResult = await exec('snap', ['refresh', '--list']);
      if (refreshResult.success) {
        const lines = refreshResult.stdout.split('\n').filter(l => l.trim() && !l.includes('Name'));
        if (lines.length > 0) {
          logger.info(`${this.name}: Snaps atualizáveis encontrados`);
          for (const line of lines) {
            logger.item(line.trim());
          }
        }
      }
    } catch {
      result.errors.push('Falha ao listar snaps');
    }

    logger.info(`${this.name}: Use 'snap remove <nome>' para remover snaps não utilizados`);

    const afterAnalysis = await this.analyze();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);

    return result;
  }
}

export const snapModule = new SnapModule();
