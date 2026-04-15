import { describe, it, expect, vi, beforeEach } from 'vitest';
import { createCleaner } from './cleaner.js';
import type { Module } from '../modules/index.js';

vi.mock('../modules/index.js', () => ({
  getModuleByIds: vi.fn((ids: string[]) => {
    return ids.map(id => ({
      id,
      name: 'Test',
      description: 'Test',
      isAvailable: vi.fn(() => true),
      analyze: vi.fn().mockResolvedValue({
        module: id,
        items: [{ path: '/test', size: 1000, type: 'file', description: 'test' }],
        totalSize: 1000,
      }),
      clean: vi.fn().mockResolvedValue({
        module: id,
        success: true,
        spaceFreed: 500,
        itemsRemoved: 1,
        errors: [],
      }),
    }));
  }),
}));

describe('Cleaner', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should create cleaner with options', () => {
    const cleaner = createCleaner(['test'], { dryRun: true, force: false, modules: ['test'] });
    expect(cleaner).toBeDefined();
  });

  it('should clean modules and collect results', async () => {
    const cleaner = createCleaner(['test'], { dryRun: false, force: false, modules: ['test'] });
    const report = await cleaner.clean();

    expect(report.modules).toHaveLength(1);
    expect(report.modules[0].success).toBe(true);
    expect(report.totalSpaceFreed).toBe(500);
    expect(report.totalItemsRemoved).toBe(1);
  });

  it('should skip unavailable modules', async () => {
    const { getModuleByIds } = await import('../modules/index.js');
    const mockModule = {
      id: 'unavailable',
      name: 'Unavailable',
      description: 'Test',
      isAvailable: vi.fn(() => false),
      analyze: vi.fn(),
      clean: vi.fn(),
    };

    vi.mocked(getModuleByIds).mockReturnValue([mockModule as unknown as Module]);

    const cleaner = createCleaner(['unavailable'], {
      dryRun: false,
      force: false,
      modules: ['unavailable'],
    });
    const report = await cleaner.clean();

    expect(report.modules).toHaveLength(0);
    expect(mockModule.clean).not.toHaveBeenCalled();
  });

  it('should collect errors from failed clean operations', async () => {
    const { getModuleByIds } = await import('../modules/index.js');
    const errorModule = {
      id: 'error',
      name: 'Error',
      description: 'Test',
      isAvailable: vi.fn(() => true),
      analyze: vi.fn(),
      clean: vi.fn().mockRejectedValue(new Error('Clean failed')),
    };

    vi.mocked(getModuleByIds).mockReturnValue([errorModule as unknown as Module]);

    const cleaner = createCleaner(['error'], { dryRun: false, force: false, modules: ['error'] });
    const report = await cleaner.clean();

    expect(report.errors).toHaveLength(1);
    expect(report.errors[0]).toContain('Error: Clean failed');
  });

  it('should respect dryRun option', async () => {
    const { getModuleByIds } = await import('../modules/index.js');
    const cleaner = createCleaner(['test'], { dryRun: true, force: false, modules: ['test'] });
    await cleaner.clean();

    const mockModules = vi.mocked(getModuleByIds).mock.results[0].value;
    const module = mockModules[0] as unknown as Module;
    expect(module.clean).toHaveBeenCalledWith(true);
  });
});
