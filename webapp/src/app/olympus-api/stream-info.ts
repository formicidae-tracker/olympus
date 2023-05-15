export class StreamInfo {
  public experiment_name: string = '';
  public stream_URL: string = '';
  public thumbnail_URL: string = '';

  static fromPlain(plain: any): StreamInfo | undefined {
    if (plain == undefined) {
      return undefined;
    }
    let res = new StreamInfo();
    res.experiment_name = plain.experiment_name || '';
    res.stream_URL = plain.stream_URL || '';
    res.thumbnail_URL = plain.thumbnail_URL || '';
    return res;
  }
}
