import { Component, OnInit, OnDestroy } from '@angular/core';
import { Title } from '@angular/platform-browser';
import { OlympusService } from '@services/olympus';
import { ZoneSummaryReport }  from '@models/zone-summary-report';
import { Subscription,timer } from 'rxjs';
@Component({
    selector: 'app-home',
    templateUrl: './home.component.html',
    styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit,OnDestroy {
    summaries: ZoneSummaryReport[];
	update: Subscription;


    constructor(private olympus : OlympusService, private title: Title) {
		this.summaries = [];
	}

    ngOnInit() {
		this.title.setTitle('Olympus: Home')

		this.update = timer(0,20000).subscribe(() => {
			this.olympus.zoneSummaries().subscribe( (list) => {
				this.summaries = list;
			});
		})
    }

	ngOnDestroy() {
		this.update.unsubscribe()
	}

}
