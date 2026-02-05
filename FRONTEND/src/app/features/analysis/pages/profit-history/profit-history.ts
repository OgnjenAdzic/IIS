import { Component, inject, OnInit, signal } from '@angular/core';
import { CommonModule, DatePipe, DecimalPipe } from '@angular/common';
import { RouterLink } from '@angular/router';
import { AnalysisService } from '../../services/analysis.service';
import { ProfitLogItem } from '../../models/analysis.model';

@Component({
  selector: 'app-profit-history',
  standalone: true,
  imports: [CommonModule, DatePipe, DecimalPipe, RouterLink],
  templateUrl: './profit-history.html',
  styleUrl: './profit-history.css'
})
export class ProfitHistoryComponent implements OnInit {
  analysisService = inject(AnalysisService);

  logs = signal<ProfitLogItem[]>([]);
  totalRevenue = signal<number>(0);

  ngOnInit() {
    this.analysisService.getHistory().subscribe({
      next: (res) => {
        const data = res.logs || [];
        this.logs.set(data);
        console.log(data)
        const total = data.reduce((acc, item) => acc + item.totalProfit, 0);
        this.totalRevenue.set(total);
      },
      error: (err) => console.error(err)
    });
  }
}