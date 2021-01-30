import { Component, OnInit, Input } from '@angular/core';
import { Zone } from '@models/zone';
import { HttpClient, HttpHeaders } from '@angular/common/http';


@Component({
  selector: 'app-zone-preview',
  templateUrl: './zone-preview.component.html',
  styleUrls: ['./zone-preview.component.css']
})

export class ZonePreviewComponent implements OnInit {
    @Input() zone: Zone;
	imagePath: string


    constructor(private httpClient : HttpClient) {
		this.imagePath = '';
	}

    ngOnInit() {
		if ( this.zone.Name == "box" ) {
			this.httpClient.get('/olympus/hls/'+ this.zone.Host + '.m3u8',{responseType: 'text'}).subscribe(
				(src) => {
					this.imagePath = '/olympus/'+ this.zone.Host + '.png';
				},
				(error) => {
					this.imagePath = '';
				},
			);
		}
    }
}
