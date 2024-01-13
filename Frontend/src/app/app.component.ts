import { Component, ViewChild, ViewEncapsulation } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';
import { IdleScreenComponent } from './idle-screen/idle-screen.component';
import { DetectPageComponent } from './detect-page/detect-page.component';
import { IntroComponent } from './intro/intro.component';
import { ClassifyComponent } from './classify/classify.component';
import { RewardingComponent } from './rewarding/rewarding.component';

@Component({
    selector: 'app-root',
    standalone: true,
    imports: [CommonModule, RouterOutlet, IdleScreenComponent, DetectPageComponent, IntroComponent, ClassifyComponent, RewardingComponent],
    templateUrl: './app.component.html',
    styleUrl: './app.component.css',
    encapsulation: ViewEncapsulation.None
})

export class AppComponent {
    state:string = "reward";
    private socket: WebSocket;
    classified = 3;

    public constructor() {
        // idle -> to idle screen
        // intro -> to start screen
        // info:class -> 
        // reward -> to awarding screen
        this.socket = new WebSocket("ws://localhost:8080/ws");
        this.socket.onopen = event => { }
        this.socket.onclose = event => { }
        this.socket.onmessage = event => {
            console.log(event.data)
            var data = event.data.split(":"); 
            if (data[0] == "info") {
                this.classified = +data[1]
            }
            this.nextComponent(data[0]);
        }
    }
    
    ngInit() : void {
        this.nextComponent(this.state);
    }

    nextComponent(state : string) : void {
        this.state = state;
    }
}


