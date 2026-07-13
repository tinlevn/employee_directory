export interface Employee {
  firstName: string;
  lastName: string;
  title: string;
  extension: string;
  phone: string;
  location: string;
  department: string;
  email: string;
}

export interface StaffSearchCriteria {
  firstName: string;
  lastName: string;
  jobTitle: string;
  department: string;
  location: string;
}
