import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

import { environment } from '../../../environments/environment';
import { Employee, StaffSearchCriteria } from '../models/employee.model';
import { mockEmployees } from '../../features/staff-directory/mock-employees';

@Injectable({
  providedIn: 'root'
})
export class EmployeeService {
  private readonly baseUrl = environment.apiUrl;

  constructor(private readonly http: HttpClient) {}

  findEmployees(criteria: StaffSearchCriteria): Observable<Employee[]> {
    return this.http.post<Employee[]>(`${this.baseUrl}/filter`, criteria, {
      withCredentials: true
    }).pipe(
      catchError(() => of(this.filterMockEmployees(criteria)))
    );
  }

  getAllEmployees(): Observable<Employee[]> {
    return this.http.get<Employee[]>(this.baseUrl, {
      withCredentials: true
    }).pipe(
      catchError(() => of(mockEmployees))
    );
  }

  private filterMockEmployees(criteria: StaffSearchCriteria): Employee[] {
    return mockEmployees.filter((employee) => {
      const matchesFirstName = !criteria.firstName || employee.firstName.toLowerCase().includes(criteria.firstName.toLowerCase());
      const matchesLastName = !criteria.lastName || employee.lastName.toLowerCase().includes(criteria.lastName.toLowerCase());
      const matchesJobTitle = !criteria.jobTitle || employee.title.toLowerCase().includes(criteria.jobTitle.toLowerCase());
      const matchesDepartment = !criteria.department || employee.department.toLowerCase() === criteria.department.toLowerCase();
      const matchesLocation = !criteria.location || employee.location.toLowerCase() === criteria.location.toLowerCase();

      return matchesFirstName && matchesLastName && matchesJobTitle && matchesDepartment && matchesLocation;
    });
  }
}
