import { Type, plainToClass } from 'class-transformer';
import { StreamInfo } from './stream-info';

export class TrackingInfo {
  public total_bytes: number = 0;
  public free_bytes: number = 0;
  public bytes_per_second: number = 0;
  @Type(() => StreamInfo)
  public stream?: StreamInfo;

  static fromPlain(plain: any): TrackingInfo {
    return plainToClass(TrackingInfo, plain, { exposeDefaultValues: true });
  }
}
