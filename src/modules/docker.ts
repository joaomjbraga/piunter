import { exec, isCommandAvailable } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';
import { parseSize } from '../utils/fs.js';

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

    const imagesResult = await exec('docker', [
      'images',
      '--format',
      '{{.Size}}\t{{.Repository}}:{{.Tag}}',
    ]);
    if (imagesResult.success) {
      const lines = imagesResult.stdout.split('\n').filter(l => l.trim());
      for (const line of lines) {
        const parts = line.split('\t');
        if (parts.length >= 2) {
          const size = parseSize(parts[0]);
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

    const containersResult = await exec('docker', [
      'ps',
      '-a',
      '--format',
      '{{.Size}}\t{{.Names}}',
    ]);
    if (containersResult.success) {
      const lines = containersResult.stdout.split('\n').filter(l => l.trim());
      for (const line of lines) {
        const parts = line.split('\t');
        if (parts.length >= 2) {
          const size = parseSize(parts[0]);
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
        const size = parseSize(match[1]);
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

  async clean(dryRun: boolean = false): Promise<CleaningResult> {
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    if (!this.isAvailable()) {
      result.errors.push('Docker não está instalado ou o serviço não está em execução');
      result.success = false;
      return result;
    }

    const beforeAnalysis = await this.analyze();
    const beforeSize = beforeAnalysis.totalSize;

    if (dryRun) {
      logger.info(`[DRY-RUN] Docker: limparía ${logger.formatBytes(beforeAnalysis.totalSize)}`);
      result.spaceFreed = beforeAnalysis.totalSize;
      return result;
    }

    try {
      const containersResult = await exec('docker', ['container', 'prune', '-f']);
      if (containersResult.success) {
        logger.item(`${this.name}: Containers parados removidos`);
      } else if (containersResult.stderr) {
        result.errors.push(`Containers: ${containersResult.stderr}`);
      }
    } catch (e) {
      result.errors.push(`Falha ao limpar containers: ${(e as Error).message}`);
    }

    try {
      const networksResult = await exec('docker', ['network', 'prune', '-f']);
      if (networksResult.success) {
        logger.item(`${this.name}: Networks não utilizadas removidas`);
      }
    } catch (e) {
      result.errors.push(`Falha ao limpar networks: ${(e as Error).message}`);
    }

    try {
      const imagesResult = await exec('docker', ['image', 'prune', '-a', '-f']);
      if (imagesResult.success) {
        logger.item(`${this.name}: Imagens não utilizadas removidas`);
      } else if (imagesResult.stderr) {
        result.errors.push(`Imagens: ${imagesResult.stderr}`);
      }
    } catch (e) {
      result.errors.push(`Falha ao limpar imagens: ${(e as Error).message}`);
    }

    try {
      const volumesResult = await exec('docker', ['volume', 'prune', '-f']);
      if (volumesResult.success) {
        logger.item(`${this.name}: Volumes não utilizados removidos`);
      } else if (volumesResult.stderr) {
        result.errors.push(`Volumes: ${volumesResult.stderr}`);
      }
    } catch (e) {
      result.errors.push(`Falha ao limpar volumes: ${(e as Error).message}`);
    }

    try {
      const systemPruneResult = await exec('docker', ['system', 'prune', '-a', '-f']);
      if (systemPruneResult.success) {
        logger.item(`${this.name}: Sistema Docker completo otimizado`);
        result.success = true;
      } else if (systemPruneResult.stderr) {
        result.errors.push(`System prune: ${systemPruneResult.stderr}`);
      }
    } catch (e) {
      result.errors.push(`Falha na limpeza completa do Docker: ${(e as Error).message}`);
    }

    const afterAnalysis = await this.analyze();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);
    result.itemsRemoved = Math.max(0, beforeAnalysis.items.length - afterAnalysis.items.length);

    return result;
  }
}

export const dockerModule = new DockerModule();
