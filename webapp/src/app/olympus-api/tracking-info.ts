import { StreamInfo } from './stream-info';

export class TrackingInfo {
  public total_bytes: number = 0;
  public free_bytes: number = 0;
  public bytes_per_second: number = 0;
  public stream?: StreamInfo;

  static fromPlain(plain: any): TrackingInfo | undefined {
    if (plain == undefined) {
      return undefined;
    }
    let ret = new TrackingInfo();
    ret.total_bytes = plain.total_bytes || 0;
    ret.free_bytes = plain.free_bytes || 0;
    ret.bytes_per_second = plain.bytes_per_second || 0;
    if (plain.stream != undefined) {
      ret.stream = StreamInfo.fromPlain(plain.stream);
    }
    return ret;
  }
}
