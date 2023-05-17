import { Component, OnInit } from '@angular/core';
import { TitleService } from '../core/services/title.service';
import { ActivatedRoute } from '@angular/router';
import { Subscription, switchMap, timer } from 'rxjs';

const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms));

@Component({
  selector: 'app-zone-view',
  templateUrl: './zone-view.component.html',
  styleUrls: ['./zone-view.component.scss'],
})
export class ZoneViewComponent implements OnInit {
  constructor(private title: TitleService, private route: ActivatedRoute) {}

  private _host: string = '';
  private _zone: string = '';
  private _update?: Subscription;

  ngOnInit(): void {
    this.title.setTitle('Zone...');
    this.route.paramMap.subscribe((params) => {
      this._host = String(params.get('host'));
      this._zone = String(params.get('zone'));
      this.title.setTitle(this._host + '.' + this._zone);
      this._update = timer(2000, 5000).subscribe(() => this.updateZone());
    });
  }

  private updateZone() {}
}
