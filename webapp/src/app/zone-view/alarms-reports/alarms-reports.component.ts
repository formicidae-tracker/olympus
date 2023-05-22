import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { Subscription, timer } from 'rxjs';
import { HumanizeDurationService } from 'src/app/core/humanize-duration.service';
import { AlarmReport } from 'src/app/olympus-api/alarm-report';

@Component({
  selector: 'app-alarms-reports',
  templateUrl: './alarms-reports.component.html',
  styleUrls: ['./alarms-reports.component.scss'],
})
export class AlarmsReportsComponent implements OnInit, OnDestroy {
  @Input() alarms: AlarmReport[] = [];

  public now: Date = new Date();
  private subscription?: Subscription;

  constructor(private humanizer: HumanizeDurationService) {}

  ngOnInit(): void {
    this.subscription = timer(0, 1000).subscribe(() => {
      this.now = new Date();
    });
  }

  ngOnDestroy(): void {
    if (this.subscription != undefined) {
      this.subscription.unsubscribe();
    }
  }

  public colorForReport(a: AlarmReport): string {
    if (a.on() == false) {
      return '';
    }
    if (a.level > 0) {
      return 'warn';
    }
    return 'accent';
  }

  public iconForReport(a: AlarmReport): string {
    return a.level > 0 ? 'error' : 'warning';
  }

  public identifyReport(index: number, item: AlarmReport): string {
    return item.identification;
  }

  public since(r: AlarmReport): string {
    return this.humanizer.humanize(this.now.getTime() - r.time().getTime());
  }
}
