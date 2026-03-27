import { describe, it, expect, vi, beforeEach } from 'vitest';
import { Analyzer, createAnalyzer } from './analyzer.js';
import type { Module } from '../modules/index.js';

describe('Analyzer', () => {
  let analyzer: Analyzer;

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('should create analyzer with specific modules', () => {
    const analyzer = createAnalyzer(['test']);
    expect(analyzer).toBeDefined();
  });

  it('should calculate summary correctly', () => {
    const analyzer = createAnalyzer(['test']);
    const mockResults = [
      {
        module: 'test',
        items: [{ path: '/test', size: 1000, type: 'file', description: 'test' }],
        totalSize: 1000,
      },
    ];

    const summary = analyzer.getSummary(mockResults);

    expect(summary.totalSize).toBe(1000);
    expect(summary.totalItems).toBe(1);
    expect(summary.byModule.test).toEqual({ size: 1000, items: 1 });
  });

  it('should calculate summary with multiple modules', () => {
    const analyzer = createAnalyzer(['test']);
    const mockResults = [
      {
        module: 'test1',
        items: [{ path: '/test1', size: 500, type: 'file', description: 't1' }],
        totalSize: 500,
      },
      {
        module: 'test2',
        items: [{ path: '/test2', size: 1500, type: 'file', description: 't2' }],
        totalSize: 1500,
      },
    ];

    const summary = analyzer.getSummary(mockResults);

    expect(summary.totalSize).toBe(2000);
    expect(summary.totalItems).toBe(2);
    expect(summary.byModule.test1).toEqual({ size: 500, items: 1 });
    expect(summary.byModule.test2).toEqual({ size: 1500, items: 1 });
  });

  it('should handle empty results', () => {
    const analyzer = createAnalyzer(['test']);
    const summary = analyzer.getSummary([]);

    expect(summary.totalSize).toBe(0);
    expect(summary.totalItems).toBe(0);
    expect(summary.byModule).toEqual({});
  });
});
