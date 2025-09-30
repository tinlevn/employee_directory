import { Component, OnInit, ViewChild } from '@angular/core';
import { Employee } from '../_models/employee';
import { EmployeeService } from '../_services/employee-service.service';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { MatTableDataSource } from '@angular/material/table';
import { Specifications } from '../_models/spec';

@Component({
  selector: 'app-searchresult',
  templateUrl: './searchresult.component.html',
  styleUrls: ['./searchresult.component.scss']
})
export class SearchResultComponent implements OnInit {

  departments = ['', 'Marketing','Sales', 'Product', 'Engineering', 'IT', 'Human Resources', 'Finance', 'Data Science', 'Design'];

  locations = ['', 'Deer Island', 'Chelsea'];

  pageSizes = [50, 100, 500];

  searchCriteria: Specifications = {
    firstName: '',
    lastName: '',
    jobTitle: '',
    department: '',
    location: ''
  };

  employeeList: Employee[] = [];
  dataSource = new MatTableDataSource<Employee>();
  displayedColumns: string[] = ['lastName', 'title', 'extension', 'phone',
    'location', 'department'];

  @ViewChild(MatSort, { static: false }) resultSort = new MatSort();
  @ViewChild(MatPaginator, { static: false }) paginator!: MatPaginator;

  constructor(private service: EmployeeService) { }

  ngOnInit(): void {

  }

  /*In Angular, this will make table load faster if initial
  loading of dataset is large, i.e 1k+ records.*/
  ngAfterViewInit() {
    this.dataSource.paginator = this.paginator;
    this.dataSource.sort = this.resultSort;
  }

  applyFilter(event: Event) {
    const filterValue = (event.target as HTMLInputElement).value;
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }

  clearFilter() {
    this.searchCriteria = {
      firstName: '',
      lastName: '',
      jobTitle: '',
      department: '',
      location: ''
    };
    this.employeeList = [];
    this.dataSource.data = this.employeeList;
  }

  onSubmit() {
    // this.service.findEmployees(this.searchCriteria).subscribe(
    //   (data) => {
    //     this.employeeList = [...data];
    //     this.dataSource.data = this.employeeList;
    //   },
    //   error => { console.log(error); }
    // )
    this.service.getAllEmployees().subscribe(
      (data) => {
        this.employeeList = [...data];
        this.dataSource.data = this.employeeList;
      },
      error => { console.log(error); }
    )
  }
}

