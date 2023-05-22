import { Component, Input, OnInit, ViewChild } from '@angular/core';
import { MatPaginator } from '@angular/material/paginator';
import { MatTableDataSource } from '@angular/material/table';
import { HumanizeDurationService } from 'src/app/core/humanize-duration.service';
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

  constructor(private humanizer: HumanizeDurationService) {}

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
    return this.humanizer.humanize(d);
  }
}
