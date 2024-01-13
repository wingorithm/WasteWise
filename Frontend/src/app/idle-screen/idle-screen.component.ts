import { Component, EventEmitter, Output, ViewEncapsulation } from '@angular/core';
import { ApiService } from '../api.service';

@Component({
  selector: 'app-idle-screen',
  standalone: true,
  imports: [],
  templateUrl: './idle-screen.component.html',
  styleUrl: './idle-screen.component.css',
  encapsulation: ViewEncapsulation.None
})

export class IdleScreenComponent {
    @Output() switch = new EventEmitter<void>();
    
    images: string[] = ['Idle - 1.jpg', 'Idle - 2.jpg'];
    idx = 0;
    currentImage!: string;
    

    constructor(private api: ApiService) {}

    ngOnInit(): void {
        this.idx = 0
        this.nextImage()
        let handle = setInterval(() => {this.nextImage()}, 5000)
        this.api.callAPI('/idle').subscribe(
            (response) => {
                clearInterval(handle);
                this.switch.emit();
            }
        )
    }
    
    nextImage() : void {
        this.currentImage = '/assets/page/' + this.images[this.idx];
        console.log(this.currentImage)
        this.idx = (this.idx + 1) % this.images.length;
    }
}
