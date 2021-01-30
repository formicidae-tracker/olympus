import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title} from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Bounds } from '@models/bounds';
import { Zone } from '@models/zone';
import { Subscription,timer } from 'rxjs';
import { ZoneService } from '../zone.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';

@Component({
	selector: 'app-zone',
	templateUrl: './zone.component.html',
	styleUrls: ['./zone.component.css']
})

export class ZoneComponent implements OnInit,OnDestroy {
    zoneName: string
    hostName: string
	zone: Zone
	notFound: boolean
	update : Subscription;
	streamUrl: string


    constructor(private route: ActivatedRoute,
				private title: Title,
				private zoneService: ZoneService,
				private httpClient: HttpClient,
			   ) {
		this.zone = null;
		this.notFound = false;
		this.streamUrl = '';
	}

    ngOnInit() {
        this.zoneName = this.route.snapshot.paramMap.get('zoneName');
        this.hostName = this.route.snapshot.paramMap.get('hostName');
        this.title.setTitle('Olympus: '+this.hostName+'.'+this.zoneName)
		this.update = timer(0,5000).subscribe( (x) => {
			this.zoneService.getZone(this.hostName,this.zoneName)
				.subscribe(
					(zone) => {
						this.zone = zone;
						this.notFound = false;
						if ( this.zone.Name == "box" ) {
							this.httpClient.get('/olympus/hls/'+ this.zone.Host + '.m3u8',{responseType: 'text'}).subscribe(
								(src) => {
									this.streamUrl = '/olympus/'+ this.zone.Host + '.m3u8';
								},
								(error) => {
									this.streamUrl = '';
								},
							);
						}
					},
					(error)  => {
						this.zone = null;
						this.notFound = true;
					},
					() => {

					});
		});
    }

	ngOnDestroy() {
		this.update.unsubscribe();
	}

}
