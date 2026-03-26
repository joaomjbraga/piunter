import { existsSync, readdirSync, statSync } from 'fs';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getHomeDir } from '../utils/os.js';
import { exec, isCommandAvailable } from '../utils/exec.js';
import { logger } from '../utils/logger.js';

export class NpmModule {
  readonly id = 'npm';
  readonly name = 'NPM';
  readonly description = 'Limpa cache do npm';

  isAvailable(): boolean {
    return isCommandAvailable('npm');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const npmCachePath = join(getHomeDir(), '.npm');

    if (existsSync(npmCachePath)) {
      try {
        const size = this.getDirSize(npmCachePath);
        items.push({
          path: npmCachePath,
          size,
          type: 'directory',
          description: 'Cache do npm (~/.npm)',
        });
        totalSize = size;
      } catch {
        // Not accessible
      }
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

    if (analysis.totalSize === 0) {
      return result;
    }

    if (dryRun) {
      logger.info(`[DRY-RUN] Limparía ${logger.formatBytes(analysis.totalSize)} do cache npm`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const commandResult = await exec('npm', ['cache', 'clean', '--force']);

    if (commandResult.success) {
      result.success = true;
      result.spaceFreed = analysis.totalSize;
      result.itemsRemoved = 1;
      logger.item(`${this.name}: Cache limpo`, logger.formatBytes(analysis.totalSize));
    } else {
      result.success = false;
      result.errors.push(commandResult.stderr || 'Falha ao limpar cache npm');
    }

    return result;
  }
}

export class YarnModule {
  readonly id = 'yarn';
  readonly name = 'Yarn';
  readonly description = 'Limpa cache do Yarn';

  isAvailable(): boolean {
    return isCommandAvailable('yarn');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const yarnCachePath = join(getHomeDir(), '.yarn', 'cache');

    if (existsSync(yarnCachePath)) {
      try {
        const size = this.getDirSize(yarnCachePath);
        items.push({
          path: yarnCachePath,
          size,
          type: 'directory',
          description: 'Cache do Yarn (~/.yarn/cache)',
        });
        totalSize = size;
      } catch {
        // Not accessible
      }
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

    if (analysis.totalSize === 0) {
      return result;
    }

    if (dryRun) {
      logger.info(`[DRY-RUN] Limparía ${logger.formatBytes(analysis.totalSize)} do cache yarn`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const commandResult = await exec('yarn', ['cache', 'clean']);

    if (commandResult.success) {
      result.success = true;
      result.spaceFreed = analysis.totalSize;
      result.itemsRemoved = 1;
      logger.item(`${this.name}: Cache limpo`, logger.formatBytes(analysis.totalSize));
    } else {
      result.success = false;
      result.errors.push(commandResult.stderr || 'Falha ao limpar cache yarn');
    }

    return result;
  }
}

export class PnpmModule {
  readonly id = 'pnpm';
  readonly name = 'PNPM';
  readonly description = 'Limpa cache do PNPM';

  isAvailable(): boolean {
    return isCommandAvailable('pnpm');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const pnpmCachePath = join(getHomeDir(), '.pnpm-store');

    if (existsSync(pnpmCachePath)) {
      try {
        const size = this.getDirSize(pnpmCachePath);
        items.push({
          path: pnpmCachePath,
          size,
          type: 'directory',
          description: 'Cache do PNPM (~/.pnpm-store)',
        });
        totalSize = size;
      } catch {
        // Not accessible
      }
    }

    const pnpmCachePath2 = join(getHomeDir(), '.local', 'share', 'pnpm', 'cache');
    if (existsSync(pnpmCachePath2)) {
      try {
        const size = this.getDirSize(pnpmCachePath2);
        items.push({
          path: pnpmCachePath2,
          size,
          type: 'directory',
          description: 'Cache do PNPM (local)',
        });
        totalSize += size;
      } catch {
        // Not accessible
      }
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

    if (analysis.totalSize === 0) {
      return result;
    }

    if (dryRun) {
      logger.info(`[DRY-RUN] Limparía ${logger.formatBytes(analysis.totalSize)} do cache pnpm`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const commandResult = await exec('pnpm', ['store', 'prune']);

    if (commandResult.success) {
      result.success = true;
      result.spaceFreed = analysis.totalSize;
      result.itemsRemoved = 1;
      logger.item(`${this.name}: Cache limpo`, logger.formatBytes(analysis.totalSize));
    } else {
      result.success = false;
      result.errors.push(commandResult.stderr || 'Falha ao limpar cache pnpm');
    }

    return result;
  }
}

export const npmModule = new NpmModule();
export const yarnModule = new YarnModule();
export const pnpmModule = new PnpmModule();
