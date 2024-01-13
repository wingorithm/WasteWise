import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { retryWhen } from 'rxjs';

@Injectable({
    providedIn: 'root'
})

export class ApiService {
    // private baseUrl = 'http://127.0.0.1:5000'
    private baseUrl = 'http://127.0.0.1:8080'

    constructor(private http: HttpClient) { }

    callAPI(relativePath : String) : Observable<any> {
        return this.http.get(this.baseUrl + relativePath).pipe(
            retryWhen(errors => errors)
        );
    }
}
