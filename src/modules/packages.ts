import { exec } from '../utils/exec.js';
import { getDistroInfo } from '../utils/os.js';
import type { AnalysisResult, CleaningResult, PackageManager } from '../types/index.js';
import { logger } from '../utils/logger.js';
import { parseSize } from '../utils/fs.js';

export class PackagesModule {
  readonly id = 'packages';
  readonly name = 'Gerenciador de Pacotes';
  readonly description = 'Limpa cache e remove pacotes órfãos do sistema';

  private packageManager: PackageManager;

  constructor() {
    this.packageManager = getDistroInfo().packageManager;
  }

  isAvailable(): boolean {
    return this.packageManager !== 'unknown';
  }

  async analyze(): Promise<AnalysisResult> {
    switch (this.packageManager) {
      case 'apt':
        return this.analyzeApt();
      case 'pacman':
        return this.analyzePacman();
      case 'dnf':
        return this.analyzeDnf();
      default:
        return { module: this.id, items: [], totalSize: 0 };
    }
  }

  private async analyzeApt(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    try {
      const cacheResult = await exec('du', ['-sb', '/var/cache/apt/archives']);
      if (cacheResult.success) {
        const match = cacheResult.stdout.match(/^(\d+)/);
        if (match) {
          const size = parseInt(match[1], 10);
          items.push({
            path: '/var/cache/apt/archives',
            size,
            type: 'apt-cache',
            description: 'Cache do APT (/var/cache/apt/archives)',
          });
          totalSize += size;
        }
      }
    } catch {
      items.push({
        path: '/var/cache/apt/archives',
        size: 0,
        type: 'apt-cache',
        description: 'Cache do APT (sem permissão para calcular)',
      });
    }

    try {
      const orphansResult = await exec('apt', ['autoremove', '--dry-run']);
      if (orphansResult.success) {
        const match = orphansResult.stdout.match(/(\d+)\s+package/i);
        if (match) {
          const count = parseInt(match[1], 10);
          items.push({
            path: 'apt-orphans',
            size: count * 50 * 1024 * 1024,
            type: 'apt-orphans',
            description: `Pacotes órfãos APT (${count})`,
          });
          totalSize += count * 50 * 1024 * 1024;
        }
      }
    } catch {
      // Dry-run pode falhar se não houver órfãos
    }

    return { module: this.id, items, totalSize };
  }

  private async analyzePacman(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const cacheResult = await exec('du', ['-sb', '/var/cache/pacman/pkg']);
    if (cacheResult.success) {
      const match = cacheResult.stdout.match(/^(\d+)/);
      if (match) {
        const size = parseInt(match[1], 10);
        items.push({
          path: '/var/cache/pacman/pkg',
          size,
          type: 'pacman-cache',
          description: 'Cache do Pacman (/var/cache/pacman/pkg)',
        });
        totalSize += size;
      }
    }

    const orphansResult = await exec('paccache', ['-dk0']);
    if (orphansResult.success) {
      const match = orphansResult.stdout.match(/([\d.]+)\s*(?:KiB|MiB|GiB)/i);
      if (match) {
        const size = parseSize(match[1] + ' KiB');
        items.push({
          path: 'pacman-orphans',
          size,
          type: 'pacman-orphans',
          description: 'Versões antigas no cache (paccache -dk0)',
        });
        totalSize += size;
      }
    }

    return { module: this.id, items, totalSize };
  }

  private async analyzeDnf(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const cacheResult = await exec('du', ['-sb', '/var/cache/dnf']);
    if (cacheResult.success) {
      const match = cacheResult.stdout.match(/^(\d+)/);
      if (match) {
        const size = parseInt(match[1], 10);
        items.push({
          path: '/var/cache/dnf',
          size,
          type: 'dnf-cache',
          description: 'Cache do DNF (/var/cache/dnf)',
        });
        totalSize += size;
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
      logger.info(`[DRY-RUN] ${this.name}: limparía ${logger.formatBytes(analysis.totalSize)}`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    switch (this.packageManager) {
      case 'apt':
        return this.cleanApt(analysis);
      case 'pacman':
        return this.cleanPacman(analysis);
      case 'dnf':
        return this.cleanDnf(analysis);
      default:
        result.errors.push('Gerenciador de pacotes não suportado');
        return result;
    }
  }

  private async cleanApt(analysis: AnalysisResult): Promise<CleaningResult> {
    const beforeSize = analysis.totalSize;
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    const cleanResult = await exec('apt', ['clean'], { sudo: true });
    if (cleanResult.success) {
      logger.item(`APT: Cache limpo`);
      result.itemsRemoved++;
    } else {
      result.errors.push('Falha ao limpar cache APT (verifique se tem privilégios sudo)');
    }

    const autoremoveResult = await exec('apt', ['autoremove', '-y'], { sudo: true });
    if (autoremoveResult.success) {
      logger.item(`APT: Pacotes órfãos removidos`);
      result.itemsRemoved++;
    }

    const afterAnalysis = await this.analyzeApt();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);

    return result;
  }

  private async cleanPacman(analysis: AnalysisResult): Promise<CleaningResult> {
    const beforeSize = analysis.totalSize;
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    const cleanResult = await exec('pacman', ['-Sc', '--noconfirm'], { sudo: true });
    if (cleanResult.success) {
      logger.item(`Pacman: Cache limpo (mantém última versão)`);
      result.itemsRemoved++;
    } else {
      result.errors.push('Falha ao limpar cache Pacman (verifique se tem privilégios sudo)');
    }

    try {
      const orphansResult = await exec('pacman', ['-Qtdq']);
      if (orphansResult.success && orphansResult.stdout.trim()) {
        const orphanPackages = orphansResult.stdout
          .trim()
          .split('\n')
          .filter(p => p);
        for (const pkg of orphanPackages) {
          const removeResult = await exec('pacman', ['-Rns', pkg, '--noconfirm'], { sudo: true });
          if (removeResult.success) {
            result.itemsRemoved++;
          }
        }
        if (orphanPackages.length > 0) {
          logger.item(`Pacman: ${orphanPackages.length} pacotes órfãos removidos`);
        }
      }
    } catch {
      logger.debug('Nenhum pacote órfão encontrado');
    }

    const afterAnalysis = await this.analyzePacman();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);

    return result;
  }

  private async cleanDnf(analysis: AnalysisResult): Promise<CleaningResult> {
    const beforeSize = analysis.totalSize;
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    const cleanResult = await exec('dnf', ['clean', 'all'], { sudo: true });
    if (cleanResult.success) {
      logger.item(`DNF: Cache limpo`);
      result.itemsRemoved++;
    } else {
      result.errors.push('Falha ao limpar cache DNF (verifique se tem privilégios sudo)');
    }

    const autoremoveResult = await exec('dnf', ['autoremove', '-y'], { sudo: true });
    if (autoremoveResult.success) {
      logger.item(`DNF: Pacotes órfãos removidos`);
      result.itemsRemoved++;
    }

    const afterAnalysis = await this.analyzeDnf();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);

    return result;
  }
}

export const packagesModule = new PackagesModule();
