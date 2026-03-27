import { describe, it, expect } from 'vitest';
import { logger } from '../utils/logger.js';

describe('Logger', () => {
  it('should format bytes correctly', () => {
    expect(logger.formatBytes(0)).toBe('0 B');
    expect(logger.formatBytes(1024)).toBe('1 KB');
    expect(logger.formatBytes(1024 * 1024)).toBe('1 MB');
    expect(logger.formatBytes(1024 * 1024 * 1024)).toBe('1 GB');
  });
});
