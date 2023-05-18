import { TestBed } from '@angular/core/testing';

import { userSettingsKey, UserSettingsService } from './user-settings.service';
import { UserSettings } from './user-settings';

describe('UserSettingsService', () => {
  let service: UserSettingsService;

  afterEach(() => {
    localStorage.clear();
  });

  describe('with cleared localStorage', () => {
    beforeEach(() => {
      localStorage.clear();
      TestBed.configureTestingModule({});
      service = TestBed.inject(UserSettingsService);
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

    it('should default to no subscription', (done) => {
      service.isSubscribedToAlarmFrom('foo').subscribe((s) => {
        expect(s).toBeFalse();
        done();
      });
    });

    it('should store to localStorage when darkTheme is modified', () => {
      expect(localStorage.getItem(userSettingsKey)).toBeNull();
      service.setDarkTheme(false);
      expect(localStorage.getItem(userSettingsKey)).toBeNull();
      service.setDarkTheme(true);
      expect(localStorage.getItem(userSettingsKey)).not.toBeNull();
    });

    it('should store to localStorage when subscription is modified', () => {
      expect(localStorage.getItem(userSettingsKey)).toBeNull();
      service.unsubscribeFromAlarmFrom('foo'); // unsubscribed by default
      expect(localStorage.getItem(userSettingsKey)).toBeNull();
      service.subscribeToAlarmFrom('foo');
      expect(localStorage.getItem(userSettingsKey)).not.toBeNull();
    });
  });

  describe('with localStorage set', () => {
    beforeEach(() => {
      localStorage.setItem(
        userSettingsKey,
        new UserSettings({
          darkMode: true,
          alarmSubscriptions: new Set<string>(['foo', 'bar']),
        }).serialize()
      );
      TestBed.configureTestingModule({});
      service = TestBed.inject(UserSettingsService);
    });

    it('should have darkTheme set', (done) => {
      service.isDarkTheme().subscribe((dark) => {
        expect(dark).toBeTrue();
        done();
      });
    });

    it('should have "foo" subscribed', (done) => {
      service.isSubscribedToAlarmFrom('foo').subscribe((subscribed) => {
        expect(subscribed).toBeTrue();
        done();
      });
    });

    it('should have "bar" subscribed', (done) => {
      service.isSubscribedToAlarmFrom('bar').subscribe((subscribed) => {
        expect(subscribed).toBeTrue();
        done();
      });
    });

    it('should have "new" unsubscribed', (done) => {
      service.isSubscribedToAlarmFrom('new').subscribe((subscribed) => {
        expect(subscribed).toBeFalse();
        done();
      });
    });
  });
});
