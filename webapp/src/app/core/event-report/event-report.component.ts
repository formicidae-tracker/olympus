import { DatePipe } from '@angular/common';
import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { HumanizeService } from 'src/app/core/humanize.service';
import { Event } from 'src/app/olympus-api/event';

@Component({
  selector: 'app-event-report',
  templateUrl: './event-report.component.html',
  styleUrls: ['./event-report.component.scss'],
})
export class EventReportComponent implements OnInit {
  @Input() events: Event[] = [];
  dataSource = new MatTableDataSource<Event>();

  @ViewChild(MatPaginator, { static: true }) paginator!: MatPaginator;

  constructor(private humanizer: HumanizeService, private date: DatePipe) {}

  public columnsToDisplay = ['start', 'end', 'duration'];

  public trackByTime(index: number, e: Event): number {
    return e.time().getTime();
  }

  ngOnInit() {
    this.dataSource = new MatTableDataSource<Event>(this.events);
    this.dataSource.paginator = this.paginator;
  }

  public duration(e: Event): string {
    const d = e.duration();
    if (d == undefined) {
      return '';
    }
    return this.humanizer.humanizeDuration(d);
  }

  public formatStartDate(e: Event): string {
    return this.date.transform(e.start, 'dd/MM/yy HH:mm:ss') || '';
  }

  public formatEndDate(e: Event): string {
    if (e.end == undefined) {
      return 'now';
    }
    if (e.end.toLocaleDateString() == e.start.toLocaleDateString()) {
      return this.date.transform(e.end, 'HH:mm:ss') || '';
    }
    return this.date.transform(e.end, 'dd/MM/yy HH:mm:ss') || '';
  }
}
