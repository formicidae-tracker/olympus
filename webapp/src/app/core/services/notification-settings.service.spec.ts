import { cases } from 'jasmine-parameterized';

import { TestBed } from '@angular/core/testing';

import {
  userSettingsKey,
  NotificationSettingsService,
} from './notification-settings.service';
import { NotificationSettings } from '../notification-settings';

describe('NotificationSettingsService', () => {
  let service: NotificationSettingsService;

  afterEach(() => {
    localStorage.clear();
  });

  describe('with cleared localStorage', () => {
    beforeEach(() => {
      localStorage.clear();
      TestBed.configureTestingModule({});
      service = TestBed.inject(NotificationSettingsService);
    });

    it('should be created', () => {
      expect(service).toBeTruthy();
    });

    it('should default to no subscription', (done) => {
      service.hasSubscription('foo').subscribe((s) => {
        expect(s).toBeFalse();
        done();
      });
    });

    cases([
      ['notifyNonGraceful', false],
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
        new NotificationSettings({
          subscriptions: new Set<string>(['foo', 'bar']),
        }).serialize()
      );
      TestBed.configureTestingModule({});
      service = TestBed.inject(NotificationSettingsService);
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
