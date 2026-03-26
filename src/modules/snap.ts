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
        const lines = listResult.stdout.split('\n').filter(l => l.trim() && !l.startsWith('Name')).slice(0, 10);
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
    } catch {
      // Snap command failed
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

  async clean(dryRun: boolean = false, force: boolean = false): Promise<CleaningResult> {
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

    const analysis = await this.analyze();

    if (dryRun) {
      logger.info(`[DRY-RUN] Snap: limparía ${logger.formatBytes(analysis.totalSize)}`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    try {
      const listResult = await exec('snap', ['list']);
      if (listResult.success) {
        const lines = listResult.stdout.split('\n').filter(l => l.trim() && !l.startsWith('Name'));
        for (const line of lines) {
          const parts = line.trim().split(/\s+/);
          if (parts.length >= 2) {
            logger.item(`Snap: ${parts[0]}`);
          }
        }
      }
    } catch {
      result.errors.push('Falha ao listar snaps');
    }

    logger.item(`${this.name}: Limpeza manual recomendada - use 'snap remove <name>' para remover`);

    return result;
  }
}

export const snapModule = new SnapModule();
