import { statSync, existsSync } from 'fs';
import { exec } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';

export class LargeFilesModule {
  readonly id = 'large-files';
  readonly name = 'Arquivos Grandes';
  readonly description = 'Detecta e limpa arquivos maiores que o limite especificado';

  private threshold: number;

  constructor(thresholdMB: number = 100) {
    this.threshold = thresholdMB * 1024 * 1024;
  }

  setThreshold(mb: number): void {
    this.threshold = mb * 1024 * 1024;
  }

  isAvailable(): boolean {
    return true;
  }

  async analyze(directory: string = '/home', thresholdMB?: number): Promise<AnalysisResult> {
    if (thresholdMB) {
      this.threshold = thresholdMB * 1024 * 1024;
    }

    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const findResult = await exec('find', [
      directory,
      '-type', 'f',
      '-size', `+${Math.floor(this.threshold / (1024 * 1024))}M`,
      '-not', '-path', '*/proc/*',
      '-not', '-path', '*/sys/*',
    ]);

    if (findResult.success && findResult.stdout) {
      const files = findResult.stdout.split('\n').filter(l => l.trim());
      
      for (const file of files) {
        try {
          const stat = statSync(file);
          items.push({
            path: file,
            size: stat.size,
            type: 'large-file',
            description: `Arquivo grande: ${file.split('/').pop()}`,
          });
          totalSize += stat.size;
        } catch {
          // Skip inaccessible files
        }
      }
    }

    return { module: this.id, items, totalSize };
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

    if (analysis.items.length === 0) {
      logger.info('Nenhum arquivo grande encontrado');
      return result;
    }

    if (dryRun) {
      logger.info(`[DRY-RUN] Removería ${analysis.items.length} arquivos (${logger.formatBytes(analysis.totalSize)})`);
      return result;
    }

    logger.warn('Use o modo interativo para selecionar arquivos específicos para remoção');

    return result;
  }
}

export const largeFilesModule = new LargeFilesModule();

export class DiskUsageModule {
  readonly id = 'disk-usage';
  readonly name = 'Uso de Disco';
  readonly description = 'Analisa uso de disco similar ao ncdu';

  isAvailable(): boolean {
    return existsSync('/usr/bin/du');
  }

  async analyze(directory: string = '/home', depth: number = 2): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const duResult = await exec('du', ['-ah', `--max-depth=${depth}`, directory]);

    if (duResult.success && duResult.stdout) {
      const lines = duResult.stdout.split('\n').filter(l => l.trim());
      
      for (const line of lines) {
        const parts = line.split('\t');
        if (parts.length >= 2) {
          const sizeStr = parts[0];
          const path = parts[1];
          
          const size = this.parseSize(sizeStr);
          items.push({
            path,
            size,
            type: 'directory',
            description: `Diretório: ${path}`,
          });
          totalSize += size;
        }
      }
    }

    items.sort((a, b) => b.size - a.size);

    return { module: this.id, items: items.slice(0, 50), totalSize };
  }

  private parseSize(sizeStr: string): number {
    const match = sizeStr.match(/([\d.]+)\s*([A-Z])/i);
    if (!match) return 0;
    
    const num = parseFloat(match[1]);
    const unit = match[2].toUpperCase();
    
    const multipliers: Record<string, number> = {
      'B': 1,
      'K': 1024,
      'M': 1024 * 1024,
      'G': 1024 * 1024 * 1024,
      'T': 1024 * 1024 * 1024 * 1024,
    };
    
    return num * (multipliers[unit] || 1);
  }

  async clean(_dryRun: boolean = false, _force: boolean = false): Promise<CleaningResult> {
    const analysis = await this.analyze();
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: analysis.totalSize,
      itemsRemoved: analysis.items.length,
      errors: [],
    };

    logger.info('Análise de uso de disco concluída');
    logger.info(`Total: ${logger.formatBytes(analysis.totalSize)} em ${analysis.items.length} diretórios`);

    return result;
  }
}

export const diskUsageModule = new DiskUsageModule();
