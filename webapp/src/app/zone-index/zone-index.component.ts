import { Component, NgZone, OnDestroy, OnInit } from '@angular/core';
import { ZoneReportSummary } from '../olympus-api/zone-report-summary';
import { OlympusService } from '../olympus-api/services/olympus.service';
import { Subscription, timer } from 'rxjs';

@Component({
  selector: 'app-zone-index',
  templateUrl: './zone-index.component.html',
  styleUrls: ['./zone-index.component.scss'],
})
export class ZoneIndexComponent implements OnInit, OnDestroy {
  public zones: ZoneReportSummary[] = [];
  public offline: boolean = false;
  private _subscription?: Subscription;

  constructor(private olympus: OlympusService, private ngZone: NgZone) {}

  ngOnInit(): void {
    this.ngZone.runOutsideAngular(() => {
      this._subscription = timer(0, 10000).subscribe(() => {
        this.ngZone.run(() => this._updateZones());
      });
    });
  }

  ngOnDestroy(): void {
    if (this._subscription == undefined) {
      return;
    }
    this._subscription.unsubscribe();
    this._subscription = undefined;
  }

  private _updateZones(): void {
    this.olympus.getZoneReportSummaries().subscribe(
      (zones) => {
        this.offline = false;
        this.zones = zones;
      },
      () => {
        this.zones = [];
        this.offline = true;
      }
    );
  }
}
