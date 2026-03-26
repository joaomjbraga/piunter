import { readFileSync, existsSync } from 'fs';
import type { DistroInfo, PackageManager } from '../types/index.js';

export function getDistroInfo(): DistroInfo {
  const osReleasePath = '/etc/os-release';

  if (!existsSync(osReleasePath)) {
    return {
      id: 'unknown',
      name: 'Unknown',
      version: 'unknown',
      packageManager: 'unknown',
    };
  }

  const content = readFileSync(osReleasePath, 'utf-8');
  const lines = content.split('\n');
  const info: Record<string, string> = {};

  for (const line of lines) {
    const [key, ...valueParts] = line.split('=');
    if (key && valueParts.length > 0) {
      info[key] = valueParts.join('=').replace(/^"|"$/g, '');
    }
  }

  const id = info.ID || 'unknown';
  const name = info.NAME || info.PRETTY_NAME || 'Unknown';
  const version = info.VERSION_ID || info.VERSION || 'unknown';
  const packageManager = detectPackageManager(id);

  return { id, name, version, packageManager };
}

function detectPackageManager(distroId: string): PackageManager {
  const debianBased = ['debian', 'ubuntu', 'pop', 'linuxmint', 'elementary', 'zorin'];
  const archBased = ['arch', 'manjaro', 'endeavouros', 'garuda', 'arcolinux'];
  const fedoraBased = ['fedora', 'rhel', 'centos', 'rocky', 'almalinux'];

  if (debianBased.some(d => distroId.toLowerCase().includes(d))) {
    return 'apt';
  }
  if (archBased.some(d => distroId.toLowerCase().includes(d))) {
    return 'pacman';
  }
  if (fedoraBased.some(d => distroId.toLowerCase().includes(d))) {
    return 'dnf';
  }

  if (existsSync('/usr/bin/apt')) return 'apt';
  if (existsSync('/usr/bin/pacman')) return 'pacman';
  if (existsSync('/usr/bin/dnf')) return 'dnf';

  return 'unknown';
}

export function getHomeDir(): string {
  return process.env.HOME || '/root';
}

export function getUsername(): string {
  return process.env.USER || 'root';
}

export function getCacheDir(): string {
  const xdgCache = process.env.XDG_CACHE_HOME;
  if (xdgCache) return xdgCache;
  return `${getHomeDir()}/.cache`;
}

export function getConfigDir(): string {
  const xdgConfig = process.env.XDG_CONFIG_HOME;
  if (xdgConfig) return xdgConfig;
  return `${getHomeDir()}/.config`;
}
