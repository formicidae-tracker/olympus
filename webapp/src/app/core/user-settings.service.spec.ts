import { cases } from 'jasmine-parameterized';

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
      service.hasSubscription('foo').subscribe((s) => {
        expect(s).toBeFalse();
        done();
      });
    });

    cases([
      ['darkTheme', false],
      ['subscribeToAll', false],
      ['notifyOnWarning', false],
    ]).it(
      'should store to localStorage when boolean is modified',
      ([property, defaultValue]) => {
        interface Indexable {
          [key: string]: boolean;
        }
        let plain: Indexable = {};
        expect(localStorage.getItem(userSettingsKey)).toBeNull();
        plain[property] = defaultValue;
        Object.assign(service, plain);
        expect(localStorage.getItem(userSettingsKey)).toBeNull();
        plain[property] = !defaultValue;
        Object.assign(service, plain);
        expect(localStorage.getItem(userSettingsKey)).not.toBeNull();
      }
    );

    it('should store to localStorage when subscription is modified', () => {
      expect(localStorage.getItem(userSettingsKey)).toBeNull();
      service.unsubscribeTo('foo'); // unsubscribed by default
      expect(localStorage.getItem(userSettingsKey)).toBeNull();
      service.subscribeTo('foo');
      expect(localStorage.getItem(userSettingsKey)).not.toBeNull();
    });
  });

  describe('with localStorage set', () => {
    beforeEach(() => {
      localStorage.setItem(
        userSettingsKey,
        new UserSettings({
          darkMode: true,
          subscriptions: new Set<string>(['foo', 'bar']),
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
      service.hasSubscription('foo').subscribe((subscribed) => {
        expect(subscribed).toBeTrue();
        done();
      });
    });

    it('should have "bar" subscribed', (done) => {
      service.hasSubscription('bar').subscribe((subscribed) => {
        expect(subscribed).toBeTrue();
        done();
      });
    });

    it('should have "new" unsubscribed', (done) => {
      service.hasSubscription('new').subscribe((subscribed) => {
        expect(subscribed).toBeFalse();
        done();
      });
    });
  });
});
