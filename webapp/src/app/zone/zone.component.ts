import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title} from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { ZoneClimateReport } from '@models/zone-climate-report';
import { Subscription,timer } from 'rxjs';
import { OlympusService } from '@services/olympus';

@Component({
	selector: 'app-zone',
	templateUrl: './zone.component.html',
	styleUrls: ['./zone.component.css']
})

export class ZoneComponent implements OnInit,OnDestroy {
    zoneName: string
    hostName: string
	zone: ZoneClimateReport
	notFound: boolean
	update : Subscription;
	streamUrl: string


    constructor(private route: ActivatedRoute,
				private title: Title,
				private olympus: OlympusService) {
		this.zone = null;
		this.notFound = false;
		this.streamUrl = '';
	}

    ngOnInit() {
        this.zoneName = this.route.snapshot.paramMap.get('zoneName');
        this.hostName = this.route.snapshot.paramMap.get('hostName');
        this.title.setTitle('Olympus: '+this.hostName+'.'+this.zoneName)
		this.update = timer(0,5000).subscribe( (x) => {
			if ( this.zoneName == 'box' ) {
				this.olympus.streamURL(this.hostName)
					.subscribe(
						(streamURL) => {
							this.streamUrl = streamURL;
						},
						(error) => {
							this.streamUrl = '';
						}
					);
			}
			this.olympus.zoneClimate(this.hostName,this.zoneName)
				.subscribe(
					(zone) => {
						this.zone = zone;
						this.notFound = false;
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
