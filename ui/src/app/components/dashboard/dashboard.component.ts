import {
  Component,
  OnInit,
  ChangeDetectorRef,
  NgZone,
  ViewChild,
} from '@angular/core';
import { HttpClient } from '@angular/common/http';
import {
  CommonModule,
  DatePipe,
  NgClass,
  DecimalPipe,
} from '@angular/common';
import { FormsModule } from '@angular/forms';
import { HostCardComponent } from '../host-card/host-card.component';
import { NgChartsModule, BaseChartDirective } from 'ng2-charts';
import {
  ChartConfiguration,
  ChartData,
  ChartType,
} from 'chart.js';
import { MetricsService } from '../../services/metrics.service';
const DEBOUNCE_DELAY = 100;
@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
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
  isCompactMode = false;
  threshold: number = 70;

  constructor(
    private http: HttpClient,
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

    
    this.http.get<{ threshold: number }>('/api/threshold').subscribe({
      next: (res) => {
        this.threshold = res.threshold;
        console.log('Threshold fetched:', this.threshold);
      },
      error: (err) => console.error('Failed to fetch threshold:', err),
    });
  }

 private updateThresholdTimeout: any;


updateThreshold(): void {
  if (this.updateThresholdTimeout) {
    clearTimeout(this.updateThresholdTimeout);
  }

  this.updateThresholdTimeout = setTimeout(() => {
    this.http.post('/api/threshold', { threshold: this.threshold }).subscribe({
      next: () => console.log('Threshold updated to', this.threshold),
      error: (err) => console.error('Failed to update threshold', err),
    });
  }, DEBOUNCE_DELAY);
}


  exportToCSV(): void {
    const headers = [
      'Host',
      'Status',
      'Latency (ms)',
      'Uptime (%)',
      'Total Checks',
      'Successful Checks',
      'Last Checked',
    ];
    const rows = this.hostEntries.map(([host, data]) => [
      `"${host.replace(/"/g, '""')}"`,
      data.up ? 'UP' : 'DOWN',
      data.latency ?? 'N/A',
      ((data.successCount / (data.totalChecks || 1)) * 100).toFixed(1),
      data.totalChecks || 0,
      data.successCount || 0,
      `"${data.lastChecked ? new Date(data.lastChecked).toLocaleString('en-US', {
        dateStyle: 'medium',
        timeStyle: 'short',
      }) : 'N/A'}"`,
    ]);
    const csvContent = [headers.join(','), ...rows.map((row) => row.join(','))].join('\n');
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `host_status_${new Date().toISOString().slice(0, 10)}.csv`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  }

  getUpHostsCount(): number {
    return this.hostEntries.filter(([_, d]: any) => d.up).length;
  }

  getDownHostsCount(): number {
    return this.hostEntries.filter(([_, d]: any) => !d.up).length;
  }

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
      if (
        hostA !== hostB ||
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
