import {
  AfterViewInit,
  Component,
  Input,
  NgZone,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import { MatButtonToggleGroup } from '@angular/material/button-toggle';
import { Subscription, timer } from 'rxjs';

import { EChartsOption } from 'echarts';
import { ClimateTimeSeries } from 'src/app/olympus-api/climate-time-series';
import { OlympusService } from 'src/app/olympus-api/services/olympus.service';
import { UserSettingsService } from 'src/app/core/services/user-settings.service';

@Component({
  selector: 'app-climate-chart',
  templateUrl: './climate-chart.component.html',
  styleUrls: ['./climate-chart.component.scss'],
})
export class ClimateChartComponent implements OnInit, OnDestroy, AfterViewInit {
  @Input() host: string = '';
  @Input() zone: string = '';

  @ViewChild(MatButtonToggleGroup) _windowGroup!: MatButtonToggleGroup;

  public timeSeries?: ClimateTimeSeries;

  private _subscription?: Subscription = undefined;
  private _window: string = '1h';

  private _dark: boolean = false;
  options: EChartsOption = {};
  updateOptions: EChartsOption = {};

  get window(): string {
    return this._window;
  }

  set window(value: string) {
    this._unsubscribe();
    this._window = value;
    this._setInterval();
  }
  constructor(
    private olympus: OlympusService,
    private settings: UserSettingsService,
    private ngZone: NgZone
  ) {}

  ngAfterViewInit(): void {
    this._windowGroup.change.subscribe((change) => {
      this.window = change.value;
    });
  }

  private _setInterval(): void {
    this.ngZone.runOutsideAngular(() => {
      this._subscription = timer(0, 10000).subscribe(() => {
        this.ngZone.run(() => this._updateData());
      });
    });
  }

  ngOnInit(): void {
    this._setInterval();
    this._setUpChartOptions();
    this.settings.isDarkTheme().subscribe((dark) => {
      this._dark = dark;
      this._updateChart();
    });
  }

  private _unsubscribe(): void {
    if (this._subscription == undefined) {
      return;
    }
    this._subscription.unsubscribe();
    this._subscription = undefined;
  }

  ngOnDestroy(): void {
    this._unsubscribe();
  }

  private _updateData(): void {
    this.olympus
      .getClimateTimeSeries(this.host, this.zone, this._window)
      .subscribe((series) => {
        this.timeSeries = series;
        this._updateChart();
      });
  }

  private _setUpChartOptions(): void {
    this.options = {
      textStyle: {
        fontFamily: 'Roboto, sans-serif',
        fontSize: '14px',
        fontWeight: 400,
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross',
          label: { show: false },
        },
      },
      legend: { data: ['Humidity', 'Temperature'] },
      axisPointer: {
        animation: false,
      },
      xAxis: {
        type: 'value',
        name: 'Time',
        nameLocation: 'middle',
        nameGap: 24,
        min: 'minData',
        max: 0,
        axisLine: { show: true },
        axisLabel: { formatter: '{value}s' },
        splitLine: { show: false },
      },
      yAxis: [
        {
          type: 'value',
          name: 'Temperature',
          position: 'right',
          min: (mm) => mm.min - 2,
          max: (mm) => mm.max + 2,
          alignTicks: true,
          axisLine: { show: true },
          axisLabel: { formatter: (value: number) => value.toFixed(1) + '°C' },
        },
        {
          type: 'value',
          name: 'Relative Humidity',
          position: 'left',
          alignTicks: true,
          min: (mm) => mm.min - 5,
          max: (mm) => mm.max + 5,
          axisLine: { show: true },
          axisLabel: {
            formatter: (value: number) => value.toFixed(1) + '%',
          },
        },
      ],
      series: [
        {
          symbol: 'none',
          name: 'Humidity',
          yAxisIndex: 1,
          type: 'line',
          data: [],
        },
        {
          symbol: 'none',
          name: 'Temperature',
          yAxisIndex: 0,
          type: 'line',
          data: [],
        },
      ],
    };
  }

  private _updateChart(): void {
    let normal: string = this._dark ? '#acaeae' : '#37393a';
    let light: string = this._dark ? '#2e3132' : '#e1e3e3';

    let primary: string = this._dark ? '#006874' : '#009fb1';
    let accent: string = this._dark ? '#c4aa31' : '#e2c54b';

    let units: string = '';
    let humidity: number[][] = [];
    let temperature: number[][] = [];
    if (this.timeSeries != undefined) {
      units = this.timeSeries.units;
      humidity = this.timeSeries.humidity;
      temperature = this.timeSeries.temperature;
    }

    this.updateOptions = {
      legend: { textStyle: { color: normal }, inactiveColor: light },
      xAxis: {
        axisLine: { lineStyle: { color: normal } },
        axisLabel: { formatter: '{value}' + units },
      },
      yAxis: [
        {
          axisLine: { lineStyle: { color: normal } },
          splitLine: { lineStyle: { color: light } },
        },
        {
          axisLine: { lineStyle: { color: normal } },
          splitLine: { lineStyle: { color: light } },
        },
      ],
      series: [
        {
          data: humidity,
          lineStyle: { color: primary, width: 4 },
          itemStyle: { color: primary },
        },
        {
          data: temperature,
          lineStyle: { color: accent, width: 4 },
          itemStyle: { color: accent },
        },
      ],
    };
  }
}
