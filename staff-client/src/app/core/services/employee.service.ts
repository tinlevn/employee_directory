import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';

import { environment } from '../../../environments/environment';
import { Employee, StaffSearchCriteria } from '../models/employee.model';

@Injectable({ providedIn: 'root' })
export class EmployeeService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = environment.apiUrl;

  /**
   * GET /api/employees with the active criteria as query-string parameters.
   * Blank criteria are omitted, so an empty form returns the full directory.
   */
  search(criteria: StaffSearchCriteria): Observable<Employee[]> {
    return this.http.get<Employee[]>(this.baseUrl, { params: this.toParams(criteria) });
  }

  getById(id: string): Observable<Employee> {
    return this.http.get<Employee>(`${this.baseUrl}/${encodeURIComponent(id)}`);
  }

  private toParams(criteria: StaffSearchCriteria): HttpParams {
    let params = new HttpParams();

    for (const [key, value] of Object.entries(criteria)) {
      const trimmed = value.trim();
      if (trimmed) {
        params = params.set(key, trimmed);
      }
    }

    return params;
  }
}
