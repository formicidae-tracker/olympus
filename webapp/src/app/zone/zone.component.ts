import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title} from '@angular/platform-browser';
import { ActivatedRoute } from '@angular/router';
import { Subscription,timer } from 'rxjs';
import { OlympusService } from '@services/olympus';
import { ZoneReport } from '@models/zone-report';
import { map } from 'rxjs/operators';

@Component({
	selector: 'app-zone',
	templateUrl: './zone.component.html',
	styleUrls: ['./zone.component.css']
})

export class ZoneComponent implements OnInit,OnDestroy {
    zoneName: string;
    hostName: string;
	zone: ZoneReport;
	update : Subscription;


    constructor(private route: ActivatedRoute,
				private title: Title,
				private olympus: OlympusService) {
		this.zone = null;
		this.zoneName = '';
		this.hostName = '';
		this.update = null;
	}

    ngOnInit() {
		this.route.paramMap.pipe(
			map((params) => {
				this.hostName = params.get('hostName');
				this.zoneName = params.get('zoneName');
				this.title.setTitle('Olympus: '+this.hostName+'.'+this.zoneName)
				this.update = timer(0,5000).subscribe( () => {
					this.olympus.zoneReport(this.hostName,this.zoneName)
						.subscribe(
							(zone) => {
								this.zone = zone;
							},
							()  => {
								this.zone = null;
							},
							() => {

							});
				});
			}),
		);
    }

	ngOnDestroy() {
		if ( this.update != null ) {
			this.update.unsubscribe();
		}
	}

}
