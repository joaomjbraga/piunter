import { existsSync } from 'fs';
import { join } from 'path';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { getHomeDir } from '../utils/os.js';
import { exec, isCommandAvailable } from '../utils/exec.js';
import { logger } from '../utils/logger.js';
import { getDirSizeAsync } from '../utils/fs.js';

abstract class PackageCacheModule {
  abstract readonly id: string;
  abstract readonly name: string;
  abstract readonly description: string;
  abstract getCachePaths(): string[];

  isAvailable(): boolean {
    return isCommandAvailable(this.id);
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];

    const results = await Promise.all(
      this.getCachePaths().map(async cachePath => {
        if (!existsSync(cachePath)) return null;
        try {
          const size = await getDirSizeAsync(cachePath);
          return { cachePath, size };
        } catch {
          return null;
        }
      })
    );

    let totalSize = 0;
    for (const result of results) {
      if (result) {
        items.push({
          path: result.cachePath,
          size: result.size,
          type: 'directory',
          description: `Cache do ${this.name} (${result.cachePath})`,
        });
        totalSize += result.size;
      }
    }

    return { module: this.id, items, totalSize };
  }

  abstract getCleanCommand(): string[];

  async clean(dryRun: boolean = false): Promise<CleaningResult> {
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
      logger.info(
        `[DRY-RUN] Limparía ${logger.formatBytes(analysis.totalSize)} do cache ${this.name.toLowerCase()}`
      );
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const commandResult = await exec(this.id, this.getCleanCommand());

    if (commandResult.success) {
      result.success = true;
      result.spaceFreed = analysis.totalSize;
      result.itemsRemoved = 1;
      logger.item(`${this.name}: Cache limpo`, logger.formatBytes(analysis.totalSize));
    } else {
      result.success = false;
      result.errors.push(commandResult.stderr || `Falha ao limpar cache ${this.name}`);
    }

    return result;
  }
}

export class NpmModule extends PackageCacheModule {
  readonly id = 'npm';
  readonly name = 'NPM';
  readonly description = 'Limpa cache do npm';

  getCachePaths(): string[] {
    return [join(getHomeDir(), '.npm')];
  }

  getCleanCommand(): string[] {
    return ['cache', 'clean', '--force'];
  }
}

export class YarnModule extends PackageCacheModule {
  readonly id = 'yarn';
  readonly name = 'Yarn';
  readonly description = 'Limpa cache do Yarn';

  getCachePaths(): string[] {
    return [join(getHomeDir(), '.cache', 'yarn'), join(getHomeDir(), '.yarn', 'cache')];
  }

  getCleanCommand(): string[] {
    return ['cache', 'clean'];
  }
}

export class PnpmModule extends PackageCacheModule {
  readonly id = 'pnpm';
  readonly name = 'PNPM';
  readonly description = 'Limpa cache do PNPM';

  getCachePaths(): string[] {
    const paths = [join(getHomeDir(), '.pnpm-store')];
    const localPath = join(getHomeDir(), '.local', 'share', 'pnpm', 'cache');
    if (existsSync(localPath)) {
      paths.push(localPath);
    }
    return paths;
  }

  getCleanCommand(): string[] {
    return ['store', 'prune'];
  }
}

export const npmModule = new NpmModule();
export const yarnModule = new YarnModule();
export const pnpmModule = new PnpmModule();
