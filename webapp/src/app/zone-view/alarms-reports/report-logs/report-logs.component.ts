import { DatePipe } from '@angular/common';
import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { HumanizeService } from 'src/app/core/humanize.service';
import { AlarmEvent } from 'src/app/olympus-api/alarm-report';

@Component({
  selector: 'app-report-logs',
  templateUrl: './report-logs.component.html',
  styleUrls: ['./report-logs.component.scss'],
})
export class ReportLogsComponent implements OnInit {
  @Input() events: AlarmEvent[] = [];
  dataSource = new MatTableDataSource<AlarmEvent>();

  @ViewChild(MatPaginator, { static: true }) paginator!: MatPaginator;

  constructor(private humanizer: HumanizeService, private date: DatePipe) {}

  public columnsToDisplay = ['start', 'end', 'duration'];

  public trackByTime(index: number, e: AlarmEvent): number {
    return e.time().getTime();
  }

  ngOnInit() {
    this.dataSource = new MatTableDataSource<AlarmEvent>(this.events);
    this.dataSource.paginator = this.paginator;
  }

  public duration(e: AlarmEvent): string {
    const d = e.duration();
    if (d == undefined) {
      return '';
    }
    return this.humanizer.humanizeDuration(d);
  }

  public formatStartDate(e: AlarmEvent): string {
    return this.date.transform(e.start, 'dd/MM/yy HH:mm:ss') || '';
  }

  public formatEndDate(e: AlarmEvent): string {
    if (e.end == undefined) {
      return 'now';
    }
    if (e.end.toLocaleDateString() == e.start.toLocaleDateString()) {
      return this.date.transform(e.end, 'HH:mm:ss') || '';
    }
    return this.date.transform(e.end, 'dd/MM/yy HH:mm:ss') || '';
  }
}
