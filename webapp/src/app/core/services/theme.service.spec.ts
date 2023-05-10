import { TestBed } from '@angular/core/testing';

import { ThemeService } from './theme.service';

describe('ThemeService', () => {
  let service: ThemeService;

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
      service.isDarkTheme.subscribe((dark) => {
        expect(dark).toBeFalsy();
        done();
      });
    });

    it('should store to local storage', () => {
      service.setDarkTheme(true);
      expect(localStorage.getItem('darkTheme')).toEqual('dark');
      service.setDarkTheme(false);
      expect(localStorage.getItem('darkTheme')).toEqual('light');
    });
  });

  describe('with localStorage set', () => {
    beforeEach(() => {
      localStorage.setItem('darkTheme', 'dark');
      TestBed.configureTestingModule({});
      service = TestBed.inject(ThemeService);
    });

    it('should be initialized from localStorage', (done) => {
      service.isDarkTheme.subscribe((dark) => {
        expect(dark).toBeTruthy();
        done();
      });
    });
  });
});
