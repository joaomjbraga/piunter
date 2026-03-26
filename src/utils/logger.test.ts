import { describe, it, expect } from 'vitest';
import { getDistroInfo } from '../utils/os.js';
import { logger } from '../utils/logger.js';

describe('OS Utils', () => {
  it('should get distro info', () => {
    const distro = getDistroInfo();
    expect(distro).toBeDefined();
    expect(distro.id).toBeDefined();
    expect(distro.name).toBeDefined();
    expect(distro.packageManager).toBeDefined();
  });
});

describe('Logger', () => {
  it('should format bytes correctly', () => {
    expect(logger.formatBytes(0)).toBe('0 B');
    expect(logger.formatBytes(1024)).toBe('1 KB');
    expect(logger.formatBytes(1024 * 1024)).toBe('1 MB');
    expect(logger.formatBytes(1024 * 1024 * 1024)).toBe('1 GB');
  });
});
