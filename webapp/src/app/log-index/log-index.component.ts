import { Component, OnInit } from '@angular/core';
import { OlympusService } from '../olympus-api/services/olympus.service';
import { ServiceEventList, ServicesLogs } from '../olympus-api/service-event';
import { HumanizeService } from '../core/humanize.service';

@Component({
  selector: 'app-log-index',
  templateUrl: './log-index.component.html',
  styleUrls: ['./log-index.component.scss'],
})
export class LogIndexComponent implements OnInit {
  public logs: ServiceEventList[] = [];

  constructor(
    private olympus: OlympusService,
    private humanize: HumanizeService
  ) {}

  ngOnInit(): void {
    this.olympus.getLogs().subscribe((logs) => (this.logs = logs));
  }

  lastEventTime(log: ServiceEventList): string {
    const lastEvent = log.events.at(-1);
    if (lastEvent == undefined) {
      return 'never';
    }
    const now = new Date();

    let ellapsed = now.getTime() - lastEvent.start.getTime();
    if (lastEvent.end != undefined) {
      ellapsed = now.getTime() - lastEvent.end.getTime();
    }
    return this.humanize.humanizeDuration(ellapsed) + ' ago';
  }
}
