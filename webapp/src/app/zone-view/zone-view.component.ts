import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Subscription, timer } from 'rxjs';

@Component({
  selector: 'app-zone-view',
  templateUrl: './zone-view.component.html',
  styleUrls: ['./zone-view.component.scss'],
})
export class ZoneViewComponent implements OnInit {
  constructor(private route: ActivatedRoute) {}

  private _host: string = '';
  private _zone: string = '';
  private _update?: Subscription;

  ngOnInit(): void {
    this.route.paramMap.subscribe((params) => {
      this._host = String(params.get('host'));
      this._zone = String(params.get('zone'));
      this._update = timer(2000, 5000).subscribe(() => this.updateZone());
    });
  }

  private updateZone() {}
}
