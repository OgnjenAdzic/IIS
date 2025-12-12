import { Component, EventEmitter, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-add-item-form',
  imports: [FormsModule, CommonModule],
  templateUrl: './add-item-form.html',
  styleUrl: './add-item-form.css',
})
export class AddItemForm {
  @Output() itemAdded = new EventEmitter<{ name: string, price: number }>();

  name: string = '';
  price: number | null = null;

  onSubmit() {
    if (this.name && this.price && this.price > 0) {
      this.itemAdded.emit({ name: this.name, price: this.price });
      // Reset form
      this.name = '';
      this.price = null;
    }
  }
}
