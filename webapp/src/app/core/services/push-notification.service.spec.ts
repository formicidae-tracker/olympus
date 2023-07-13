import { TestBed } from '@angular/core/testing';

import {
  Observable,
  of,
  firstValueFrom,
  throwError,
  concat,
  map,
  BehaviorSubject,
  finalize,
  startWith,
  reduce,
  defer,
  interval,
  take,
} from 'rxjs';

import { SwPush } from '@angular/service-worker';
import { HttpClientModule } from '@angular/common/http';

import {
  PushNotificationService,
  PushSubscriptionStatus,
  serverPublicKey,
} from './push-notification.service';
import { NotificationSettingsService } from './notification-settings.service';
import { OlympusService } from 'src/app/olympus-api/services/olympus.service';
import { NotificationSettings } from '../notification-settings';

class StubSwPush {
  public isEnabled: boolean = false;
  public reject: boolean = true;
  public subscription: Observable<PushSubscription | null>;

  private _subscription$: BehaviorSubject<PushSubscription | null>;

  constructor() {
    this._subscription$ = new BehaviorSubject<PushSubscription | null>(null);
    this.subscription = this._subscription$.asObservable();
  }

  public unsubscribe(): Promise<void> {
    this._subscription$.next(null);
    return of(void 0).toPromise();
  }

  public requestSubscription({
    serverPublicKey = '',
  }): Promise<PushSubscription> {
    return this._requestSubscription({ serverPublicKey: serverPublicKey });
  }

  public _requestSubscription({
    serverPublicKey = '',
  }): Promise<PushSubscription> {
    if (this.reject == true) {
      this._subscription$.next(null);
      return firstValueFrom(throwError('user rejected'));
    }
    const s = { endpoint: 'there' } as PushSubscription;
    this._subscription$.next(s);
    return firstValueFrom(this._subscription$ as Observable<PushSubscription>);
  }

  public complete() {
    this._subscription$.complete();
  }
}

