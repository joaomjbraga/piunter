export * from './cache.js';
export * from './npm.js';
export * from './flatpak.js';
export * from './docker.js';
export * from './logs.js';
export * from './packages.js';
export * from './disk.js';

import { cacheModule } from './cache.js';
import { npmModule, yarnModule, pnpmModule } from './npm.js';
import { flatpakModule } from './flatpak.js';
import { dockerModule } from './docker.js';
import { logsModule } from './logs.js';
import { packagesModule } from './packages.js';
import { largeFilesModule, diskUsageModule } from './disk.js';
import type { ModuleInfo } from '../types/index.js';

export interface Module {
  id: string;
  name: string;
  description: string;
  isAvailable(): boolean;
  analyze(): Promise<import('../types/index.js').AnalysisResult>;
  clean(dryRun?: boolean, force?: boolean): Promise<import('../types/index.js').CleaningResult>;
}

export const modules: Module[] = [
  packagesModule,
  cacheModule,
  npmModule,
  yarnModule,
  pnpmModule,
  flatpakModule,
  dockerModule,
  logsModule,
  largeFilesModule,
  diskUsageModule,
];

export function getAvailableModules(): ModuleInfo[] {
  return modules.map(m => ({
    id: m.id,
    name: m.name,
    description: m.description,
    available: m.isAvailable(),
  }));
}

export function getModule(id: string): Module | undefined {
  return modules.find(m => m.id === id);
}

export function getModuleByIds(ids: string[]): Module[] {
  return modules.filter(m => ids.includes(m.id));
}
