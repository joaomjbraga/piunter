import { existsSync } from 'fs';
import { exec, isCommandAvailable } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';

export class LogsModule {
  readonly id = 'logs';
  readonly name = 'Logs do Sistema';
  readonly description = 'Limpa logs do systemd journal e logs antigos';

  isAvailable(): boolean {
    return existsSync('/var/log') && isCommandAvailable('journalctl');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];
    let totalSize = 0;

    const journalResult = await exec('du', ['-sb', '/var/log/journal']);
    if (journalResult.success) {
      const match = journalResult.stdout.match(/^(\d+)/);
      if (match) {
        const size = parseInt(match[1], 10);
        items.push({
          path: '/var/log/journal',
          size,
          type: 'journal',
          description: 'Logs do systemd journal',
        });
        totalSize += size;
      }
    }

    const logDirs = ['/var/log', '/var/log.old'];
    for (const logDir of logDirs) {
      if (existsSync(logDir)) {
        const findResult = await exec('find', [logDir, '-type', 'f', '-name', '*.log', '-size', '+1M']);
        if (findResult.success) {
          const files = findResult.stdout.split('\n').filter(l => l.trim());
          for (const file of files) {
            const sizeResult = await exec('du', ['-sb', file]);
            if (sizeResult.success) {
              const match = sizeResult.stdout.match(/^(\d+)/);
              if (match) {
                const size = parseInt(match[1], 10);
                items.push({
                  path: file,
                  size,
                  type: 'log-file',
                  description: `Arquivo de log: ${file}`,
                });
                totalSize += size;
              }
            }
          }
        }
      }
    }

    return { module: this.id, items, totalSize };
  }

  async clean(dryRun: boolean = false, _force: boolean = false): Promise<CleaningResult> {
    const analysis = await this.analyze();
    const beforeSize = analysis.totalSize;
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    if (dryRun) {
      logger.info(`[DRY-RUN] Logs: limparía ${logger.formatBytes(analysis.totalSize)}`);
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const vacuumResult = await exec('journalctl', ['--vacuum-size=500M'], { sudo: true });
    if (vacuumResult.success) {
      logger.item(`${this.name}: Journal limpo (limite 500MB)`);
      result.itemsRemoved++;
    } else {
      result.errors.push('Falha ao limpar journalctl (verifique se tem privilégios sudo)');
    }

    const vacuumTimeResult = await exec('journalctl', ['--vacuum-time=7d'], { sudo: true });
    if (vacuumTimeResult.success) {
      logger.item(`${this.name}: Logs anteriores a 7 dias removidos`);
    }

    const oldLogsResult = await exec('find', ['/var/log', '-type', 'f', '-name', '*.log', '-mtime', '+30', '-delete'], { sudo: true });
    if (oldLogsResult.success) {
      logger.item(`${this.name}: Logs antigos (>30 dias) removidos`);
    }

    const afterAnalysis = await this.analyze();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);

    return result;
  }

  async cleanOldLogs(days: number = 30, dryRun: boolean = false): Promise<CleaningResult> {
    const analysis = await this.analyze();
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    if (dryRun) {
      const findResult = await exec('find', ['/var/log', '-type', 'f', '-name', '*.log', '-mtime', `+${days}`]);
      if (findResult.success) {
        const files = findResult.stdout.split('\n').filter(l => l.trim());
        logger.info(`[DRY-RUN] Removería ${files.length} logs com mais de ${days} dias`);
      }
      result.spaceFreed = analysis.totalSize;
      return result;
    }

    const cleanResult = await exec('find', ['/var/log', '-type', 'f', '-name', '*.log', '-mtime', `+${days}`, '-delete']);
    if (cleanResult.success) {
      logger.item(`${this.name}: Logs com mais de ${days} dias removidos`);
      result.success = true;
      result.spaceFreed = analysis.totalSize;
    } else {
      result.errors.push('Falha ao limpar logs antigos');
    }

    return result;
  }
}

export const logsModule = new LogsModule();
