import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-menu-list',
  imports: [CommonModule],
  templateUrl: './menu-list.html',
  styleUrl: './menu-list.css',
})
export class MenuList {
  @Input() items: any[] = [];
  @Output() priceUpdated = new EventEmitter<{ id: string, price: number }>();
  @Output() itemDeleted = new EventEmitter<string>();

  onPriceChange(id: string, event: Event) {
    const input = event.target as HTMLInputElement;
    const newPrice = parseFloat(input.value);

    if (newPrice > 0) {
      this.priceUpdated.emit({ id, price: newPrice });
    }
  }

  onDelete(id: string) {
    if (confirm('Are you sure you want to remove this item from the menu?')) {
      this.itemDeleted.emit(id);
    }
  }
}
