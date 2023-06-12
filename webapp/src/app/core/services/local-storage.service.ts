import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class LocalStorageService {
  setItem(key: string, value: string): void {
    localStorage.setItem(key, value);
  }

  getItem(key: string): string | null {
    return localStorage.getItem(key);
  }
}

@Injectable()
export class NullLocalStorageService extends LocalStorageService {
  override setItem(key: string, value: string): void {}

  override getItem(key: string): string | null {
    return null;
  }
}
