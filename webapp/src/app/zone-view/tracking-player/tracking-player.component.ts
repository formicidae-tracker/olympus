import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-tracking-player',
  templateUrl: './tracking-player.component.html',
  styleUrls: ['./tracking-player.component.scss'],
})
export class TrackingPlayerComponent {
  @Input() src: string = '';
  @Input() thumbnail: string = '';
}
