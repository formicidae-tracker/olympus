import { Component, OnInit } from '@angular/core';
import { olympus_version } from 'src/environments/version';
import { OlympusService } from '../olympus-api/services/olympus.service';

@Component({
  selector: 'app-version',
  templateUrl: './version.component.html',
  styleUrls: ['./version.component.scss'],
})
export class VersionComponent implements OnInit {
  public frontend_version: string = olympus_version;
  public backend_version: string = '<unknown>';

  constructor(private olympus: OlympusService) {}

  ngOnInit(): void {
    this.olympus
      .getVersion()
      .subscribe((version) => (this.backend_version = version));
  }
}
