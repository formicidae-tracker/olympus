import { Component, Input } from '@angular/core';
import { HumanizeService } from 'src/app/core/services/humanize.service';
import { AlarmReport } from 'src/app/olympus-api/alarm-report';

@Component({
  selector: 'app-alarms-reports',
  templateUrl: './alarms-reports.component.html',
  styleUrls: ['./alarms-reports.component.scss'],
})
export class AlarmsReportsComponent {
  @Input() alarms: AlarmReport[] = [];

  @Input() now: Date = new Date();

  constructor(private humanizer: HumanizeService) {}

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
    return this.humanizer.humanizeDuration(
      this.now.getTime() - r.time().getTime(),
      1
    );
  }
}
