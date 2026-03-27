import type { CliFlags } from '../types/index.js';
import { getDefaultFlags, DEFAULT_THRESHOLD, MIN_THRESHOLD, MAX_THRESHOLD } from './flags.js';

export function parseFlags(args: string[]): CliFlags {
  const flags = getDefaultFlags();

  flags.all = args.includes('--all') || args.includes('-a');
  flags.cache = args.includes('--cache');
  flags.npm = args.includes('--npm');
  flags.yarn = args.includes('--yarn');
  flags.pnpm = args.includes('--pnpm');
  flags.flatpak = args.includes('--flatpak');
  flags.snap = args.includes('--snap');
  flags.docker = args.includes('--docker');
  flags.logs = args.includes('--logs');
  flags.packages = args.includes('--packages');
  flags.analyze = args.includes('--analyze');
  flags.dryRun = args.includes('--dry-run') || args.includes('-n');
  flags.force = args.includes('--force') || args.includes('-f');
  flags.interactive = args.includes('--interactive') || args.includes('-i');
  flags.largeFiles = args.includes('--large-files');
  flags.appimage = args.includes('--appimage');
  flags.thumbs = args.includes('--thumbs');
  flags.recent = args.includes('--recent');

  const thresholdArg = args.find(a => a.startsWith('--threshold='));
  if (thresholdArg) {
    const val = thresholdArg.split('=')[1];
    const parsed = val ? parseInt(val) : NaN;
    flags.largeFilesThreshold = isNaN(parsed) ? DEFAULT_THRESHOLD : clampThreshold(parsed);
  }

  return flags;
}

function clampThreshold(value: number): number {
  if (value < MIN_THRESHOLD) return MIN_THRESHOLD;
  if (value > MAX_THRESHOLD) return MAX_THRESHOLD;
  return value;
}

export function getModulesFromFlags(flags: CliFlags): string[] {
  const modules: string[] = [];

  if (flags.all) {
    return [
      'packages',
      'cache',
      'npm',
      'yarn',
      'pnpm',
      'flatpak',
      'snap',
      'docker',
      'logs',
      'large-files',
      'appimage',
      'thumbs',
      'recent',
    ];
  }

  if (flags.cache) modules.push('cache');
  if (flags.npm) modules.push('npm');
  if (flags.yarn) modules.push('yarn');
  if (flags.pnpm) modules.push('pnpm');
  if (flags.flatpak) modules.push('flatpak');
  if (flags.snap) modules.push('snap');
  if (flags.docker) modules.push('docker');
  if (flags.logs) modules.push('logs');
  if (flags.packages) modules.push('packages');
  if (flags.largeFiles) modules.push('large-files');
  if (flags.appimage) modules.push('appimage');
  if (flags.thumbs) modules.push('thumbs');
  if (flags.recent) modules.push('recent');

  return modules;
}
