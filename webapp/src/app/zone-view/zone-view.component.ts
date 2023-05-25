import { Component, NgZone, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription, timer } from 'rxjs';
import { ZoneReport } from '../olympus-api/zone-report';
import { OlympusService } from '../olympus-api/services/olympus.service';
import { AlarmReport } from '../olympus-api/alarm-report';
import { HumanizeService } from '../core/humanize.service';
import { HttpErrorResponse } from '@angular/common/http';

type ViewState = 'loading' | 'success' | 'offline';

interface HasSince {
  since: Date;
}

@Component({
  selector: 'app-zone-view',
  templateUrl: './zone-view.component.html',
  styleUrls: ['./zone-view.component.scss'],
})
export class ZoneViewComponent implements OnInit, OnDestroy {
  public zone: ZoneReport = new ZoneReport();
  public state: ViewState = 'loading';

  private _identifier: [string, string] = ['', ''];
  private _subscription?: Subscription;

  public now: Date = new Date(0);

  constructor(
    private route: ActivatedRoute,
    private olympus: OlympusService,
    private humanizer: HumanizeService,
    private ngZone: NgZone
  ) {}

  ngOnInit(): void {
    this.state = 'loading';
    this.route.paramMap.subscribe((params) => {
      let host = String(params.get('host'));
      let zone = String(params.get('zone'));
      this._identifier = [host, zone];
      this.ngZone.runOutsideAngular(() => {
        this._subscription = timer(0, 10000).subscribe(() => {
          this.ngZone.run(() => this.updateZone());
        });
      });
    });
  }

  ngOnDestroy(): void {
    if (!this._subscription) {
      return;
    }
    this._subscription.unsubscribe();
  }

  private updateZone() {
    this.olympus
      .getZoneReport(this._identifier[0], this._identifier[1])
      .subscribe(
        (report) => {
          this.state = 'success';
          this.now = new Date();
          this.zone = report;
          this.zone.alarms.sort(AlarmReport.compareFunction);
        },
        (e: HttpErrorResponse) => {
          this.now = new Date();
          if (e.status == 404) {
            this.zone = new ZoneReport();
            this.state = 'offline';
          } else {
          }
        }
      );
  }

  public zoneName(): string {
    if (this._identifier[0].length == 0 && this._identifier[1].length == 0) {
      return '<unknown_zone>';
    }
    return this._identifier[0] + '.' + this._identifier[1];
  }

  public formatSince(s: HasSince | undefined): string {
    return this.humanizer.humanizeDuration(
      this.now.getTime() - (s?.since.getTime() || Infinity),
      1
    );
  }
}
