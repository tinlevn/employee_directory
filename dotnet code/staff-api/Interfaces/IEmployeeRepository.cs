using StaffSearch.Dtos;
using StaffSearch.Entities;

namespace StaffSearch.Interfaces;

public interface IEmployeeRepository
{
    Task<IReadOnlyList<Employee>> GetAllAsync(CancellationToken cancellationToken = default);

    Task<Employee?> GetByIdAsync(string id, CancellationToken cancellationToken = default);

    Task<IReadOnlyList<Employee>> SearchAsync(EmployeeQuery query, CancellationToken cancellationToken = default);
}
