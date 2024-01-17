import { Component, EventEmitter, Output, ViewChild, ViewEncapsulation } from '@angular/core';
import { WebcamImage, WebcamInitError, WebcamModule, WebcamUtil } from 'ngx-webcam';
import { CommonModule } from '@angular/common';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Subject, tap } from 'rxjs';

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

    triggerObservable = new Subject<void>();

    started: string = "";
    chat: string = "";
    showCamera = false;
    allowCameraSwitch = false
    captureImageData = true;

    handle! : any

    ngOnInit(): void {
        this.showCamera = true;
        this.chat = "Sekarang periksa yuk sampah kamu itu jenis apa"
        this.handle = setInterval(() => {
            this.triggerObservable.next()
        }, 3000);
        // 1000
    }

    ngOnDestroy(): void {
        clearInterval(this.handle)
    }

    constructor(private http: HttpClient) { }

    handleImage(image: WebcamImage): void {

        const imageData = image.imageAsBase64;
        console.log(imageData);
        
        const apiUrl = 'http://localhost:8080/imageHandler';

        const headers = new HttpHeaders()
            .set('Content-Type', 'application/octet-stream')
            .set('X-Content-Transfer-Encoding', 'base64'); 
        // Set the image data in the request body

        // Set the headers
        // const headers = new HttpHeaders({
        //     'Content-Type': 'application/json',
        // });

        this.http.post(apiUrl, imageData, { headers: headers }).subscribe(val => console.log(val));
    }
}