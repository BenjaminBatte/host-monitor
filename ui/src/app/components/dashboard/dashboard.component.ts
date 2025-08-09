import {
  Component,
  OnInit,
  OnDestroy,
  ChangeDetectorRef,
  NgZone,
  ViewChild,
  HostListener,
} from '@angular/core';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import {
  CommonModule,
  DatePipe,
  NgClass,
  DecimalPipe,
} from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Subject, takeUntil } from 'rxjs';

import { HostCardComponent } from '../host-card/host-card.component';
import { MetricsService } from '../../services/metrics.service';

import { NgChartsModule, BaseChartDirective } from 'ng2-charts';
import { ChartConfiguration, ChartData, ChartType } from 'chart.js';

const DEBOUNCE_DELAY = 100;

type HostMetrics = {
  up: boolean;
  latency: number | null;
  lastChecked: string | number | Date | null;
  successCount: number;
  totalChecks: number;
};

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    NgClass,
    DecimalPipe,
    DatePipe,
    HttpClientModule,
    HostCardComponent,
    NgChartsModule,
  ],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
})
export class DashboardComponent implements OnInit, OnDestroy {
  hostEntries: Array<[string, HostMetrics]> = [];
  isCompactMode = false;
  threshold = 70;

  // Scroll helpers
  showScrollTop = false;

  private destroy$ = new Subject<void>();
  private updateThresholdTimeout?: number;

  constructor(
    private http: HttpClient,
    private metricsService: MetricsService,
    private cdr: ChangeDetectorRef,
    private zone: NgZone
  ) {}

  /* ==== Lifecycle ==== */
  ngOnInit(): void {
    // Live metrics subscription
    this.metricsService
      .getMetrics()
      .pipe(takeUntil(this.destroy$))
      .subscribe((data) => {
        this.zone.run(() => {
          const newEntries = this.normalizeEntries(
            Object.entries(data as Record<string, HostMetrics>)
          );
          if (!this.areHostEntriesEqual(this.hostEntries, newEntries)) {
            this.hostEntries = newEntries;
            this.updateChartData();
          }
        });
      });

    // Fetch initial threshold
    this.http
      .get<{ threshold: number }>('/api/threshold')
      .pipe(takeUntil(this.destroy$))
      .subscribe({
        next: (res) => (this.threshold = res.threshold),
        error: (err) => console.error('Failed to fetch threshold:', err),
      });
  }

  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
    if (this.updateThresholdTimeout !== undefined) {
      window.clearTimeout(this.updateThresholdTimeout);
    }
  }

  /* ==== Scrolling / UX ==== */

  @HostListener('window:scroll')
  onWindowScroll() {
    const y = window.scrollY || document.documentElement.scrollTop || 0;
    this.showScrollTop = y > 400;
  }

  private scrollBehavior(): ScrollBehavior {
    return window.matchMedia('(prefers-reduced-motion: reduce)').matches
      ? 'auto'
      : 'smooth';
  }

  scrollToTop(): void {
    window.scrollTo({ top: 0, behavior: this.scrollBehavior() });
  }

scrollToHistory(): void {
  setTimeout(() => {
    document
      .getElementById('history-section')
      ?.scrollIntoView({ behavior: this.scrollBehavior() });
  }, 50);
}

scrollToHistoryTop(): void {
  setTimeout(() => {
    document
      .getElementById('history-section')
      ?.scrollIntoView({ behavior: this.scrollBehavior(), block: 'start' });
  }, 50);
}

  /* ==== Threshold (debounced) ==== */
  updateThreshold(): void {
    if (this.updateThresholdTimeout !== undefined) {
      window.clearTimeout(this.updateThresholdTimeout);
    }
    this.updateThresholdTimeout = window.setTimeout(() => {
      this.http
        .post('/api/threshold', { threshold: this.threshold })
        .pipe(takeUntil(this.destroy$))
        .subscribe({
          next: () => console.log('Threshold updated to', this.threshold),
          error: (err) => console.error('Failed to update threshold', err),
        });
    }, DEBOUNCE_DELAY);
  }

  /* ==== CSV Export ==== */
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
      `"${data.lastChecked
        ? new Date(data.lastChecked).toLocaleString('en-US', {
            dateStyle: 'medium',
            timeStyle: 'short',
          })
        : 'N/A'}"`,
    ]);
    const csv = [headers.join(','), ...rows.map((r) => r.join(','))].join('\n');
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' });
    const url = URL.createObjectURL(blob);

    const a = document.createElement('a');
    a.href = url;
    a.download = `host_status_${new Date().toISOString().slice(0, 10)}.csv`;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
  }

  /* ==== Stats helpers ==== */
  getUpHostsCount(): number {
    return this.hostEntries.filter(([_, d]) => d.up).length;
  }

  getDownHostsCount(): number {
    return this.hostEntries.filter(([_, d]) => !d.up).length;
  }

  /* ==== Chart config ==== */
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
  animation: false
};


  @ViewChild(BaseChartDirective) chart?: BaseChartDirective<'pie'>;

  private updateChartData(): void {
    const up = this.getUpHostsCount();
    const down = this.getDownHostsCount();
    const current = this.pieChartData.datasets[0].data as number[];
    if (current[0] !== up || current[1] !== down) {
      this.pieChartData.datasets[0].data = [up, down];
      this.chart?.update();
    }
  }

  /* ==== ngFor helpers ==== */
  trackByHost = (index: number, item: [string, HostMetrics]) => item[0];

  /* ==== Equality / Ordering ==== */
  private normalizeEntries(entries: Array<[string, HostMetrics]>) {
    return entries.sort((a, b) => a[0].localeCompare(b[0]));
  }

  private areHostEntriesEqual(
    a: Array<[string, HostMetrics]>,
    b: Array<[string, HostMetrics]>
  ): boolean {
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
