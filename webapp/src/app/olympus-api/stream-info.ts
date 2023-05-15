import { plainToClass } from 'class-transformer';

export class StreamInfo {
  public experiment_name: string = '';
  public stream_URL: string = '';
  public thumbnail_URL: string = '';

  static fromPlain(plain: any): StreamInfo {
    return plainToClass(StreamInfo, plain, {
      exposeDefaultValues: true,
    });
  }
}
