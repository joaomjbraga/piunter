export * from './cache.js';
export * from './npm.js';
export * from './flatpak.js';
export * from './docker.js';
export * from './logs.js';
export * from './packages.js';
export * from './disk.js';
export * from './snap.js';
export * from './appimage.js';
export * from './thumbs.js';
export * from './recent.js';

import { cacheModule } from './cache.js';
import { npmModule, yarnModule, pnpmModule } from './npm.js';
import { flatpakModule } from './flatpak.js';
import { dockerModule } from './docker.js';
import { logsModule } from './logs.js';
import { packagesModule } from './packages.js';
import { largeFilesModule } from './disk.js';
import { snapModule } from './snap.js';
import { appimageModule } from './appimage.js';
import { thumbsModule } from './thumbs.js';
import { recentFilesModule } from './recent.js';
import type { ModuleInfo } from '../types/index.js';

export interface Module {
  id: string;
  name: string;
  description: string;
  isAvailable(): boolean;
  analyze(threshold?: number): Promise<import('../types/index.js').AnalysisResult>;
  clean(dryRun?: boolean): Promise<import('../types/index.js').CleaningResult>;
}

export const modules: Module[] = [
  packagesModule,
  cacheModule,
  npmModule,
  yarnModule,
  pnpmModule,
  flatpakModule,
  snapModule,
  dockerModule,
  logsModule,
  largeFilesModule,
  appimageModule,
  thumbsModule,
  recentFilesModule,
];

let cachedModules: ModuleInfo[] | null = null;
let moduleCacheTime = 0;
let cachePromise: Promise<ModuleInfo[]> | null = null;
const MODULE_CACHE_TTL = 5000;

export function getAvailableModules(): ModuleInfo[] {
  const now = Date.now();
  if (cachedModules && now - moduleCacheTime < MODULE_CACHE_TTL) {
    return cachedModules;
  }

  cachedModules = modules.map(m => ({
    id: m.id,
    name: m.name,
    description: m.description,
    available: m.isAvailable(),
  }));
  moduleCacheTime = now;
  return cachedModules;
}

export function getAvailableModulesAsync(): Promise<ModuleInfo[]> {
  if (cachePromise) {
    return cachePromise;
  }

  cachePromise = Promise.resolve().then(() => {
    const now = Date.now();
    if (cachedModules && now - moduleCacheTime < MODULE_CACHE_TTL) {
      return cachedModules;
    }

    cachedModules = modules.map(m => ({
      id: m.id,
      name: m.name,
      description: m.description,
      available: m.isAvailable(),
    }));
    moduleCacheTime = now;
    cachePromise = null;
    return cachedModules;
  });

  return cachePromise;
}

export function clearModuleCache(): void {
  cachedModules = null;
  moduleCacheTime = 0;
  cachePromise = null;
}

const moduleMap = new Map(modules.map(m => [m.id, m]));

export function getModule(id: string): Module | undefined {
  return moduleMap.get(id);
}

export function getModuleByIds(ids: string[]): Module[] {
  return ids.map(id => moduleMap.get(id)).filter((m): m is Module => m !== undefined);
}
