export class ClimateTimeSeries {
  units: string = '';
  reference: Date = new Date(0);
  humidity: any;
  temperature: any;
  temperatureAux: any;

  static fromPlain(plain: any): ClimateTimeSeries {
    let res = new ClimateTimeSeries();
    res.units = plain.units || '';
    res.reference = new Date(plain.reference || 0);
    res.humidity = plain.humidity;
    res.temperature = plain.temperature;
    res.temperatureAux = plain.temperatureAux;
    return res;
  }
}
