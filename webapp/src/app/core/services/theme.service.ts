import { Injectable } from '@angular/core';
import { LocalStorageService } from './local-storage.service';
import { BehaviorSubject, Observable } from 'rxjs';

export const themeKey = 'theme';

@Injectable({
  providedIn: 'root',
})
export class ThemeService {
  private dark$: BehaviorSubject<boolean>;

  constructor(private localStorage: LocalStorageService) {
    const dark = this.localStorage.getItem(themeKey) == 'dark';
    this.dark$ = new BehaviorSubject<boolean>(dark);
  }

  public isDarkTheme(): Observable<boolean> {
    return this.dark$.asObservable();
  }

  public set darkTheme(dark: boolean) {
    if (this.dark$.value == dark) {
      return;
    }
    this.localStorage.setItem(themeKey, dark ? 'dark' : 'light');
    this.dark$.next(dark);
  }
}
