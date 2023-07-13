import { Injectable } from '@angular/core';
import { SwPush } from '@angular/service-worker';
import { NotificationSettingsService } from './notification-settings.service';
import {
  Observable,
  RetryConfig,
  catchError,
  combineLatest,
  concat,
  delay,
  filter,
  from,
  map,
  mapTo,
  of,
  retry,
  skip,
  startWith,
  switchMap,
  switchMapTo,
  take,
} from 'rxjs';
import { NotificationSettings } from '../notification-settings';
import { OlympusService } from 'src/app/olympus-api/services/olympus.service';
import { NotificationSettingsUpdate } from 'src/app/olympus-api/notification-settings-update';
import { HttpErrorResponse } from '@angular/common/http';

export type PushSubscriptionStatus = 'non-accepted' | 'not-updated' | 'updated';

export const serverPublicKey: string = 'serverPublicKey';

@Injectable({
  providedIn: 'root',
})
export class PushNotificationService {
  private serverPublicKey: string = '';

  public retryDelay: number = 2000;
  public retryIncrease: boolean = true;

  constructor(
    private push: SwPush,
    private notifications: NotificationSettingsService,
    private olympus: OlympusService
  ) {}

  public updateNotificationsOnDemand(): Observable<PushSubscriptionStatus> {
    if (this.push.isEnabled == false) {
      return of();
    }

    return this.updatePushSubscription();
  }

  public requestSubscriptionOnDemand(): Observable<void> {
    if (this.push.isEnabled == false) {
      return of();
    }
    return this.requestPushSubscriptionWhenHasKey();
  }

  private updatePushSubscription(): Observable<PushSubscriptionStatus> {
    return this.updatedPushSubscription().pipe(
      switchMap((subscription: PushSubscription | null) => {
        if (subscription == null) {
          return of('non-accepted' as PushSubscriptionStatus);
        }
        return this.notifications
          .getSettings()
          .pipe(
            switchMap((settings: NotificationSettings) =>
              this.updateNotificationSettings(subscription, settings).pipe(
                retry(this.retryConfig())
              )
            )
          );
      })
    );
  }

  private updatedPushSubscription(): Observable<PushSubscription | null> {
    return combineLatest([
      this.olympus.getPushServerPublicKey(),
      this.push.subscription,
    ]).pipe(
      switchMap(([key, subscription]: [string, PushSubscription | null]) => {
        this.serverPublicKey = key;
        const localKey = localStorage.getItem(serverPublicKey);
        if (subscription == null || localKey == this.serverPublicKey) {
          localStorage.setItem(serverPublicKey, this.serverPublicKey);
          return this.updateServerSubscription(1).pipe(startWith(subscription));
        }

        console.warn(
          'mismatched local and server public key server:' +
            this.serverPublicKey +
            ' local: ' +
            localKey
        );

        return from(this.push.unsubscribe()).pipe(
          switchMapTo(this.updateServerSubscription(0))
        );
      })
    );
  }

  private updateServerSubscription(
    toSkip: number
  ): Observable<PushSubscription | null> {
    return this.push.subscription.pipe(
      skip(toSkip),
      switchMap((subscription) => {
        if (subscription == null) {
          return of(null);
        }
        return this.olympus
          .registerPushSubscription(subscription)
          .pipe(mapTo(subscription));
      })
    );
  }

  private updateNotificationSettings(
    subscription: PushSubscription,
    settings: Required<NotificationSettings>
  ): Observable<PushSubscriptionStatus> {
    return concat(
      of('not-updated' as PushSubscriptionStatus),
      this.olympus
        .updateNotificationSettings(
          new NotificationSettingsUpdate(subscription.endpoint, settings)
        )
        .pipe(
          catchError((err: HttpErrorResponse, caught: Observable<void>) => {
            if (err.status == 404) {
              return concat(
                this.olympus.registerPushSubscription(subscription),
                caught
              );
            }
            throw err;
          }),
          map(() => 'updated' as PushSubscriptionStatus)
        )
    );
  }

  private requestPushSubscriptionWhenHasKey(): Observable<void> {
    return this.pushSubscriptionRequired().pipe(
      switchMap(() =>
        this.push.requestSubscription({
          serverPublicKey: this.serverPublicKey,
        })
      ),
      mapTo(void 0)
    );
  }

  private pushSubscriptionRequired(): Observable<void> {
    return this.updatedPushSubscription().pipe(
      filter((sub: PushSubscription | null) => sub == null),
      take(1),
      switchMap(() => this.notifications.getSettings()),
      filter((settings) => settings.needPushSubscription()),
      take(1),
      map(() => void 0)
    );
  }

  private retryConfig(): RetryConfig {
    if (this.retryIncrease == false) {
      return { delay: this.retryDelay };
    }
    return {
      delay: (err: any, retryCount: number) => {
        return of(null).pipe(
          delay(Math.min(30000, 2 ** retryCount * this.retryDelay))
        );
      },
    };
  }
}

@Injectable()
export class NullPushNotificationService {
  public updateNotificationsOnDemand(): Observable<PushSubscriptionStatus> {
    return of();
  }

  public requestSubscriptionOnDemand(): Observable<void> {
    return of();
  }
}
