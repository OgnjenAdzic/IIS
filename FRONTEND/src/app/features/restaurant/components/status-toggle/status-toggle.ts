import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-status-toggle',
  imports: [CommonModule],
  templateUrl: './status-toggle.html',
  styleUrl: './status-toggle.css',
})
export class StatusToggle {
  @Input() isOpen: boolean = false;
  @Output() statusChanged = new EventEmitter<boolean>();

  toggle() {
    this.statusChanged.emit(!this.isOpen);
  }
}
