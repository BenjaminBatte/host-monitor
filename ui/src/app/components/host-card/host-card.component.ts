import { Component, Input } from '@angular/core';
import { CommonModule, NgClass, DecimalPipe, DatePipe } from '@angular/common';

@Component({
  selector: 'app-host-card',
  standalone: true,
  imports: [CommonModule, NgClass, DecimalPipe, DatePipe],
  templateUrl: './host-card.component.html',
  styleUrls: ['./host-card.component.scss'],
})
export class HostCardComponent {
  @Input() host!: string;
  @Input() data!: any;

  toastMessage = '';  

  copyToClipboard(value: string): void {
    navigator.clipboard.writeText(value).then(() => {
      this.toastMessage = 'Copied to clipboard!';
      setTimeout(() => this.toastMessage = '', 2000);
    });
  }
}
