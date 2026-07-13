import { NO_ERRORS_SCHEMA } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';

import { EmployeeService } from '../../core/services/employee.service';
import { SearchResultComponent } from './searchresult.component';

describe('SearchresultComponent', () => {
  let component: SearchResultComponent;
  let fixture: ComponentFixture<SearchResultComponent>;
  let employeeService: jasmine.SpyObj<EmployeeService>;

  beforeEach(async () => {
    employeeService = jasmine.createSpyObj<EmployeeService>('EmployeeService', ['getAllEmployees', 'findEmployees']);
    employeeService.getAllEmployees.and.returnValue(of([
      {
        firstName: 'Ada',
        lastName: 'Lovelace',
        title: 'Engineer',
        extension: '1234',
        phone: '555-0100',
        location: 'Chelsea',
        department: 'Engineering',
        email: 'ada@example.com'
      }
    ]));

    await TestBed.configureTestingModule({
      declarations: [SearchResultComponent],
      providers: [{ provide: EmployeeService, useValue: employeeService }],
      schemas: [NO_ERRORS_SCHEMA]
    }).compileComponents();

    fixture = TestBed.createComponent(SearchResultComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load employees and mark the directory as ready', () => {
    component.onSubmit();

    expect(employeeService.getAllEmployees).toHaveBeenCalled();
    expect(component.hasLoaded).toBeTrue();
    expect(component.errorMessage).toBeNull();
    expect(component.employeeList.length).toBe(1);
  });
});
