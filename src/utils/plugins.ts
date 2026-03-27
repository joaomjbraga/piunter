import { existsSync, mkdirSync, readdirSync, readFileSync, writeFileSync } from 'fs';
import { join } from 'path';
import { getHomeDir } from './os.js';
import type { AnalysisResult, CleaningResult } from '../types/index.js';

export interface Plugin {
  id: string;
  name: string;
  description: string;
  version: string;
  isAvailable(): boolean;
  analyze?(): Promise<AnalysisResult>;
  clean?(dryRun?: boolean, force?: boolean): Promise<CleaningResult>;
}

export interface PluginConfig {
  enabled: boolean;
  options: Record<string, unknown>;
}

const PLUGIN_DIR = join(getHomeDir(), '.piunter', 'plugins');
const PLUGIN_CONFIG_FILE = join(getHomeDir(), '.piunter', 'plugin-config.json');

class PluginManager {
  private plugins: Map<string, Plugin> = new Map();
  private config: Record<string, PluginConfig> = {};

  constructor() {
    this.loadConfig();
  }

  private loadConfig(): void {
    if (existsSync(PLUGIN_CONFIG_FILE)) {
      try {
        const content = readFileSync(PLUGIN_CONFIG_FILE, 'utf-8');
        this.config = JSON.parse(content);
      } catch {
        this.config = {};
      }
    }
  }

  private saveConfig(): void {
    try {
      const configDir = join(getHomeDir(), '.piunter');
      if (!existsSync(configDir)) {
        mkdirSync(configDir, { recursive: true });
      }
      writeFileSync(PLUGIN_CONFIG_FILE, JSON.stringify(this.config, null, 2), 'utf-8');
    } catch {
      // Silently fail if config cannot be saved
    }
  }

  isPluginEnabled(pluginId: string): boolean {
    return this.config[pluginId]?.enabled ?? true;
  }

  enablePlugin(pluginId: string): void {
    if (!this.config[pluginId]) {
      this.config[pluginId] = { enabled: true, options: {} };
    } else {
      this.config[pluginId].enabled = true;
    }
    this.saveConfig();
  }

  disablePlugin(pluginId: string): void {
    if (!this.config[pluginId]) {
      this.config[pluginId] = { enabled: false, options: {} };
    } else {
      this.config[pluginId].enabled = false;
    }
    this.saveConfig();
  }

  getPluginConfig(pluginId: string): PluginConfig {
    return this.config[pluginId] || { enabled: true, options: {} };
  }

  updatePluginOptions(pluginId: string, options: Record<string, unknown>): void {
    if (!this.config[pluginId]) {
      this.config[pluginId] = { enabled: true, options: {} };
    }
    this.config[pluginId].options = { ...this.config[pluginId].options, ...options };
    this.saveConfig();
  }

  registerPlugin(plugin: Plugin): void {
    this.plugins.set(plugin.id, plugin);
  }

  unregisterPlugin(pluginId: string): void {
    this.plugins.delete(pluginId);
  }

  getPlugin(pluginId: string): Plugin | undefined {
    return this.plugins.get(pluginId);
  }

  getAllPlugins(): Plugin[] {
    return Array.from(this.plugins.values());
  }

  getEnabledPlugins(): Plugin[] {
    return this.getAllPlugins().filter(p => this.isPluginEnabled(p.id));
  }

  async loadExternalPlugins(): Promise<void> {
    if (!existsSync(PLUGIN_DIR)) {
      return;
    }

    try {
      const files = readdirSync(PLUGIN_DIR);
      for (const file of files) {
        if (file.endsWith('.js') || file.endsWith('.ts')) {
          try {
            const pluginPath = join(PLUGIN_DIR, file);
            const plugin = await import(pluginPath);
            if (plugin.default) {
              this.registerPlugin(plugin.default);
            }
          } catch {
            // Skip invalid plugins
          }
        }
      }
    } catch {
      // Plugin directory not accessible
    }
  }

  listPlugins(): { id: string; name: string; enabled: boolean }[] {
    return this.getAllPlugins().map(p => ({
      id: p.id,
      name: p.name,
      enabled: this.isPluginEnabled(p.id),
    }));
  }
}

export const pluginManager = new PluginManager();
