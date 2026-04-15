import type { CliFlags } from '../types/index.js';
import { getDefaultFlags, MIN_THRESHOLD, MAX_THRESHOLD } from './flags.js';

const THRESHOLD_REGEX = /^\d+$/;

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
    if (val && THRESHOLD_REGEX.test(val)) {
      const parsed = parseInt(val, 10);
      flags.largeFilesThreshold = clampThreshold(parsed);
    }
  }

  return flags;
}

function clampThreshold(value: number): number {
  if (value < MIN_THRESHOLD) return MIN_THRESHOLD;
  if (value > MAX_THRESHOLD) return MAX_THRESHOLD;
  return value;
}

export function getModulesFromFlags(flags: CliFlags): string[] {
  const moduleFlags = [
    { flag: flags.cache, id: 'cache' },
    { flag: flags.npm, id: 'npm' },
    { flag: flags.yarn, id: 'yarn' },
    { flag: flags.pnpm, id: 'pnpm' },
    { flag: flags.flatpak, id: 'flatpak' },
    { flag: flags.snap, id: 'snap' },
    { flag: flags.docker, id: 'docker' },
    { flag: flags.logs, id: 'logs' },
    { flag: flags.packages, id: 'packages' },
    { flag: flags.largeFiles, id: 'large-files' },
    { flag: flags.appimage, id: 'appimage' },
    { flag: flags.thumbs, id: 'thumbs' },
    { flag: flags.recent, id: 'recent' },
  ];

  if (flags.all) {
    return moduleFlags.map(m => m.id);
  }

  return moduleFlags.filter(m => m.flag).map(m => m.id);
}
