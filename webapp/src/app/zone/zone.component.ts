import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title} from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Subscription,timer } from 'rxjs';
import { OlympusService } from '@services/olympus';
import { ZoneReport } from '@models/zone-report';

@Component({
	selector: 'app-zone',
	templateUrl: './zone.component.html',
	styleUrls: ['./zone.component.css']
})


export class ZoneComponent implements OnInit,OnDestroy {
    zoneName: string;
    hostName: string;
	zone: ZoneReport;
	loading: boolean;
	update: Subscription;


    constructor(private route: ActivatedRoute,
				private title: Title,
				private olympus: OlympusService) {
		this.loading = true;
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
			this.loading = true;
			this.zone = null;
			return;
		}
		this.olympus.zoneReport(this.hostName,this.zoneName)
			.subscribe(
				(r) => {
					this.zone = r;
					this.loading = false;
				},
				() => {
					this.zone = null;
					this.loading = false;
				}
			);
	}

	hasClimate(): boolean {
		return this.zone != null && this.zone.climate != null;
	}

	loaded(): boolean {
		return this.loading == false && this.zone != null;
	}

	unavailable(): boolean {
		return this.loading == false && this.zone == null;
	}

	ngOnDestroy() {
		if ( this.update != null ) {
			this.update.unsubscribe();
		}
	}

}
