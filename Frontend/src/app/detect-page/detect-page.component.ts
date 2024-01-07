import { Component, EventEmitter, Output, ViewEncapsulation } from '@angular/core';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-detect-page',
  standalone: true,
  imports: [],
  templateUrl: './detect-page.component.html',
  styleUrl: './detect-page.component.css',
  encapsulation: ViewEncapsulation.None
})
export class DetectPageComponent implements BaseComponent {
    @Output() switch = new EventEmitter<void>();
    
    currentImage: string = '';

    constructor(private api: ApiService) {}

    init(): void {
        this.api.callAPI('/classify').subscribe(
            (response) => {
                this.switch.emit()
            }
        )
    }

    
}
