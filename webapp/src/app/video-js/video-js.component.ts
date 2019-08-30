import { Component, OnInit, OnChanges, Input } from '@angular/core';

import videojs from 'video.js';

@Component({
	selector: 'app-video-js',
	templateUrl: './video-js.component.html',
	styleUrls: ['./video-js.component.css']
})
export class VideoJsComponent implements OnInit,OnChanges {
	public vjs: videojs.Player;
	@Input() urlVideo: string;
	@Input() urlPoster: string;
	constructor() {	}

	ngOnInit() {
	}

	ngOnChanges() {
		if ( this.urlVideo.length == 0 ) {
			return;
		}
		const options = {
			'sources' : [{
				'src' : this.urlVideo,
				'type' : 'application/x-mpegURL'
			}
						],
			'poster' : this.urlPoster
		};
		this.vjs = videojs('my-player', options);
	}
}
