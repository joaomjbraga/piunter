export type PackageManager = 'apt' | 'pacman' | 'dnf' | 'unknown';

export interface DistroInfo {
  id: string;
  name: string;
  version: string;
  packageManager: PackageManager;
}

export interface CommandResult {
  success: boolean;
  stdout: string;
  stderr: string;
  code: number;
}

export interface CleanableItem {
  path: string;
  size: number;
  type: string;
  description: string;
}

export interface CleaningResult {
  module: string;
  success: boolean;
  spaceFreed: number;
  itemsRemoved: number;
  errors: string[];
}

export interface AnalysisResult {
  module: string;
  items: CleanableItem[];
  totalSize: number;
}

export interface ModuleInfo {
  id: string;
  name: string;
  description: string;
  available: boolean;
}

export interface CleanOptions {
  dryRun: boolean;
  force: boolean;
  modules: string[];
}

export interface CliFlags {
  all: boolean;
  cache: boolean;
  npm: boolean;
  yarn: boolean;
  pnpm: boolean;
  flatpak: boolean;
  snap: boolean;
  docker: boolean;
  logs: boolean;
  packages: boolean;
  analyze: boolean;
  dryRun: boolean;
  force: boolean;
  interactive: boolean;
  largeFiles: boolean;
  largeFilesThreshold: number;
  appimage: boolean;
  thumbs: boolean;
  recent: boolean;
}

export interface Report {
  startTime: Date;
  endTime: Date;
  modules: CleaningResult[];
  totalSpaceFreed: number;
  totalItemsRemoved: number;
  errors: string[];
}
