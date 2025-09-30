import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Employee } from '../_models/employee';
import { environment } from '../../environments/environment';
import { Specifications } from '../_models/spec';
@Injectable({
  providedIn: 'root'
})
export class EmployeeService {

  constructor(private http: HttpClient) { }

  findEmployees(criteria: Specifications) {
    return this.http.post<Employee[]>(environment.apiUrl + '/filter', criteria, {withCredentials: true});
  }
  getAllEmployees() {
    return this.http.get<Employee[]>(environment.apiUrl);
  }
}

/*
Add ", {withCredentials: true}" to post request to allow negotiation */
