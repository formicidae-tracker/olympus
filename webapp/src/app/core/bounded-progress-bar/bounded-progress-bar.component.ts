import { Component, Input } from '@angular/core';
import { ProgressBarMode } from '@angular/material/progress-bar';

@Component({
  selector: 'app-bounded-progress-bar',
  templateUrl: './bounded-progress-bar.component.html',
  styleUrls: ['./bounded-progress-bar.component.scss'],
})
export class BoundedProgressBarComponent {
  @Input() public value: number | undefined = undefined;
  @Input() public minimum: number | undefined = undefined;
  @Input() public maximum: number | undefined = undefined;

  _minimum(): number {
    return this.minimum || 0.0;
  }
  _maximum(): number {
    return this.maximum || 100.0;
  }

  public progressValue(): number {
    if (this.value == undefined) {
      return 0.0;
    }

    let ratio =
      (100.0 * (this.value - this._minimum())) /
      (this._maximum() - this._minimum());

    return Math.min(Math.max(ratio, 0.0), 100.0);
  }

  public progressMode(): ProgressBarMode {
    if (this.value == undefined) {
      return 'query';
    }
    return 'determinate';
  }

  public progressColor(): string {
    if (this.value == undefined) {
      return 'primary';
    }
    if (
      (this.minimum != undefined && this.value < this.minimum) ||
      (this.maximum != undefined && this.value > this.maximum)
    ) {
      return 'warn';
    }
    const range = 0.05 * (this._maximum() - this._minimum());
    if (
      (this.maximum != undefined && this.value > this.maximum - range) ||
      (this.minimum != undefined && this.value < this.minimum + range)
    ) {
      return 'accent';
    }

    return 'primary';
  }
}
