export class Bounds {
  public minimum?: number;
  public maximum?: number;

  static fromPlain(plain: Partial<Bounds>): Bounds {
    let res = new Bounds();
    res.minimum = plain.minimum;
    res.maximum = plain.maximum;
    return res;
  }
}
