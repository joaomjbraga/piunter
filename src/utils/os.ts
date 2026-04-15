import { readFileSync, existsSync, statSync } from 'fs';
import type { DistroInfo, PackageManager } from '../types/index.js';

let cachedDistroInfo: DistroInfo | null = null;
let cachedDistroMtime: number | null = null;
let distroPromise: Promise<DistroInfo> | null = null;

export function getDistroInfo(): DistroInfo {
  const osReleasePath = '/etc/os-release';

  if (existsSync(osReleasePath)) {
    try {
      const stat = statSync(osReleasePath);
      const currentMtime = stat.mtimeMs;

      if (cachedDistroInfo && cachedDistroMtime === currentMtime) {
        return cachedDistroInfo;
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

      cachedDistroInfo = { id, name, version, packageManager };
      cachedDistroMtime = currentMtime;
      return cachedDistroInfo;
    } catch {
      cachedDistroInfo = {
        id: 'unknown',
        name: 'Unknown',
        version: 'unknown',
        packageManager: 'unknown',
      };
      cachedDistroMtime = null;
      return cachedDistroInfo;
    }
  }

  cachedDistroInfo = {
    id: 'unknown',
    name: 'Unknown',
    version: 'unknown',
    packageManager: 'unknown',
  };
  cachedDistroMtime = null;
  return cachedDistroInfo;
}

export function getDistroInfoAsync(): Promise<DistroInfo> {
  if (distroPromise) {
    return distroPromise;
  }

  distroPromise = Promise.resolve().then(() => {
    const result = getDistroInfo();
    distroPromise = null;
    return result;
  });

  return distroPromise;
}

export function clearDistroCache(): void {
  cachedDistroInfo = null;
  cachedDistroMtime = null;
  distroPromise = null;
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
  const home = process.env.HOME || process.env.USERPROFILE;
  if (home && existsSync(home)) {
    return home;
  }
  if (existsSync('/root')) {
    return '/root';
  }
  const passwd = '/etc/passwd';
  if (existsSync(passwd)) {
    try {
      const content = readFileSync(passwd, 'utf-8');
      const lines = content.split('\n');
      for (const line of lines) {
        const parts = line.split(':');
        if (parts[0] === process.env.USER && parts[5]) {
          return parts[5];
        }
      }
    } catch {
      // ignore read errors, fallback below
    }
  }
  return '/tmp';
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
