import { existsSync } from 'fs';
import { exec, isCommandAvailable } from '../utils/exec.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';
import { logger } from '../utils/logger.js';
import { loadConfig } from '../utils/config.js';

export class LogsModule {
  readonly id = 'logs';
  readonly name = 'Logs do Sistema';
  readonly description = 'Limpa logs do systemd journal e logs antigos';

  isAvailable(): boolean {
    return existsSync('/var/log') && isCommandAvailable('journalctl');
  }

  async analyze(): Promise<AnalysisResult> {
    const items: AnalysisResult['items'] = [];

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
      }
    }

    const logDirs = ['/var/log', '/var/log.old'];

    const dirResults = await Promise.all(
      logDirs.map(async logDir => {
        if (!existsSync(logDir)) return { files: [], sizes: [] };

        const findResult = await exec('find', [
          logDir,
          '-type',
          'f',
          '-name',
          '*.log',
          '-size',
          '+1M',
        ]);

        if (!findResult.success) return { files: [], sizes: [] };

        const files = findResult.stdout.split('\n').filter(l => l.trim());
        if (files.length === 0) return { files: [], sizes: [] };

        const sizeResults = await Promise.all(
          files.map(async file => {
            const sizeResult = await exec('du', ['-sb', file]);
            if (sizeResult.success) {
              const match = sizeResult.stdout.match(/^(\d+)/);
              return match ? parseInt(match[1], 10) : 0;
            }
            return 0;
          })
        );

        return { files, sizes: sizeResults };
      })
    );

    let totalSize = 0;
    for (const result of dirResults) {
      for (let i = 0; i < result.files.length; i++) {
        const file = result.files[i];
        const size = result.sizes[i];
        items.push({
          path: file,
          size,
          type: 'log-file',
          description: `Arquivo de log: ${file}`,
        });
        totalSize += size;
      }
    }

    return { module: this.id, items, totalSize };
  }

  async clean(dryRun: boolean = false): Promise<CleaningResult> {
    const config = loadConfig();
    const journalSizeMB = config.thresholds.journalSizeMB;
    const logDays = config.thresholds.logDays;

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

    const [vacuumResult, vacuumTimeResult, oldLogsResult] = await Promise.all([
      exec('journalctl', [`--vacuum-size=${journalSizeMB}M`], { sudo: true }),
      exec('journalctl', [`--vacuum-time=${logDays}d`], { sudo: true }),
      exec(
        'find',
        ['/var/log', '-type', 'f', '-name', '*.log', '-mtime', `+${logDays}`, '-delete'],
        { sudo: true }
      ),
    ]);

    if (vacuumResult.success) {
      logger.item(`${this.name}: Journal limpo (limite ${journalSizeMB}MB)`);
      result.itemsRemoved++;
    } else {
      result.errors.push('Falha ao limpar journalctl (verifique se tem privilégios sudo)');
    }

    if (vacuumTimeResult.success) {
      logger.item(`${this.name}: Logs anteriores a ${logDays} dias removidos`);
    }

    if (oldLogsResult.success) {
      logger.item(`${this.name}: Logs antigos (>${logDays} dias) removidos`);
    }

    const afterAnalysis = await this.analyze();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);

    return result;
  }

  async cleanOldLogs(days?: number, dryRun: boolean = false): Promise<CleaningResult> {
    const config = loadConfig();
    const logDays = days ?? config.thresholds.logDays;

    const beforeAnalysis = await this.analyze();
    const beforeSize = beforeAnalysis.totalSize;
    const result: CleaningResult = {
      module: this.id,
      success: true,
      spaceFreed: 0,
      itemsRemoved: 0,
      errors: [],
    };

    if (dryRun) {
      const findResult = await exec('find', [
        '/var/log',
        '-type',
        'f',
        '-name',
        '*.log',
        '-mtime',
        `+${logDays}`,
      ]);
      if (findResult.success) {
        const files = findResult.stdout.split('\n').filter(l => l.trim());
        logger.info(`[DRY-RUN] Removería ${files.length} logs com mais de ${logDays} dias`);
      }
      result.spaceFreed = beforeSize;
      return result;
    }

    const cleanResult = await exec('find', [
      '/var/log',
      '-type',
      'f',
      '-name',
      '*.log',
      '-mtime',
      `+${logDays}`,
      '-delete',
    ]);
    if (cleanResult.success) {
      logger.item(`${this.name}: Logs com mais de ${logDays} dias removidos`);
      result.success = true;
    } else {
      result.errors.push('Falha ao limpar logs antigos');
    }

    const afterAnalysis = await this.analyze();
    result.spaceFreed = Math.max(0, beforeSize - afterAnalysis.totalSize);

    return result;
  }
}

export const logsModule = new LogsModule();
