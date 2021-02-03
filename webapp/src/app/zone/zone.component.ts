import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title} from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Subscription,timer } from 'rxjs';
import { OlympusService } from '@services/olympus';
import { ZoneReport } from '@models/zone-report';

export enum ZoneState {
	Loading = 1,
	Loaded = 2,
	Unavailable = 3,
};

@Component({
	selector: 'app-zone',
	templateUrl: './zone.component.html',
	styleUrls: ['./zone.component.css']
})


export class ZoneComponent implements OnInit,OnDestroy {
    zoneName: string;
    hostName: string;
	zone: ZoneReport;
	state: ZoneState;
	update: Subscription;


    constructor(private route: ActivatedRoute,
				private title: Title,
				private olympus: OlympusService) {
		this.state = ZoneState.Loading;
		this.zone = null;
		this.update = null;
		this.hostName = '';
		this.zoneName = '';
	}

    ngOnInit() {
		this.route.paramMap
			.subscribe((params) => {
				this.hostName = params.get('hostName');
				this.zoneName = params.get('zoneName');
				this.title.setTitle('Olympus: '+this.hostName+'.'+this.zoneName);
				this.update = timer(0,5000).subscribe( () => {
					this.updateZone();
				});
			});
    }

	updateZone(): void {
		if ( this.hostName.length == 0
			|| this.zoneName.length == 0 ) {
			this.state = ZoneState.Loading;
			this.zone = null;
			return;
		}
		this.olympus.zoneReport(this.hostName,this.zoneName)
			.subscribe(
				(r) => {
					this.zone = r;
					this.state = ZoneState.Loaded;
				},
				() => {
					this.zone = null;
					this.state = ZoneState.Unavailable;
				}
			);
	}

	ngOnDestroy() {
		if ( this.update != null ) {
			this.update.unsubscribe();
		}
	}

}
