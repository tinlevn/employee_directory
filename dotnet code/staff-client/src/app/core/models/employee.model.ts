export interface Employee {
  id: string;
  firstName: string | null;
  lastName: string | null;
  title: string | null;
  extension: string | null;
  phone: string | null;
  location: string | null;
  department: string | null;
  email: string | null;
  hiredDate: string | null;
}

export interface StaffSearchCriteria {
  firstName: string;
  lastName: string;
  jobTitle: string;
  department: string;
  location: string;
}

export const EMPTY_SEARCH_CRITERIA: StaffSearchCriteria = {
  firstName: '',
  lastName: '',
  jobTitle: '',
  department: '',
  location: ''
};
