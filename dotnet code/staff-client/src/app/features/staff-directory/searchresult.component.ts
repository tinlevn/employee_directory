import { CommonModule } from '@angular/common';
import { Component, DestroyRef, OnInit, ViewChild, computed, effect, inject, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatPaginator, MatPaginatorModule } from '@angular/material/paginator';
import { MatSelectModule } from '@angular/material/select';
import { MatSort, MatSortModule } from '@angular/material/sort';
import { MatTableDataSource, MatTableModule } from '@angular/material/table';

import { EMPTY_SEARCH_CRITERIA, Employee, StaffSearchCriteria } from '../../core/models/employee.model';
import { EmployeeService } from '../../core/services/employee.service';
import { STAFF_DIRECTORY_CONSTANTS } from '../../core/constants/staff-directory.constants';

@Component({
  selector: 'app-searchresult',
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
  private readonly employeeService = inject(EmployeeService);
  private readonly destroyRef = inject(DestroyRef);

  readonly departments = STAFF_DIRECTORY_CONSTANTS.DEPARTMENTS;
  readonly locations = STAFF_DIRECTORY_CONSTANTS.LOCATIONS;
  readonly pageSizes = STAFF_DIRECTORY_CONSTANTS.PAGE_SIZES;
  readonly displayedColumns = STAFF_DIRECTORY_CONSTANTS.DISPLAY_COLUMNS;
  readonly defaultPageSize = STAFF_DIRECTORY_CONSTANTS.DEFAULT_PAGE_SIZE;
  readonly errorText = STAFF_DIRECTORY_CONSTANTS.ERROR_MESSAGE;

  readonly searchCriteria = signal<StaffSearchCriteria>({ ...EMPTY_SEARCH_CRITERIA });

  readonly employeeList = signal<Employee[]>([]);
  readonly isLoading = signal(false);
  readonly hasLoaded = signal(false);
  readonly hasError = signal(false);

  readonly hasActiveFilters = computed(() =>
    Object.values(this.searchCriteria()).some((value) => value.trim().length > 0)
  );

  readonly dataSource = new MatTableDataSource<Employee>([]);

  @ViewChild(MatSort) private set sort(sort: MatSort) {
    this.dataSource.sort = sort;
  }

  @ViewChild(MatPaginator) private set paginator(paginator: MatPaginator) {
    this.dataSource.paginator = paginator;
  }

  constructor() {
    this.dataSource.filterPredicate = this.createFilterPredicate();

    // Keep the Material table in sync with the reactive employee list.
    effect(() => {
      this.dataSource.data = this.employeeList();
    });
  }

  ngOnInit(): void {
    this.loadEmployees();
  }

  applyQuickFilter(event: Event): void {
    this.dataSource.filter = (event.target as HTMLInputElement).value.trim().toLowerCase();
    this.dataSource.paginator?.firstPage();
  }

  updateCriteria(field: keyof StaffSearchCriteria, value: string): void {
    this.searchCriteria.update((criteria) => ({ ...criteria, [field]: value }));
  }

  reset(): void {
    this.searchCriteria.set({ ...EMPTY_SEARCH_CRITERIA });
    this.dataSource.filter = '';
    this.loadEmployees();
  }

  onSubmit(): void {
    this.loadEmployees();
  }

  private loadEmployees(): void {
    this.isLoading.set(true);
    this.hasError.set(false);

    this.employeeService
      .search(this.searchCriteria())
      .pipe(takeUntilDestroyed(this.destroyRef))
      .subscribe({
        next: (employees) => {
          this.employeeList.set(employees);
          this.hasLoaded.set(true);
          this.dataSource.paginator?.firstPage();
          this.isLoading.set(false);
        },
        error: () => {
          this.employeeList.set([]);
          this.hasLoaded.set(false);
          this.hasError.set(true);
          this.isLoading.set(false);
        }
      });
  }

  private createFilterPredicate() {
    return (employee: Employee, filter: string): boolean => {
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
        .filter((value): value is string => !!value)
        .join(' ')
        .toLowerCase();

      return searchableText.includes(normalizedFilter);
    };
  }
}
