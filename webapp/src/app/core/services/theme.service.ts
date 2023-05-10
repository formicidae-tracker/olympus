import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  private darkTheme = new BehaviorSubject<boolean>(
    localStorage.getItem('darkTheme') == 'dark'
  );

  public isDarkTheme = this.darkTheme.asObservable();

  public setDarkTheme(dark: boolean): void {
    this.darkTheme.next(dark);
    localStorage.setItem('darkTheme', dark ? 'dark' : 'light');
  }
}
