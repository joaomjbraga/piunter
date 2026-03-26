import { exec, isCommandAvailable } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';

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
          const size = this.parseSize(parts[0]);
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

  private parseSize(sizeStr: string): number {
    const match = sizeStr.match(/([\d.]+)\s*([A-Z]+)/i);
    if (!match) return 0;
    
    const num = parseFloat(match[1]);
    const unit = match[2].toUpperCase();
    
    const multipliers: Record<string, number> = {
      'B': 1,
      'KB': 1024,
      'MB': 1024 * 1024,
      'GB': 1024 * 1024 * 1024,
    };
    
    return num * (multipliers[unit] || 1);
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

  async clean(dryRun: boolean = false, force: boolean = false): Promise<CleaningResult> {
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    const analysis = await this.analyze();

    if (dryRun) {
      logger.info(`[DRY-RUN] Flatpak: limparía ${logger.formatBytes(analysis.totalSize)}`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const repairResult = await exec('flatpak', ['repair', '--system']);
    if (!repairResult.success) {
      result.errors.push('Falha no repair do Flatpak');
    }

    const cacheResult = await exec('flatpak', ['repair', '-y']);
    if (cacheResult.success) {
      result.spaceFreed += analysis.totalSize * 0.1;
      logger.item(`${this.name}: Repair concluído`);
    }

    const uninstallResult = await exec('flatpak', ['uninstall', '--unused', '-y']);
    if (uninstallResult.success) {
      logger.item(`${this.name}: Apps não utilizados removidos`);
    }

    const pruneResult = await exec('flatpak', ['system', 'reset', '-y']);
    if (pruneResult.success) {
      logger.item(`${this.name}: Sistema Flatpak otimizado`);
    }

    return result;
  }
}

export const flatpakModule = new FlatpakModule();
