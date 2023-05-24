import { StreamInfo } from './stream-info';

export class TrackingInfo {
  public total_bytes: number = 0;
  public free_bytes: number = 0;
  public bytes_per_second: number = 0;
  public stream?: StreamInfo;
  public since: Date = new Date(0);

  public get used_bytes(): number {
    return Math.max(0, this.total_bytes - this.free_bytes);
  }

  public filledUpEta(): number {
    if (this.free_bytes <= 0) {
      return 0;
    }
    if (this.bytes_per_second <= 0) {
      return Infinity;
    }
    return (this.free_bytes / this.bytes_per_second) * 1000;
  }

  static fromPlain(plain: any): TrackingInfo | undefined {
    if (plain == undefined) {
      return undefined;
    }
    let ret = new TrackingInfo();
    ret.since = new Date(plain.since || 0);
    ret.total_bytes = plain.total_bytes || 0;
    ret.free_bytes = plain.free_bytes || 0;
    ret.bytes_per_second = plain.bytes_per_second || 0;
    if (plain.stream != undefined) {
      ret.stream = StreamInfo.fromPlain(plain.stream);
    }
    return ret;
  }
}
