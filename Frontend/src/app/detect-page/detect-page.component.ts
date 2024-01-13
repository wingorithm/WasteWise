import { Component, EventEmitter, Output, ViewChild, ViewEncapsulation } from '@angular/core';
import { ApiService } from '../api.service';
import { WebcamImage, WebcamInitError, WebcamModule, WebcamUtil } from 'ngx-webcam';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-detect-page',
  standalone: true,
  imports: [WebcamModule, CommonModule],
  templateUrl: './detect-page.component.html',
  styleUrl: './detect-page.component.css',
  encapsulation: ViewEncapsulation.None
})
export class DetectPageComponent {
    @Output() switch = new EventEmitter<void>();

    started: string = "";
    chat: string = "";
    showCamera = false;
    allowCameraSwitch = false

    constructor(private api: ApiService) {}

    ngOnInit(): void {
        this.showCamera = true;
        this.chat = "Sekarang periksa yuk sampah kamu itu jenis apa"

        this.api.callAPI('/classify').subscribe(
            (response) => {
                this.switch.emit()
                this.started = "";
                this.showCamera = false;
            }
        )
    }

}