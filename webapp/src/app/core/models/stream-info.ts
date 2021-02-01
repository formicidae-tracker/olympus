export class StreamInfo {
	constructor(public streamURL: string = '',
				public thumbnailURL: string = '') {
	}

	static adapt(item: any): StreamInfo {
		if ( item == null ) {
			return null;
		}
		return new StreamInfo(item.StreamURL,item.ThumbnailURL);
	}
}
