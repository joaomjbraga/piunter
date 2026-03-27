import { exec, isCommandAvailable } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';
import { parseSize } from '../utils/fs.js';

export class FlatpakModule {
  readonly id = 'flatpak';
  readonly name = 'Flatpak';
  readonly description = 'Remove apps Flatpak não utilizados e limpa cache';

  isAvailable(): boolean {
    return isCommandAvailable('flatpak');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const listResult = await exec('flatpak', ['list', '--app', '--columns=size,name']);

    if (listResult.success) {
      const lines = listResult.stdout.split('\n').filter(l => l.trim());
      for (const line of lines) {
        const parts = line.split('\t');
        if (parts.length >= 2) {
          const size = parseSize(parts[0]);
          const name = parts[1];
          items.push({
            path: name,
            size,
            type: 'flatpak-app',
            description: `App Flatpak: ${name}`,
          });
          totalSize += size;
        }
      }
    }

    const cacheResult = await this.getCacheSize();
    if (cacheResult > 0) {
      items.push({
        path: 'flatpak-cache',
        size: cacheResult,
        type: 'flatpak-cache',
        description: 'Cache do Flatpak',
      });
      totalSize += cacheResult;
    }

    return { module: this.id, items, totalSize };
  }

  private async getCacheSize(): Promise<number> {
    const cachePath = '/var/tmp/flatpak-cache';
    const cacheResult = await exec('du', ['-sb', cachePath]);

    if (cacheResult.success) {
      const match = cacheResult.stdout.match(/^(\d+)/);
      if (match) {
        return parseInt(match[1], 10);
      }
    }

    return 0;
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
      result.errors.push('Flatpak não está instalado no sistema');
      result.success = false;
      return result;
    }

    const beforeAnalysis = await this.analyze();
    const beforeSize = beforeAnalysis.items.reduce((sum, i) => sum + i.size, 0);

    if (dryRun) {
      logger.info(`[DRY-RUN] Flatpak: limparía ${logger.formatBytes(beforeSize)}`);
      result.spaceFreed = beforeSize;
      return result;
    }

    try {
      const uninstallResult = await exec('flatpak', ['uninstall', '--unused', '-y'], {
        sudo: true,
      });
      if (uninstallResult.success) {
        const match = uninstallResult.stdout.match(/(\d+)/);
        if (match) {
          result.itemsRemoved += parseInt(match[1], 10);
        }
        logger.item(`${this.name}: Apps não utilizados removidos`);
      }
    } catch {
      result.errors.push('Falha ao desinstalar Flatpaks não utilizados');
    }

    try {
      const repairResult = await exec('flatpak', ['repair', '--system'], { sudo: true });
      if (repairResult.success) {
        logger.item(`${this.name}: Repair concluído`);
      }
    } catch {
      result.errors.push('Falha no repair do Flatpak');
    }

    const afterAnalysis = await this.analyze();
    const afterSize = afterAnalysis.items.reduce((sum, i) => sum + i.size, 0);
    result.spaceFreed = Math.max(0, beforeSize - afterSize);

    return result;
  }
}

export const flatpakModule = new FlatpakModule();
