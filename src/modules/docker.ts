import { exec, isCommandAvailable } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';

export class DockerModule {
  readonly id = 'docker';
  readonly name = 'Docker';
  readonly description = 'Remove imagens, containers e volumes Docker não utilizados';

  isAvailable(): boolean {
    return isCommandAvailable('docker');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const imagesResult = await exec('docker', ['images', '--format', '{{.Size}}\t{{.Repository}}:{{.Tag}}']);
    if (imagesResult.success) {
      const lines = imagesResult.stdout.split('\n').filter(l => l.trim());
      for (const line of lines) {
        const parts = line.split('\t');
        if (parts.length >= 2) {
          const size = this.parseSize(parts[0]);
          const name = parts[1];
          if (name !== '<none>:<none>') {
            items.push({
              path: name,
              size,
              type: 'docker-image',
              description: `Imagem Docker: ${name}`,
            });
            totalSize += size;
          }
        }
      }
    }

    const containersResult = await exec('docker', ['ps', '-a', '--format', '{{.Size}}\t{{.Names}}']);
    if (containersResult.success) {
      const lines = containersResult.stdout.split('\n').filter(l => l.trim());
      for (const line of lines) {
        const parts = line.split('\t');
        if (parts.length >= 2) {
          const size = this.parseSize(parts[0]);
          const name = parts[1];
          items.push({
            path: name,
            size,
            type: 'docker-container',
            description: `Container Docker: ${name}`,
          });
          totalSize += size;
        }
      }
    }

    const volumesResult = await exec('docker', ['volume', 'ls', '--format', '{{.Name}}']);
    if (volumesResult.success) {
      const lines = volumesResult.stdout.split('\n').filter(l => l.trim());
      for (const line of lines) {
        items.push({
          path: line,
          size: 0,
          type: 'docker-volume',
          description: `Volume Docker: ${line}`,
        });
      }
    }

    const systemDfResult = await exec('docker', ['system', 'df']);
    if (systemDfResult.success) {
      const match = systemDfResult.stdout.match(/Total\s+([\d.]+\s*[A-Z]+)/i);
      if (match) {
        const size = this.parseSize(match[1]);
        items.push({
          path: 'docker-system',
          size,
          type: 'docker-total',
          description: 'Total de recursos Docker',
        });
        totalSize += size;
      }
    }

    return { module: this.id, items, totalSize };
  }

  private parseSize(sizeStr: string): number {
    const match = sizeStr.match(/([\d.]+)\s*([A-Z]+)?B?/i);
    if (!match) return 0;
    
    const num = parseFloat(match[1]);
    const unit = (match[2] || 'MB').toUpperCase();
    
    const multipliers: Record<string, number> = {
      'B': 1,
      'KB': 1024,
      'MB': 1024 * 1024,
      'GB': 1024 * 1024 * 1024,
    };
    
    return num * (multipliers[unit] || 1024 * 1024);
  }

  async clean(dryRun: boolean = false, force: boolean = false): Promise<CleaningResult> {
    const analysis = await this.analyze();
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    if (dryRun) {
      logger.info(`[DRY-RUN] Docker: limparía ${logger.formatBytes(analysis.totalSize)}`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const containersResult = await exec('docker', ['container', 'prune', '-f']);
    if (containersResult.success) {
      logger.item(`${this.name}: Containers parados removidos`);
    } else {
      result.errors.push('Falha ao limpar containers');
    }

    const networksResult = await exec('docker', ['network', 'prune', '-f']);
    if (networksResult.success) {
      logger.item(`${this.name}: Networks não utilizadas removidas`);
    }

    const imagesResult = await exec('docker', ['image', 'prune', '-a', '-f']);
    if (imagesResult.success) {
      logger.item(`${this.name}: Imagens não utilizadas removidas`);
    }

    const volumesResult = await exec('docker', ['volume', 'prune', '-f']);
    if (volumesResult.success) {
      logger.item(`${this.name}: Volumes não utilizados removidos`);
    }

    const systemPruneResult = await exec('docker', ['system', 'prune', '-a', '-f', '--volumes']);
    if (systemPruneResult.success) {
      logger.item(`${this.name}: Sistema Docker completo otimizado`);
      result.success = true;
      result.spaceFreed = analysis.totalSize;
      result.itemsRemoved = analysis.items.length;
    } else {
      result.errors.push('Falha na limpeza completa do Docker');
    }

    return result;
  }
}

export const dockerModule = new DockerModule();
