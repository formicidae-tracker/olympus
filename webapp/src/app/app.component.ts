import { Component, OnDestroy, OnInit } from '@angular/core';
import {
  Observable,
  Subscription,
  concat,
  map,
  of,
  retry,
  switchMap,
} from 'rxjs';

import { MatSnackBar, MatSnackBarRef } from '@angular/material/snack-bar';

import { SnackNetworkOfflineComponent } from './core/snack-network-offline/snack-network-offline.component';

import { NetworkStatusService } from './core/services/network-status.service';
import { ThemeService } from './core/services/theme.service';
import { SwPush } from '@angular/service-worker';
import { NotificationSettingsService } from './core/services/notification-settings.service';
import { NotificationSettings } from './core/notification-settings';

export type PushSubscriptionStatus = 'non-accepted' | 'not-updated' | 'updated';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent implements OnInit, OnDestroy {
  public darkTheme: boolean = false;

  private _snackOffline?: MatSnackBarRef<any>;
  private _subscriptions: Subscription[] = [];

  constructor(
    private theme: ThemeService,
    private snackBar: MatSnackBar,
    private networkStatus: NetworkStatusService,
    private notifications: NotificationSettingsService,
    private push: SwPush
  ) {}

  updatePushSubscriptionStatus(): Observable<PushSubscriptionStatus> {
    return this.push.subscription.pipe(
      switchMap((subscription: PushSubscription | null) => {
        if (subscription == null) {
          return of('non-accepted' as PushSubscriptionStatus);
        }
        return this.notifications.getSettings().pipe(
          switchMap((settings: Required<NotificationSettings>) => {
            return concat(
              of('not-updated' as PushSubscriptionStatus),
              //TODO: this of should call olympus to update the
              //endpoint notification settings.
              of({
                endpoint: subscription.endpoint,
                settings: settings,
              }).pipe(map(() => 'updated' as PushSubscriptionStatus))
            );
          }),
          // retry every second until either: a) succeed, b) have a
          // new settings or c) have no endpoint anymore.
          retry({ delay: 1000 })
        );
      })
    );
  }

  ngOnInit(): void {
    this._subscriptions.push(
      this.theme.isDarkTheme().subscribe((dark) => (this.darkTheme = dark))
    );
    this._subscriptions.push(
      this.networkStatus.online.subscribe((online) => {
        if (online == true) {
          this._dismissOffline();
        } else {
          this._indicateOffline();
        }
      })
    );

    this._subscriptions.push(
      this.updatePushSubscriptionStatus().subscribe(
        (status: PushSubscriptionStatus) => {
          console.log('PushSubscription status is :' + status);
        }
      )
    );
  }

  ngOnDestroy(): void {
    for (const s of this._subscriptions) {
      s.unsubscribe();
    }
  }

  private _indicateOffline() {
    if (this._snackOffline != undefined) {
      return;
    }
    this._snackOffline = this.snackBar.openFromComponent(
      SnackNetworkOfflineComponent
    );
  }

  private _dismissOffline() {
    this._snackOffline?.dismiss();
    this._snackOffline = undefined;
  }
}
