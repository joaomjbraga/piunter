import type { CliFlags } from '../types/index.js';

export const DEFAULT_THRESHOLD = 100;
export const MIN_THRESHOLD = 1;
export const MAX_THRESHOLD = 10000;

export function getDefaultFlags(): CliFlags {
  return {
    all: false,
    cache: false,
    npm: false,
    yarn: false,
    pnpm: false,
    flatpak: false,
    snap: false,
    docker: false,
    logs: false,
    packages: false,
    analyze: false,
    dryRun: false,
    force: false,
    interactive: false,
    largeFiles: false,
    largeFilesThreshold: DEFAULT_THRESHOLD,
    appimage: false,
    thumbs: false,
    recent: false,
  };
}

export const MODULE_FLAG_MAP: Record<keyof CliFlags, string | null> = {
  all: null,
  cache: 'cache',
  npm: 'npm',
  yarn: 'yarn',
  pnpm: 'pnpm',
  flatpak: 'flatpak',
  snap: 'snap',
  docker: 'docker',
  logs: 'logs',
  packages: 'packages',
  analyze: null,
  dryRun: null,
  force: null,
  interactive: null,
  largeFiles: 'large-files',
  largeFilesThreshold: null,
  appimage: 'appimage',
  thumbs: 'thumbs',
  recent: 'recent',
};

export const MODULES_REQUIRING_SUDO = ['packages', 'logs', 'flatpak'];

export function requiresSudo(moduleIds: string[]): boolean {
  return moduleIds.some(id => MODULES_REQUIRING_SUDO.includes(id));
}
