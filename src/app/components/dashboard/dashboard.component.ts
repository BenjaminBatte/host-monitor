import {
  Component,
  OnInit,
  ViewChild,
  ChangeDetectorRef,
  NgZone,
} from '@angular/core';
import {
  CommonModule,
  DatePipe,
  DecimalPipe,
  NgClass,
  NgFor,
  NgIf,
} from '@angular/common';
import { MetricsService } from '../../services/metrics.service';
import { HostCardComponent } from '../host-card/host-card.component';

import { NgChartsModule, BaseChartDirective } from 'ng2-charts';
import {
  ChartType,
  ChartConfiguration,
  ChartData,
} from 'chart.js';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    NgIf,
    NgFor,
    NgClass,
    DecimalPipe,
    DatePipe,
    HostCardComponent,
    NgChartsModule,
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit {
  hostEntries: any[] = [];

  public pieChartType: ChartType = 'pie';
  public pieChartLabels: string[] = ['UP', 'DOWN'];

  public pieChartData: ChartData<'pie'> = {
    labels: this.pieChartLabels,
    datasets: [
      {
        data: [0, 0],
        backgroundColor: ['#2ecc71', '#e74c3c'],
      },
    ],
  };

  public pieChartOptions: ChartConfiguration<'pie'>['options'] = {
    responsive: true,
    maintainAspectRatio: false,
    animation: {
      duration: 500,
      easing: 'easeOutQuart',
    },
  };

  @ViewChild(BaseChartDirective) chart?: BaseChartDirective;

  constructor(
    private metricsService: MetricsService,
    private cdr: ChangeDetectorRef,
    private zone: NgZone
  ) {}

  ngOnInit(): void {
    this.metricsService.getMetrics().subscribe((data) => {
      this.zone.run(() => {
        const newEntries = Object.entries(data);

        if (!this.areHostEntriesEqual(this.hostEntries, newEntries)) {
          this.hostEntries = newEntries;
          this.updateChartData();
        }
      });
    });
  }

  getUpHostsCount(): number {
    return this.hostEntries.filter(([_, d]: any) => d.up).length;
  }

  getDownHostsCount(): number {
    return this.hostEntries.filter(([_, d]: any) => !d.up).length;
  }

  private updateChartData(): void {
    const up = this.getUpHostsCount();
    const down = this.getDownHostsCount();

    const currentData = this.pieChartData.datasets[0].data;

    if (currentData[0] !== up || currentData[1] !== down) {
      this.pieChartData.datasets[0].data = [up, down];

     
      this.chart?.update();
    }
  }

  private areHostEntriesEqual(a: any[], b: any[]): boolean {
    if (a.length !== b.length) return false;

    for (let i = 0; i < a.length; i++) {
      const [hostA, dataA] = a[i];
      const [hostB, dataB] = b[i];

      if (hostA !== hostB) return false;

      if (
        dataA.up !== dataB.up ||
        dataA.latency !== dataB.latency ||
        dataA.lastChecked !== dataB.lastChecked ||
        dataA.successCount !== dataB.successCount ||
        dataA.totalChecks !== dataB.totalChecks
      ) {
        return false;
      }
    }

    return true;
  }
}