describe('PushNotificationService', () => {
  let service: PushNotificationService;
  let push: StubSwPush;
  let notifications: jasmine.SpyObj<NotificationSettingsService>;
  let olympus: jasmine.SpyObj<OlympusService>;
  beforeEach(() => {
    push = new StubSwPush();
    notifications = jasmine.createSpyObj('NotificationSettingsService', [
      'getSettings',
    ]);
    olympus = jasmine.createSpyObj('OlympusService', [
      'getPushServerPublicKey',
      'updateNotificationSettings',
      'registerPushSubscription',
    ]);

    olympus.registerPushSubscription.and.returnValue(of(void 0));

    TestBed.configureTestingModule({
      imports: [HttpClientModule],
      providers: [
        { provide: SwPush, useValue: push },
        { provide: NotificationSettingsService, useValue: notifications },
        { provide: OlympusService, useValue: olympus },
      ],
    });

    service = TestBed.inject(PushNotificationService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('without service worker', () => {
    beforeEach(() => {
      localStorage.removeItem(serverPublicKey);
      push.isEnabled = false;
      olympus.getPushServerPublicKey.and.returnValue(of(''));
    });

    it('should not updateNotificationOnDemand', (done: DoneFn) => {
      service.updateNotificationsOnDemand().subscribe({
        next: (status) => {
          fail('Got unexpected status: ' + status);
        },
        error: (e) => {
          fail('Got unexpected error: ' + e);
        },
        complete: () => {
          expect(true).toBe(true);
          done();
        },
      });
    });

    it('should not subscribeOnDemand', (done: DoneFn) => {
      service.requestSubscriptionOnDemand().subscribe({
        next: (value) => {
          fail('Got unexpected value: ' + value);
        },
        error: (e) => {
          fail('Got unexpected error: ' + e);
        },
        complete: () => {
          expect(true).toBe(true);
          done();
        },
      });
    });
  });

  describe('with a service worker but no subscription', () => {
    beforeEach(() => {
      localStorage.removeItem(serverPublicKey);
      push.isEnabled = true;
      push.reject = true;
      olympus.getPushServerPublicKey.and.returnValue(of('youGotTheMagicWord'));
    });

    it('should not ask the user settings until needed ', (done: DoneFn) => {
      let modifyingSettings = false;

      spyOn(push, 'requestSubscription').and.callFake(() => {
        if (modifyingSettings == false) {
          fail('should not have been called');
        }
        return push._requestSubscription({ serverPublicKey: 'foo' });
      });

      notifications.getSettings.and.returnValue(
        concat(
          of(new NotificationSettings()),
          of(new NotificationSettings({ subscribeToAll: true })).pipe(
            map((s) => {
              modifyingSettings = true;
              return s;
            })
          )
        )
      );

      service
        .requestSubscriptionOnDemand()
        .pipe(
          finalize(() => {
            push.complete();
          })
        )
        .subscribe({
          next: () => {
            fail('Unexpected subscription completed');
          },
          error: (e) => {
            expect(e).toBe('user rejected');
          },
        });

      service
        .updateNotificationsOnDemand()
        .pipe(
          reduce((acc, value) => {
            acc.push(value);
            return acc;
          }, [] as PushSubscriptionStatus[])
        )
        .subscribe((statuses) => {
          expect(statuses).toEqual(['non-accepted']);
          done();
        });
    });

    it('should update once user subscribe', (done: DoneFn) => {
      push.reject = false;

      notifications.getSettings.and.returnValue(
        of(new NotificationSettings({ subscribeToAll: true })).pipe(
          map((v) => {
            return v;
          })
        )
      );

      olympus.updateNotificationSettings.and.callFake(() => {
        push.complete();
        return of(void 0);
      });

      service.requestSubscriptionOnDemand().subscribe({
        next: (v) => {
          expect(v).toBeUndefined();
        },
        error: (e) => {
          fail('unexpected subscription error: ' + e);
        },
      });

      service
        .updateNotificationsOnDemand()
        .pipe(
          reduce((acc, value) => {
            acc.push(value);
            return acc;
          }, [] as PushSubscriptionStatus[])
        )
        .subscribe((statuses) => {
          expect(statuses).toEqual(['not-updated', 'updated']);
          expect(localStorage.getItem(serverPublicKey)).toEqual(
            'youGotTheMagicWord'
          );
          done();
        });
    });
  });

  describe('with a service worker and a subscription', () => {
    beforeEach(() => {
      localStorage.setItem(serverPublicKey, 'youGotTheMagicWord');

      push.isEnabled = true;
      push.reject = false;
      push
        .requestSubscription({ serverPublicKey: 'youGotTheMagicWord' })
        .then();
      olympus.getPushServerPublicKey.and.returnValue(of('youGotTheMagicWord'));
    });

    it('should update on any push', (done: DoneFn) => {
      notifications.getSettings.and.returnValue(
        of(
          new NotificationSettings({ subscribeToAll: true }),
          new NotificationSettings()
        )
      );

      olympus.updateNotificationSettings.and.returnValue(of(void 0));

      service.requestSubscriptionOnDemand().subscribe({
        next: (value) => {
          fail('should not have subscribed: ' + value);
        },
        error: (e) => {
          fail('should not have tried to subscribe: ' + e);
        },
      });

      service
        .updateNotificationsOnDemand()
        .pipe(
          reduce((acc, value) => {
            acc.push(value);
            if (acc.length >= 4) {
              push.complete();
            }
            return acc;
          }, [] as PushSubscriptionStatus[])
        )
        .subscribe((statuses) => {
          expect(statuses).toEqual([
            'not-updated',
            'updated',
            'not-updated',
            'updated',
          ]);
          done();
        });
    });

    it('should retry if something goes wrong', (done: DoneFn) => {
      service.retryDelay = 10;
      service.retryIncrease = false;
      notifications.getSettings.and.returnValue(of(new NotificationSettings()));

      let count = 0;
      const httpFakeCall = () =>
        new Promise<void>((resolve, reject) => {
          count += 1;
          if (count < 3) {
            reject('500');
          }
          push.complete();
          resolve(void 0);
        });

      olympus.updateNotificationSettings.and.returnValue(defer(httpFakeCall));
      service.requestSubscriptionOnDemand().subscribe({
        next: (value) => {
          fail('should not have subscribed: ' + value);
        },
        error: (e) => {
          fail('should not have tried to subscribe: ' + e);
        },
      });

      service
        .updateNotificationsOnDemand()
        .pipe(
          reduce((acc, value) => {
            acc.push(value);
            return acc;
          }, [] as PushSubscriptionStatus[])
        )
        .subscribe((statuses) => {
          expect(statuses).toEqual([
            'not-updated',
            'not-updated',
            'not-updated',
            'updated',
          ]);
          done();
        });
    });

    it('should stop retrying on new settings', (done: DoneFn) => {
      service.retryDelay = 11;
      const toSend = [
        new NotificationSettings({ subscribeToAll: true }),
        new NotificationSettings(),
      ];

      notifications.getSettings.and.returnValue(
        interval(20).pipe(
          take(toSend.length),
          map((i) => toSend[i] as Required<NotificationSettings>)
        )
      );

      olympus.updateNotificationSettings.and.callFake((update) => {
        if (update.settings.subscribeToAll == true) {
          return defer(() => throwError('500'));
        }
        push.complete();
        return of(void 0);
      });

      service.requestSubscriptionOnDemand().subscribe({
        next: (value) => {
          fail('should not have subscribed: ' + value);
        },
        error: (e) => {
          fail('should not have tried to subscribe: ' + e);
        },
      });

      service
        .updateNotificationsOnDemand()
        .pipe(
          reduce((acc, value) => {
            acc.push(value);
            return acc;
          }, [] as PushSubscriptionStatus[])
        )
        .subscribe((statuses) => {
          expect(statuses).toEqual(['not-updated', 'not-updated', 'updated']);
          done();
        });
    });
  });
});
