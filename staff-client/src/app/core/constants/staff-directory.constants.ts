/**
 * Constants for the staff directory feature
 */
export const STAFF_DIRECTORY_CONSTANTS = {
  DEPARTMENTS: [
    'Marketing',
    'Sales',
    'Product',
    'Engineering',
    'IT',
    'Human Resources',
    'Finance',
    'Data Science',
    'Design'
  ] as const,
  LOCATIONS: ['Deer Island', 'Chelsea'] as const,
  PAGE_SIZES: [25, 50, 100, 250] as const,
  DEFAULT_PAGE_SIZE: 50,
  DISPLAY_COLUMNS: ['lastName', 'title', 'extension', 'phone', 'location', 'department'] as const,
  ERROR_MESSAGE: 'We could not fetch the staff directory right now.'
};

export type Department = (typeof STAFF_DIRECTORY_CONSTANTS.DEPARTMENTS)[number];
export type Location = (typeof STAFF_DIRECTORY_CONSTANTS.LOCATIONS)[number];
