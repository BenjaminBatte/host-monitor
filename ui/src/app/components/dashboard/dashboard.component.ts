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
import { ChartData, ChartOptions } from 'chart.js';
import DataLabelsPlugin from 'chartjs-plugin-datalabels';

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

  
  showScrollTop = false;

  private destroy$ = new Subject<void>();
  private updateThresholdTimeout?: number;

  constructor(
    private http: HttpClient,
    private metricsService: MetricsService,
    private cdr: ChangeDetectorRef,
    private zone: NgZone
  ) {}


  ngOnInit(): void {

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


  getUpHostsCount(): number {
    return this.hostEntries.filter(([_, d]) => d.up).length;
  }
  getDownHostsCount(): number {
    return this.hostEntries.filter(([_, d]) => !d.up).length;
  }
  get totalHosts(): number {
    return this.hostEntries.length;
  }

  /* ==== Chart config ==== */
  @ViewChild(BaseChartDirective) chart?: BaseChartDirective; 

  pieChartPlugins = [DataLabelsPlugin];

  pieChartData: ChartData<'doughnut'> = {
    labels: ['Up', 'Down'],
    datasets: [
      {
        data: [0, 0],
        backgroundColor: ['#2ecc71', '#e74c3c'],
        borderWidth: 0,
        hoverOffset: 10,
      },
    ],
  };

  pieChartOptions: ChartOptions<'doughnut'> = {
    responsive: true,
    maintainAspectRatio: false,
    cutout: '65%',
    animation: { duration: 400, easing: 'easeOutQuart' },
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label: (ctx) => {
            const total = (ctx.dataset.data as number[]).reduce((a, b) => a + b, 0) || 1;
            const val = ctx.parsed as number;
            const pct = ((val / total) * 100).toFixed(1);
            return `${ctx.label}: ${val} (${pct}%)`;
          },
        },
      },
      datalabels: {
        color: '#2c3e50',
        formatter: (value, ctx) => {
          const data = ctx.dataset.data as number[];
          const total = data.reduce((a, b) => a + b, 0) || 1;
          const pct = (Number(value) / total) * 100;
          return pct >= 5 ? `${pct.toFixed(0)}%` : '';
        },
   font: { weight: 700, size: 12 },
        clamp: true,
      },
    },
  };

  private updateChartData(): void {
    const up = this.getUpHostsCount();
    const down = this.getDownHostsCount();
    const current = this.pieChartData.datasets[0].data as number[];
    if (current[0] !== up || current[1] !== down) {
  
      (this.pieChartData.datasets[0].data as number[])[0] = up;
      (this.pieChartData.datasets[0].data as number[])[1] = down;
      // trigger redraw
      this.chart?.update();
    }
  }

overallAvailability(): number {
  const totals = this.hostEntries.reduce((acc, [, d]) => {
    acc.s += d.successCount || 0;
    acc.t += d.totalChecks || 0;
    return acc;
  }, {s: 0, t: 0});
  return totals.t ? (totals.s / totals.t) * 100 : 0;
}


overallLossPct(): number {
  const a = this.overallAvailability();
  return 100 - a;
}

avgLatencyMs(): number {
  const vals = this.hostEntries
    .filter(([, d]) => d.up && typeof d.latency === 'number')
    .map(([, d]) => Number(d.latency));
  if (!vals.length) return 0;
  return vals.reduce((a, b) => a + b, 0) / vals.length;
}


topLoss(): Array<[string, number, number]> {
  const items = this.hostEntries.map(([h, d]) => {
    const t = d.totalChecks || 0;
    const s = d.successCount || 0;
    const avail = t ? (s / t) * 100 : 0;
    const loss = 100 - avail;
    return [h, avail, loss] as [string, number, number];
  });
  return items.sort((a, b) => b[2] - a[2]).slice(0, 5);
}


slowestHosts(): Array<[string, number]> {
  const items = this.hostEntries
    .filter(([, d]) => typeof d.latency === 'number' && d.latency !== null)
    .map(([h, d]) => [h, Number(d.latency)] as [string, number]);
  return items.sort((a, b) => b[1] - a[1]).slice(0, 5);
}


  trackByHost = (index: number, item: [string, HostMetrics]) => item[0];

 
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
