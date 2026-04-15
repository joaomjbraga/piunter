import { readdirSync, statSync } from 'fs';
import { readdir, stat } from 'fs/promises';
import { join } from 'path';

const SIZE_REGEX = /^([\d.]+)\s*(B|KB|MB|GB|TB|KiB|MiB|GiB|TiB)?$/i;
const SIZE_MULTIPLIERS: Record<string, number> = {
  B: 1,
  KB: 1024,
  KIB: 1024,
  MB: 1024 * 1024,
  MIB: 1024 * 1024,
  GB: 1024 * 1024 * 1024,
  GIB: 1024 * 1024 * 1024,
  TB: 1024 * 1024 * 1024 * 1024,
  TIB: 1024 * 1024 * 1024 * 1024,
};

export function getDirSize(dirPath: string): number {
  let size = 0;
  try {
    const entries = readdirSync(dirPath);
    for (const entry of entries) {
      const fullPath = join(dirPath, entry);
      try {
        const statInfo = statSync(fullPath);
        if (statInfo.isDirectory()) {
          size += getDirSize(fullPath);
        } else {
          size += statInfo.size;
        }
      } catch {
        // Skip inaccessible entries
      }
    }
  } catch {
    // Skip inaccessible directories
  }
  return size;
}

export async function getDirSizeAsync(dirPath: string): Promise<number> {
  let size = 0;
  try {
    const entries = await readdir(dirPath);
    const stats = await Promise.all(
      entries.map(async entry => {
        const fullPath = join(dirPath, entry);
        try {
          const statInfo = await stat(fullPath);
          return { isDir: statInfo.isDirectory(), size: statInfo.size, path: fullPath };
        } catch {
          return null;
        }
      })
    );

    const subdirs: string[] = [];
    for (const s of stats) {
      if (!s) continue;
      if (s.isDir) {
        subdirs.push(s.path);
      } else {
        size += s.size;
      }
    }

    const subSizes = await Promise.all(subdirs.map(sub => getDirSizeAsync(sub)));
    size += subSizes.reduce((sum, s) => sum + s, 0);
  } catch {
    // Skip inaccessible directories
  }
  return size;
}

export function parseSize(sizeStr: string): number {
  const match = sizeStr.match(SIZE_REGEX);
  if (!match) return 0;

  const num = parseFloat(match[1]);
  const unit = (match[2] || 'B').toUpperCase();

  return num * (SIZE_MULTIPLIERS[unit] || 1);
}
