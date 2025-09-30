using Microsoft.AspNetCore.Mvc;
using Staff_Search.Entity;

namespace Staff_Search.Interface
{
    public interface IEmployeeRepository
    {
        Task<IReadOnlyList<Employee>> GetEmployeesAsync();

        Task<Employee> CreateEmployee(Employee empToAdd);

        Task<Employee> UpdateEmployee(Employee specs);

        Task<IActionResult> DeleteEmployee(int id);

        Task<IReadOnlyList<Employee>> GetEmployeeResult(Filter specifications);

    }
}
