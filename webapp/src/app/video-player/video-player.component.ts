import { Component, OnInit, OnChanges, OnDestroy, Input } from '@angular/core';

@Component({
	selector: 'app-video-js',
	templateUrl: './video-player.component.html',
	styleUrls: ['./video-player.component.css']
})

export class VideoPlayerComponent implements OnInit,OnChanges,OnDestroy {
	@Input() url: string;
	@Input() thumbnail: string;
	constructor() {	}

	ngOnInit() {
	}

	ngOnChanges() {
	}

	ngOnDestroy() {
	}
}
