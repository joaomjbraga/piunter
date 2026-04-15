import cliSpinners from 'cli-spinners';

export interface ProgressOptions {
  text: string;
  total?: number;
  color?: 'cyan' | 'green' | 'yellow' | 'red' | 'magenta';
}

export class ProgressBar {
  private current = 0;
  private total = 100;
  private text: string;
  private color: string;
  private width = 40;
  private interval: ReturnType<typeof setInterval> | null = null;
  private spinnerIndex = 0;

  constructor(options: ProgressOptions) {
    this.text = options.text;
    this.total = options.total || 100;
    this.color = options.color || 'cyan';
  }

  start(): void {
    this.render();
    this.interval = setInterval(() => this.tick(), 100);
  }

  tick(): void {
    this.spinnerIndex = (this.spinnerIndex + 1) % cliSpinners.dots.frames.length;
    this.render();
  }

  update(current: number): void {
    this.current = Math.min(Math.max(current, 0), this.total);
    this.render();
  }

  increment(amount = 1): void {
    this.update(this.current + amount);
  }

  setTotal(total: number): void {
    this.total = total;
  }

  private clearLine(): void {
    const columns = process.stdout.columns || 80;
    process.stdout.write(`\r${' '.repeat(columns)}\r`);
  }

  render(): void {
    if (this.total === 0) {
      const spinner = cliSpinners.dots.frames[this.spinnerIndex];
      process.stdout.write(`\r${spinner} ${this.text}`);
      return;
    }

    const filled = Math.round((this.current / this.total) * this.width);
    const empty = this.width - filled;
    const bar = '█'.repeat(filled) + '░'.repeat(empty);
    const percent = Math.round((this.current / this.total) * 100);
    const spinner = cliSpinners.dots.frames[this.spinnerIndex];

    process.stdout.write(`\r${spinner} [${bar}] ${percent}% ${this.text}`);
  }

  stop(message?: string): void {
    if (this.interval) {
      clearInterval(this.interval);
      this.interval = null;
    }
    const text = message || `${this.text} - Concluído`;
    this.clearLine();
    console.log(`  ${text}`);
  }

  fail(message?: string): void {
    if (this.interval) {
      clearInterval(this.interval);
      this.interval = null;
    }
    const text = message || `${this.text} - Falhou`;
    this.clearLine();
    console.log(`  ✗ ${text}`);
  }
}

export function createProgressBar(options: ProgressOptions): ProgressBar {
  return new ProgressBar(options);
}
