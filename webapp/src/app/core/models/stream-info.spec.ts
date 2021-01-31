import { StreamInfo } from './stream-info';


describe('StreamData', () => {

	it('should create an instance',() => {
		expect(new StreamInfo('/olympus/hls/someserver.m3u8')).toBeTruthy();
	});

	it('should compute the thumbnail address from the streamURL',() => {
		let testdata = [
			{
				streamURL:'https://example.com/hls/foo.m3u8',
				thumbnailURL:'https://example.com/foo.png',
			},
			{
				streamURL:'/olympus/hls/foo-bar.m3u8',
				thumbnailURL:'/olympus/foo-bar.png',
			},
			{
				streamURL:'',
				thumbnailURL:'',
			},

		]

		for ( let d of testdata) {
			let sd = new StreamInfo(d.streamURL);
			expect(sd.streamURL).toBe(d.streamURL);
			expect(sd.thumbnailURL).toBe(d.thumbnailURL);
		}

	});

});
