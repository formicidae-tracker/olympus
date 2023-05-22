import { Injectable } from '@angular/core';
import {
  HumanizeDuration,
  HumanizeDurationLanguage,
} from 'humanize-duration-ts';

@Injectable({
  providedIn: 'root',
})
export class HumanizeDurationService {
  private language: HumanizeDurationLanguage = new HumanizeDurationLanguage();
  private humanizer: HumanizeDuration = new HumanizeDuration(this.language);
  constructor() {}

  public humanize(ms: number): string {
    return this.humanizer.humanize(ms, { largest: 2 });
  }
}
