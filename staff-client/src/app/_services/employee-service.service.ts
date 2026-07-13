import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

import { Employee } from '../_models/employee';
import { environment } from '../../environments/environment';
import { Specifications } from '../_models/spec';

@Injectable({
  providedIn: 'root'
})
export class EmployeeService {
  private readonly baseUrl = environment.apiUrl;

  constructor(private readonly http: HttpClient) { }

  findEmployees(criteria: Specifications): Observable<Employee[]> {
    return this.http.post<Employee[]>(`${this.baseUrl}/filter`, criteria, {
      withCredentials: true
    });
  }

  getAllEmployees(): Observable<Employee[]> {
    return this.http.get<Employee[]>(this.baseUrl, {
      withCredentials: true
    });
  }
}
