import { Injectable } from '@angular/core';
import {
  HumanizeDuration,
  HumanizeDurationLanguage,
} from 'humanize-duration-ts';

const prefixes = ['', 'Ki', 'Mi', 'Gi', 'Ti', 'Pi', 'Ei', 'Zi', 'Yi'];

@Injectable({
  providedIn: 'root',
})
export class HumanizeService {
  private language: HumanizeDurationLanguage = new HumanizeDurationLanguage();
  private humanizer: HumanizeDuration = new HumanizeDuration(this.language);
  constructor() {}

  public humanizeDuration(ms: number, largest: number = 2): string {
    return this.humanizer.humanize(ms, { largest: largest });
  }

  public findBinaryPrefixAndDivider(
    value: number | number[]
  ): [string, number] {
    if (value instanceof Array) {
      value = Math.max(...value);
    }
    let prefix: string = '';
    let divider: number = 1;
    for (prefix of prefixes) {
      if (Math.abs(value) < 1024) {
        break;
      }
      value /= 1024;
      divider *= 1024;
    }

    return [prefix, divider];
  }

  public humanizeBytes(value: number): string {
    const [prefix, divider] = this.findBinaryPrefixAndDivider(value);
    return (value / divider).toFixed(1) + ' ' + prefix + 'B';
  }

  public humanizeByteFraction(numerator: number, divisor: number): string {
    const [prefix, divider] = this.findBinaryPrefixAndDivider([
      numerator,
      divisor,
    ]);
    return (
      (numerator / divider).toFixed(1) +
      ' / ' +
      (divisor / divider).toFixed(1) +
      ' ' +
      prefix +
      'B'
    );
  }
}
