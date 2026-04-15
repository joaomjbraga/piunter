import { existsSync, readFileSync, writeFileSync, statSync } from 'fs';
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

export const VERSION = '1.2.3';

const DEFAULT_CONFIG: PiunterConfig = {
  version: VERSION,
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

let cachedConfig: PiunterConfig | null = null;
let cachedConfigMtime: number | null = null;
let configLoadPromise: Promise<PiunterConfig> | null = null;

export function loadConfig(): PiunterConfig {
  const configPath = getConfigPath();

  try {
    if (existsSync(configPath)) {
      const stat = statSync(configPath);
      const currentMtime = stat.mtimeMs;

      if (cachedConfig && cachedConfigMtime === currentMtime) {
        return cachedConfig;
      }

      const content = readFileSync(configPath, 'utf-8');
      const config = JSON.parse(content);
      cachedConfig = { ...DEFAULT_CONFIG, ...config } as PiunterConfig;
      cachedConfigMtime = currentMtime;
      return cachedConfig;
    }
  } catch {
    cachedConfig = DEFAULT_CONFIG;
    cachedConfigMtime = null;
  }

  if (!cachedConfig) {
    cachedConfig = DEFAULT_CONFIG;
    cachedConfigMtime = null;
  }
  return cachedConfig;
}

export async function loadConfigAsync(): Promise<PiunterConfig> {
  if (configLoadPromise) {
    return configLoadPromise;
  }

  configLoadPromise = Promise.resolve().then(() => {
    const result = loadConfig();
    configLoadPromise = null;
    return result;
  });

  return configLoadPromise;
}

export function clearConfigCache(): void {
  cachedConfig = null;
  cachedConfigMtime = null;
  configLoadPromise = null;
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
