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

    const [imagesResult, containersResult, volumesResult] = await Promise.all([
      exec('docker', ['images', '--format', '{{.Size}}\t{{.Repository}}:{{.Tag}}']),
      exec('docker', ['ps', '-a', '--format', '{{.Size}}\t{{.Names}}']),
      exec('docker', ['volume', 'ls', '--format', '{{.Name}}']),
    ]);

    let totalSize = 0;

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

    const pruneCommands = [
      { cmd: ['container', 'prune', '-f'], name: 'Containers', key: 'containers' },
      { cmd: ['network', 'prune', '-f'], name: 'Networks', key: 'networks' },
      { cmd: ['image', 'prune', '-a', '-f'], name: 'Imagens', key: 'images' },
      { cmd: ['volume', 'prune', '-f'], name: 'Volumes', key: 'volumes' },
    ];

    const pruneResults = await Promise.all(
      pruneCommands.map(async ({ cmd, name, key }) => {
        try {
          const cmdResult = await exec('docker', cmd);
          if (cmdResult.success) {
            logger.item(`${this.name}: ${name} não utilizados removidos`);
            return { success: true, key, count: 1 };
          } else {
            logger.debug(`${this.name}: Falha ao limpar ${name}: ${cmdResult.stderr}`);
            return { success: false, key, error: cmdResult.stderr, count: 0 };
          }
        } catch (e) {
          return { success: false, key, error: (e as Error).message, count: 0 };
        }
      })
    );

    result.itemsRemoved = pruneResults.reduce((sum, r) => sum + r.count, 0);

    const failures = pruneResults.filter(r => !r.success);
    if (failures.length > 0) {
      for (const pruneResult of failures) {
        result.errors.push(`${pruneResult.key}: ${pruneResult.error}`);
      }
      if (failures.length === pruneResults.length) {
        result.success = false;
      }
    }

    const afterAnalysis = await this.analyze();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);
    result.itemsRemoved = Math.max(0, beforeAnalysis.items.length - afterAnalysis.items.length);

    return result;
  }
}

export const dockerModule = new DockerModule();
