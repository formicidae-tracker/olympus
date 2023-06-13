import { Injectable } from '@angular/core';
import { SwPush } from '@angular/service-worker';
import { NotificationSettingsService } from './notification-settings.service';
import {
  Observable,
  concat,
  first,
  from,
  map,
  merge,
  of,
  retry,
  switchMap,
} from 'rxjs';
import { NotificationSettings } from '../notification-settings';
import { OlympusService } from 'src/app/olympus-api/services/olympus.service';

export type PushSubscriptionStatus = 'non-accepted' | 'not-updated' | 'updated';

@Injectable({
  providedIn: 'root',
})
export class PushNotificationService {
  private pushSubscriptionRequired$: Observable<null>;
  private updatePushSubscription$: Observable<PushSubscriptionStatus>;
  private requestPushSubscriptionWhenHasKey$: Observable<PushSubscriptionStatus>;

  private serverPublicKey: string = '';

  constructor(
    private push: SwPush,
    private notifications: NotificationSettingsService,
    private olympus: OlympusService
  ) {
    this.pushSubscriptionRequired$ = this.notifications.getSettings().pipe(
      first((settings) => settings.needPushSubscription()),
      switchMap(() => this.push.subscription),
      first((sub: PushSubscription | null) => sub == null),
      map(() => null)
    );

    this.updatePushSubscription$ = this.push.subscription.pipe(
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

    this.requestPushSubscriptionWhenHasKey$ = this.olympus
      .getPushServerPublicKey()
      .pipe(
        switchMap((key: string) => {
          this.serverPublicKey = key;
          if (key.length == 0) {
            return of();
          }
          return this.pushSubscriptionRequired$;
        }),
        switchMap(() => this.requestPushSubscription()),
        map(() => 'not-updated' as PushSubscriptionStatus)
      );
  }

  public updateNotificationsOnDemand(): Observable<PushSubscriptionStatus> {
    if (this.push.isEnabled == false) {
      console.log('Push Notification disabled');
      return of();
    }

    return merge(
      this.requestPushSubscriptionWhenHasKey$,
      this.updatePushSubscription$
    );
  }

  private requestPushSubscription(): Observable<PushSubscription> {
    return from(
      this.push.requestSubscription({ serverPublicKey: this.serverPublicKey })
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
}
