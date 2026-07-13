import { CommonModule } from '@angular/common';
import { Component, OnInit, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';
import { catchError, finalize, tap } from 'rxjs/operators';
import { of } from 'rxjs';

import { Employee, StaffSearchCriteria } from '../../core/models/employee.model';
import { EmployeeService } from '../../core/services/employee.service';

@Component({
  selector: 'app-searchresult',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    MatFormFieldModule,
    MatIconModule,
    MatInputModule,
    MatPaginatorModule,
    MatSelectModule,
    MatSortModule,
    MatTableModule
  ],
  templateUrl: './searchresult.component.html',
  styleUrls: ['./searchresult.component.scss']
})
export class SearchResultComponent implements OnInit {
  departments = ['', 'Marketing', 'Sales', 'Product', 'Engineering', 'IT', 'Human Resources', 'Finance', 'Data Science', 'Design'];
  locations = ['', 'Deer Island', 'Chelsea'];
  pageSizes = [25, 50, 100, 250];
  displayedColumns: string[] = ['lastName', 'title', 'extension', 'phone', 'location', 'department'];

  searchCriteria: StaffSearchCriteria = {
    firstName: '',
    lastName: '',
    jobTitle: '',
    department: '',
    location: ''
  };

  employeeList: Employee[] = [];
  dataSource = new MatTableDataSource<Employee>();
  isLoading = false;
  hasLoaded = false;
  errorMessage: string | null = null;
  currentFilterValue = '';
  totalCount = 0;

  @ViewChild(MatSort, { static: false }) resultSort!: MatSort;
  @ViewChild(MatPaginator, { static: false }) paginator!: MatPaginator;

  constructor(private readonly service: EmployeeService) { }

  ngOnInit(): void {
    this.loadEmployees();
  }

  ngAfterViewInit(): void {
    this.dataSource.sort = this.resultSort;
    this.dataSource.paginator = this.paginator;
    this.dataSource.filterPredicate = (employee: Employee, filter: string) => {
      const normalizedFilter = filter.trim().toLowerCase();
      if (!normalizedFilter) {
        return true;
      }

      const searchableText = [
        employee.firstName,
        employee.lastName,
        employee.title,
        employee.department,
        employee.location,
        employee.phone,
        employee.extension
      ]
        .filter(Boolean)
        .join(' ')
        .toLowerCase();

      return searchableText.includes(normalizedFilter);
    };
  }

  get hasActiveFilters(): boolean {
    return Object.values(this.searchCriteria).some((value) => value.trim().length > 0);
  }

  applyFilter(event: Event): void {
    const filterValue = (event.target as HTMLInputElement).value;
    this.currentFilterValue = filterValue;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  clearFilter(): void {
    this.searchCriteria = {
      firstName: '',
      lastName: '',
      jobTitle: '',
      department: '',
      location: ''
    };
    this.currentFilterValue = '';
    this.dataSource.filter = '';
    this.loadEmployees();
  }

  onSubmit(): void {
    this.loadEmployees();
  }

  private loadEmployees(): void {
    this.isLoading = true;
    this.errorMessage = null;

    const request$ = this.hasActiveFilters
      ? this.service.findEmployees(this.searchCriteria)
      : this.service.getAllEmployees();

    request$.pipe(
      tap((data) => {
        this.employeeList = [...data];
        this.totalCount = data.length;
        this.hasLoaded = true;
        this.syncDataSource();
      }),
      catchError(() => {
        this.errorMessage = 'We could not fetch the staff directory right now.';
        this.employeeList = [];
        this.totalCount = 0;
        this.hasLoaded = false;
        this.syncDataSource();
        return of([]);
      }),
      finalize(() => {
        this.isLoading = false;
      })
    ).subscribe();
  }

  private syncDataSource(): void {
    this.dataSource.data = this.employeeList;
    this.dataSource.filter = this.currentFilterValue.trim().toLowerCase();

    if (this.paginator) {
      this.paginator.firstPage();
      this.dataSource.paginator = this.paginator;
    }

    if (this.resultSort) {
      this.dataSource.sort = this.resultSort;
    }
  }
}

