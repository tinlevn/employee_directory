import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideZonelessChangeDetection } from '@angular/core';
import { of, throwError } from 'rxjs';

import { Employee } from '../../core/models/employee.model';
import { EmployeeService } from '../../core/services/employee.service';
import { SearchResultComponent } from './searchresult.component';

const testEmployees: Employee[] = [
  {
    id: '1',
    firstName: 'Ada',
    lastName: 'Lovelace',
    title: 'Engineer',
    extension: '1234',
    phone: '555-0100',
    location: 'Chelsea',
    department: 'Engineering',
    email: 'ada@example.com',
    hiredDate: '2019-04-12T00:00:00'
  }
];

describe('SearchResultComponent', () => {
  let component: SearchResultComponent;
  let fixture: ComponentFixture<SearchResultComponent>;
  let employeeService: jasmine.SpyObj<EmployeeService>;

  beforeEach(async () => {
    employeeService = jasmine.createSpyObj<EmployeeService>('EmployeeService', ['search', 'getById']);
    employeeService.search.and.returnValue(of(testEmployees));

    await TestBed.configureTestingModule({
      imports: [SearchResultComponent],
      providers: [
        provideZonelessChangeDetection(),
        { provide: EmployeeService, useValue: employeeService }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(SearchResultComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load employees on init and mark the directory as ready', () => {
    expect(employeeService.search).toHaveBeenCalledWith(component.searchCriteria());
    expect(component.employeeList().length).toBe(1);
    expect(component.hasLoaded()).toBeTrue();
    expect(component.hasError()).toBeFalse();
  });

  it('should expose an error state when the search fails', () => {
    employeeService.search.and.returnValue(throwError(() => new Error('boom')));

    component.onSubmit();

    expect(component.hasError()).toBeTrue();
    expect(component.hasLoaded()).toBeFalse();
    expect(component.employeeList().length).toBe(0);
  });

  it('should report active filters only when a criterion has a value', () => {
    expect(component.hasActiveFilters()).toBeFalse();

    component.updateCriteria('department', 'Engineering');

    expect(component.hasActiveFilters()).toBeTrue();
  });

  it('should reset criteria and reload', () => {
    component.updateCriteria('firstName', 'Ada');
    employeeService.search.calls.reset();

    component.reset();

    expect(component.searchCriteria().firstName).toBe('');
    expect(employeeService.search).toHaveBeenCalled();
  });
});
