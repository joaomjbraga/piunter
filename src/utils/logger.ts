import chalk from 'chalk';
import ora, { type Ora } from 'ora';
import cliSpinners from 'cli-spinners';

export type LogLevel = 'info' | 'success' | 'warn' | 'error' | 'debug';

class Logger {
  private spinner: Ora | null = null;
  private enabled: boolean = true;

  setEnabled(enabled: boolean): void {
    this.enabled = enabled;
  }

  info(message: string): void {
    if (!this.enabled) return;
    console.log(`${chalk.blue('ℹ')} ${message}`);
  }

  success(message: string): void {
    if (!this.enabled) return;
    console.log(`${chalk.green('✓')} ${chalk.green(message)}`);
  }

  warn(message: string): void {
    if (!this.enabled) return;
    console.log(`${chalk.yellow('⚠')} ${chalk.yellow(message)}`);
  }

  error(message: string): void {
    if (!this.enabled) return;
    console.log(`${chalk.red('✗')} ${chalk.red(message)}`);
  }

  debug(message: string): void {
    if (!this.enabled) return;
    if (process.env.DEBUG) {
      console.log(`${chalk.gray('[DEBUG]')} ${message}`);
    }
  }

  startSpinner(text: string): void {
    this.spinner = ora({
      text,
      spinner: cliSpinners.dots,
      color: 'cyan',
    }).start();
  }

  updateSpinner(text: string): void {
    if (this.spinner) {
      this.spinner.text = text;
    }
  }

  stopSpinner(success: boolean = true, message?: string): void {
    if (!this.spinner) return;
    if (success) {
      this.spinner.succeed(message || 'Concluído');
    } else {
      this.spinner.fail(message || 'Falhou');
    }
    this.spinner = null;
  }

  title(text: string): void {
    if (!this.enabled) return;
    console.log(`\n${chalk.bold.cyan('━'.repeat(50))}`);
    console.log(`  ${chalk.bold.cyan(text)}`);
    console.log(`${chalk.bold.cyan('━'.repeat(50))}\n`);
  }

  subtitle(text: string): void {
    if (!this.enabled) return;
    console.log(`\n${chalk.cyan('› ')}${chalk.bold(text)}`);
  }

  item(message: string, size?: string): void {
    if (!this.enabled) return;
    if (size) {
      console.log(`  ${chalk.gray('•')} ${message} ${chalk.dim(`(${size})`)}`);
    } else {
      console.log(`  ${chalk.gray('•')} ${message}`);
    }
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
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  }
}

export const logger = new Logger();
