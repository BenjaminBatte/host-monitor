import {
  Component,
  Input,
  OnChanges,
  SimpleChanges,
} from '@angular/core';
import { CommonModule, NgClass, DecimalPipe, DatePipe } from '@angular/common';
import { ChartConfiguration, ChartData, ChartType } from 'chart.js';
import { NgChartsModule } from 'ng2-charts';


@Component({
  selector: 'app-host-card',
  standalone: true,
  imports: [
    CommonModule,
    NgClass,
    DecimalPipe,
    DatePipe,
    NgChartsModule, 
  ],
  templateUrl: './host-card.component.html',
  styleUrls: ['./host-card.component.scss'],
})
export class HostCardComponent implements OnChanges {
  @Input() host!: string;
  @Input() data!: any;

  toastMessage = '';

  public latencyChartData: ChartData<'line'> = {
    labels: [],
    datasets: [
      {
        data: [],
        label: 'Latency (ms)',
        fill: false,
        tension: 0.3,
        pointRadius: 1,
        borderWidth: 1,
      },
    ],
  };
  public chartType: ChartType = 'line';

ngOnChanges(changes: SimpleChanges): void {
  if (this.data?.latencyHistory) {
    const history = this.data.latencyHistory.slice(-20); 
this.latencyChartData.labels = history.map((_: number, i: number) => i.toString());
    this.latencyChartData.datasets[0].data = history;
  }
}


  copyToClipboard(value: string): void {
    navigator.clipboard.writeText(value).then(() => {
      this.toastMessage = 'Copied to clipboard!';
      setTimeout(() => (this.toastMessage = ''), 2000);
    });
  }
}
