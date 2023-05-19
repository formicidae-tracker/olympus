import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription, timer } from 'rxjs';
import { ZoneReport } from '../olympus-api/zone-report';
import { OlympusService } from '../olympus-api/services/olympus.service';
import { AlarmReport } from '../olympus-api/alarm-report';

type ViewState = 'loading' | 'success' | 'error';

@Component({
  selector: 'app-zone-view',
  templateUrl: './zone-view.component.html',
  styleUrls: ['./zone-view.component.scss'],
})
export class ZoneViewComponent implements OnInit, OnDestroy {
  constructor(private route: ActivatedRoute, private olympus: OlympusService) {}

  public zone: ZoneReport = new ZoneReport();
  public state: ViewState = 'loading';

  private _identifier: [string, string] = ['', ''];
  private _update?: Subscription;

  ngOnInit(): void {
    this.state = 'loading';
    this.route.paramMap.subscribe((params) => {
      let host = String(params.get('host'));
      let zone = String(params.get('zone'));
      this._identifier = [host, zone];
      this._update = timer(0, 5000).subscribe(() => this.updateZone());
    });
  }

  ngOnDestroy(): void {
    if (!this._update) {
      return;
    }
    this._update.unsubscribe();
  }

  private updateZone() {
    this.olympus
      .getZoneReport(this._identifier[0], this._identifier[1])
      .subscribe(
        (report) => {
          this.state = 'success';
          this.zone = report;
          this.zone.alarms.sort(AlarmReport.compareFunction);
        },
        () => {
          this.zone = new ZoneReport();
          this.state = 'error';
        }
      );
  }

  public zoneName(): string {
    if (this._identifier[0].length == 0 && this._identifier[1].length == 0) {
      return '<unknown_zone>';
    }
    return this._identifier[0] + '.' + this._identifier[1];
  }
}
