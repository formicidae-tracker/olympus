import { Injectable } from '@angular/core';
import { SwPush } from '@angular/service-worker';
import { NotificationSettingsService } from './notification-settings.service';
import {
  Observable,
  concat,
  filter,
  first,
  from,
  map,
  of,
  retry,
  switchMap,
  take,
} from 'rxjs';
import { NotificationSettings } from '../notification-settings';
import { OlympusService } from 'src/app/olympus-api/services/olympus.service';

export type PushSubscriptionStatus = 'non-accepted' | 'not-updated' | 'updated';

@Injectable({
  providedIn: 'root',
})
export class PushNotificationService {
  private serverPublicKey: string = '';

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

  public requestSubscriptionOnDemand(): Observable<boolean> {
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
        return this.notifications.getSettings().pipe(
          switchMap(
            (notifications: NotificationSettings) =>
              this.updateNotificationSettings(
                subscription.endpoint,
                notifications
              ).pipe(retry({ delay: 1000 }))
            // retry every second until either: a) succeed, b) have a
            // new settings or c) have no endpoint anymore.
          )
        );
      })
    );
  }

  private updateNotificationSettings(
    endpoint: string,
    notifications: Required<NotificationSettings>
  ): Observable<PushSubscriptionStatus> {
    return concat(
      of('not-updated' as PushSubscriptionStatus),
      this.olympus
        .updateNotificationSettings(endpoint, notifications)
        .pipe(map(() => 'updated' as PushSubscriptionStatus))
    );
  }

  private requestPushSubscriptionWhenHasKey(): Observable<boolean> {
    return this.olympus.getPushServerPublicKey().pipe(
      switchMap((key: string) => {
        this.serverPublicKey = key;
        if (key.length == 0) {
          return of();
        }
        return this.pushSubscriptionRequired();
      }),
      switchMap(() => this.requestPushSubscription())
    );
  }

  private pushSubscriptionRequired(): Observable<null> {
    return this.notifications.getSettings().pipe(
      filter((settings) => settings.needPushSubscription()),
      take(1),
      switchMap(() => this.push.subscription),
      filter((sub: PushSubscription | null) => sub == null),
      take(1),
      map(() => null)
    );
  }

  private requestPushSubscription(): Observable<boolean> {
    return from(
      this.push.requestSubscription({ serverPublicKey: this.serverPublicKey })
    ).pipe(map(() => true));
  }
}
