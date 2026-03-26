import { existsSync, readFileSync, writeFileSync } from 'fs';
import { join } from 'path';
import { getHomeDir } from './os.js';

export interface PiunterConfig {
  version: string;
  defaults: {
    dryRun: boolean;
    force: boolean;
    modules: string[];
  };
  thresholds: {
    largeFilesMB: number;
    logDays: number;
    journalSizeMB: number;
  };
  sudo: {
    autoPrompt: boolean;
  };
}

const DEFAULT_CONFIG: PiunterConfig = {
  version: '1.0.0',
  defaults: {
    dryRun: false,
    force: false,
    modules: ['packages', 'cache', 'npm'],
  },
  thresholds: {
    largeFilesMB: 100,
    logDays: 30,
    journalSizeMB: 500,
  },
  sudo: {
    autoPrompt: true,
  },
};

export function getConfigPath(): string {
  return join(getHomeDir(), '.piunter.json');
}

export function loadConfig(): PiunterConfig {
  const configPath = getConfigPath();

  if (!existsSync(configPath)) {
    return DEFAULT_CONFIG;
  }

  try {
    const content = readFileSync(configPath, 'utf-8');
    const config = JSON.parse(content);
    return { ...DEFAULT_CONFIG, ...config };
  } catch {
    return DEFAULT_CONFIG;
  }
}

export function saveConfig(config: PiunterConfig): void {
  const configPath = getConfigPath();
  writeFileSync(configPath, JSON.stringify(config, null, 2), 'utf-8');
}

export function createDefaultConfig(): void {
  const configPath = getConfigPath();

  if (!existsSync(configPath)) {
    saveConfig(DEFAULT_CONFIG);
  }
}

export function validateThreshold(value: number, min: number, max: number): number {
  if (isNaN(value) || value < min) return min;
  if (value > max) return max;
  return value;
}
