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
import { PushNotificationService } from './core/services/push-notification.service';

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
    private push: PushNotificationService
  ) {}

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
      this.push.requestSubscriptionOnDemand().subscribe()
    );
    this._subscriptions.push(
      this.push.updateNotificationsOnDemand().subscribe()
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
