import { Component, OnInit, Input } from '@angular/core';
import { Zone } from '../core/zone.model';
import { ZoneService } from '../zone.service';


@Component({
  selector: 'app-zone-preview',
  templateUrl: './zone-preview.component.html',
  styleUrls: ['./zone-preview.component.css']
})
export class ZonePreviewComponent implements OnInit {
    @Input() zone: Zone;
	imagePath: string


    constructor(private zs : ZoneService) {
		this.imagePath = '';
	}

    ngOnInit() {
		if ( this.zone.Name == "box" ) {
			this.zs.hasStream(this.zone.Host).subscribe(
				(src) => {
					if ( src == true ) {
						this.imagePath = '/olympus/'+ this.zone.Host + '.png';
					}
				});
		}
    }
}
