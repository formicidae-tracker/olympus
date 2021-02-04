export class StreamInfo {
	constructor(public streamURL: string = '',
				public thumbnailURL: string = '') {
	}

	static adapt(item: any): StreamInfo {
		if ( item == null ) {
			return new StreamInfo();
		}
		return new StreamInfo(item.StreamURL,item.ThumbnailURL);
	}

	hasStream(): boolean {
		return this.streamURL.length>0
	}

	hasThumbnail(): boolean {
		return this.thumbnailURL.length>0
	}
}
