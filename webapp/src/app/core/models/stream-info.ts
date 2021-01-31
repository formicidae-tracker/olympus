export class StreamInfo {
	public streamURL: string;
	public thumbnailURL: string;
	constructor(streamURL: string){
		this.streamURL = streamURL;
		this.thumbnailURL = StreamInfo.thumbnailFromStreamURL(streamURL);
	}

	static thumbnailFromStreamURL(streamURL: string) {
		if ( streamURL.length == 0 ) {
			return '';
		}
		let filename = streamURL.substring(streamURL.lastIndexOf('/')+1);
		let nodename = filename.replace('.m3u8','.png');
		let base = streamURL.substring(0,streamURL.lastIndexOf('/'));
		base = base.substring(0,base.lastIndexOf('/'));
		return base + '/' + nodename;
	}
}
