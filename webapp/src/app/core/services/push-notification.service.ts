import { Injectable } from '@angular/core';
import { SwPush } from '@angular/service-worker';
import { NotificationSettingsService } from './notification-settings.service';
import {
  Observable,
  RetryConfig,
  catchError,
  concat,
  delay,
  filter,
  from,
  map,
  of,
  retry,
  switchMap,
  take,
} from 'rxjs';

import { NotificationSettings } from '../notification-settings';
import { OlympusService } from 'src/app/olympus-api/services/olympus.service';
import { NotificationSettingsUpdate } from 'src/app/olympus-api/notification-settings-update';
import { HttpErrorResponse } from '@angular/common/http';

export type PushSubscriptionStatus = 'non-accepted' | 'not-updated' | 'updated';

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
    return this.push.subscription.pipe(
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
    return this.olympus.getPushServerPublicKey().pipe(
      switchMap((key: string) => {
        this.serverPublicKey = key;
        if (key.length == 0) {
          return of();
        }
        const savedKey = localStorage.getItem('serverPublicKey');

        return this.push.subscription.pipe(
          switchMap((subscription: PushSubscription | null) => {
            if (subscription != null && key != savedKey) {
              return concat(
                from(this.push.unsubscribe()),
                this.pushSubscriptionRequired()
              );
            }
            return this.pushSubscriptionRequired();
          })
        );
      }),
      switchMap(() => this.requestPushSubscription())
    );
  }

  private pushSubscriptionRequired(): Observable<void> {
    return this.notifications.getSettings().pipe(
      filter((settings) => settings.needPushSubscription()),
      take(1),
      switchMap(() => this.push.subscription),
      filter((sub: PushSubscription | null) => sub == null),
      take(1),
      map(() => void 0)
    );
  }

  private requestPushSubscription(): Observable<void> {
    return from(
      this.push.requestSubscription({ serverPublicKey: this.serverPublicKey })
    ).pipe(
      switchMap((s: PushSubscription) => {
        localStorage.setItem('serverPublicKey', this.serverPublicKey);
        return this.olympus.registerPushSubscription(s);
      })
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
