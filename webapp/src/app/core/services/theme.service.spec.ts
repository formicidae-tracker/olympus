import { TestBed } from '@angular/core/testing';

import { ThemeService, themeKey } from './theme.service';

describe('ThemeService', () => {
  let service: ThemeService;

  afterEach(() => {
    localStorage.clear();
  });

  describe('with cleared localStorage', () => {
    beforeEach(() => {
      localStorage.clear();
      TestBed.configureTestingModule({});
      service = TestBed.inject(ThemeService);
    });

    it('should be created', () => {
      expect(service).toBeTruthy();
    });

    it('should default to light theme', (done) => {
      service.isDarkTheme().subscribe((dark) => {
        expect(dark).toBeFalsy();
        done();
      });
    });

    it('should store to localStorage when darkTheme is modified', () => {
      expect(localStorage.getItem(themeKey)).toBeNull();
      service.darkTheme = false;
      expect(localStorage.getItem(themeKey)).toBeNull();
      service.darkTheme = true;
      expect(localStorage.getItem(themeKey)).toEqual('dark');
    });
  });

  describe('with localStorage set', () => {
    beforeEach(() => {
      localStorage.setItem(themeKey, 'dark');
      TestBed.configureTestingModule({});
      service = TestBed.inject(ThemeService);
    });

    it('should have darkTheme set', (done) => {
      service.isDarkTheme().subscribe((dark) => {
        expect(dark).toBeTrue();
        done();
      });
    });
  });
});
