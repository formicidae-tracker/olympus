import {
  ApplicationRef,
  Component,
  NgZone,
  OnDestroy,
  OnInit,
} from '@angular/core';
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
  public zone: ZoneReport = new ZoneReport();
  public state: ViewState = 'loading';

  private _identifier: [string, string] = ['', ''];
  private _subscription?: Subscription;

  public now: Date = new Date(0);

  constructor(
    private route: ActivatedRoute,
    private olympus: OlympusService,
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
        () => {
          this.zone = new ZoneReport();
          this.now = new Date();
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
