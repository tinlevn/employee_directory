import { TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection } from '@angular/core';
import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';

import { environment } from '../../../environments/environment';
import { Employee } from '../models/employee.model';
import { EmployeeService } from './employee.service';

const mockEmployee: Employee = {
  id: 'id-ada',
  firstName: 'Ada',
  lastName: 'Lovelace',
  title: 'Principal Engineer',
  extension: '1001',
  phone: '617-555-0101',
  location: 'Chelsea',
  department: 'Engineering',
  email: 'ada.lovelace@example.com',
  hiredDate: '2019-04-12T00:00:00'
};

describe('EmployeeService', () => {
  let service: EmployeeService;
  let http: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        provideZonelessChangeDetection(),
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });

    service = TestBed.inject(EmployeeService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('should GET the collection without params when all criteria are blank', () => {
    service
      .search({ firstName: '', lastName: '  ', jobTitle: '', department: '', location: '' })
      .subscribe((employees) => expect(employees).toEqual([mockEmployee]));

    const req = http.expectOne(environment.apiUrl);
    expect(req.request.method).toBe('GET');
    expect(req.request.params.keys().length).toBe(0);
    req.flush([mockEmployee]);
  });

  it('should map active criteria to query params and trim values', () => {
    service
      .search({ firstName: '  Ada ', lastName: '', jobTitle: 'Engineer', department: 'Engineering', location: '' })
      .subscribe();

    const req = http.expectOne((r) =>
      r.url === environment.apiUrl
      && r.params.get('firstName') === 'Ada'
      && r.params.get('jobTitle') === 'Engineer'
      && r.params.get('department') === 'Engineering'
      && !r.params.has('lastName')
      && !r.params.has('location'));

    expect(req.request.method).toBe('GET');
    req.flush([]);
  });

  it('should GET a single employee by id with encoding', () => {
    service.getById('id with spaces').subscribe((employee) => expect(employee).toEqual(mockEmployee));

    const req = http.expectOne(`${environment.apiUrl}/id%20with%20spaces`);
    expect(req.request.method).toBe('GET');
    req.flush(mockEmployee);
  });

  it('should propagate HTTP errors to the caller', () => {
    let errored = false;

    service.search({ firstName: '', lastName: '', jobTitle: '', department: '', location: '' })
      .subscribe({ error: () => (errored = true) });

    http.expectOne(environment.apiUrl).flush('boom', { status: 500, statusText: 'Server Error' });

    expect(errored).toBeTrue();
  });
});
