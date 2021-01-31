import { Component, OnInit, Input } from '@angular/core';
import { ZoneSummaryReport } from '@models/zone-summary-report';


@Component({
	selector: 'app-zone-preview',
	templateUrl: './zone-preview.component.html',
	styleUrls: ['./zone-preview.component.css']
})

export class ZonePreviewComponent implements OnInit {
    @Input() summary: ZoneSummaryReport;


    constructor() {
	}

    ngOnInit() {
    }
}
