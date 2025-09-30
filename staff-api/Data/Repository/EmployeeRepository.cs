using Microsoft.AspNetCore.Mvc;
using Microsoft.EntityFrameworkCore;
using Staff_Search.Entity;
using Staff_Search.Interface;

namespace Staff_Search.Data.Repository
{
    public class EmployeeRepository : IEmployeeRepository
    {
        private readonly staffContext _context;

        public EmployeeRepository(staffContext context)
        {
            _context = context;
        }

        public Task<Employee> CreateEmployee(Employee empToAdd)
        {
            throw new NotImplementedException();
        }

        public Task<IActionResult> DeleteEmployee(int id)
        {
            throw new NotImplementedException();
        }

        public async Task<IReadOnlyList<Employee>> GetEmployeesAsync()
        {
            return await _context.Employees.ToListAsync();
        }

        public Task<Employee> UpdateEmployee(Employee specs)
        {
            throw new NotImplementedException();
        }

        public async Task<IReadOnlyList<Employee>> GetEmployeeResult(Filter specs)
        {
            return await _context.Employees
            .Where(t => t.FirstName!.Contains(specs.firstName!))
            .Where(t => t.LastName!.Contains(specs.lastName!))
            .Where(t => t.Department!.Contains(specs.Department!))
            .Where(t => t.Location!.Contains(specs.Location!))
            .Where(t => t.Title!.Contains(specs.jobTitle!))
            .ToListAsync();
        }
    }
}
