import chalk from 'chalk';

export type LogLevel = 'info' | 'success' | 'warn' | 'error' | 'debug';

class Logger {
  private enabled: boolean = true;

  setEnabled(enabled: boolean): void {
    this.enabled = enabled;
  }

  private getWidth(): number {
    return process.stdout.columns || 60;
  }

  info(message: string): void {
    if (!this.enabled) return;
    console.log(`  ${chalk.blue('*')} ${message}`);
  }

  success(message: string): void {
    if (!this.enabled) return;
    console.log(`  ${chalk.green('*')} ${chalk.green(message)}`);
  }

  warn(message: string): void {
    if (!this.enabled) return;
    console.log(`  ${chalk.yellow('!')} ${chalk.yellow(message)}`);
  }

  error(message: string): void {
    if (!this.enabled) return;
    console.log(`  ${chalk.red('x')} ${chalk.red(message)}`);
  }

  debug(message: string): void {
    if (!this.enabled) return;
    if (process.env.DEBUG) {
      console.log(`  ${chalk.gray('[debug]')} ${message}`);
    }
  }

  title(text: string): void {
    if (!this.enabled) return;
    console.log();
    console.log(`  ${chalk.cyan(text)}`);
    console.log(`  ${chalk.dim('─'.repeat(Math.min(this.getWidth() - 4, 40)))}`);
    console.log();
  }

  subtitle(text: string): void {
    if (!this.enabled) return;
    console.log();
    console.log(`  ${chalk.bold(text)}`);
    console.log();
  }

  item(message: string, value?: string): void {
    if (!this.enabled) return;
    const msg = value 
      ? `${message} ${chalk.cyan(`(${value})`)}`
      : message;
    console.log(`    ${chalk.dim('-')} ${msg}`);
  }

  list(items: { name: string; value: string; success?: boolean }[]): void {
    if (!this.enabled) return;
    items.forEach((item) => {
      const status = item.success === false 
        ? chalk.red('x')
        : item.success === true
          ? chalk.green('*')
          : chalk.dim('-');
      
      const line = `${item.name} ${chalk.cyan(`(${item.value})`)}`;
      console.log(`    ${status} ${line}`);
    });
  }

  space(): void {
    if (!this.enabled) return;
    console.log();
  }

  formatBytes(bytes: number): string {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(1))} ${sizes[i]}`;
  }
}

export const logger = new Logger();
