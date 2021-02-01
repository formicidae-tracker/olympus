import { StreamInfo } from './stream-info';


describe('StreamInfo', () => {

	it('should create an instance',() => {
		expect(new StreamInfo('/olympus/hls/someserver.m3u8','')).toBeTruthy();
	});

	it('should adapt from null',() => {
		expect(StreamInfo.adapt(null)).toBeNull()
	});

	it('should adapt without a thumbnail',() => {
		let s = StreamInfo.adapt({StreamURL: 'foo.m3u8'});

		expect(s).toBeTruthy();
		expect(s.streamURL).toBe('foo.m3u8');
		expect(s.thumbnailURL).toBe('');
	});

	it('should adapt with a thumbnail',() => {
		let s = StreamInfo.adapt({StreamURL: 'foo.m3u8',ThumbnailURL:'foo.png'})

		expect(s).toBeTruthy();
		expect(s.streamURL).toBe('foo.m3u8');
		expect(s.thumbnailURL).toBe('foo.png');

	});

});
