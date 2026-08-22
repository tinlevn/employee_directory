using Microsoft.EntityFrameworkCore;
using StaffSearch.Dtos;
using StaffSearch.Entities;
using StaffSearch.Interfaces;

namespace StaffSearch.Data.Repositories;

public class EmployeeRepository(StaffSearchDbContext context) : IEmployeeRepository
{
    public async Task<IReadOnlyList<Employee>> GetAllAsync(CancellationToken cancellationToken = default)
        => await context.Employees
            .AsNoTracking()
            .OrderBy(e => e.LastName)
            .ThenBy(e => e.FirstName)
            .ToListAsync(cancellationToken);

    public async Task<Employee?> GetByIdAsync(string id, CancellationToken cancellationToken = default)
        => await context.Employees
            .AsNoTracking()
            .FirstOrDefaultAsync(e => e.Guid == id, cancellationToken);

    public async Task<IReadOnlyList<Employee>> SearchAsync(EmployeeQuery query, CancellationToken cancellationToken = default)
    {
        IQueryable<Employee> employees = context.Employees.AsNoTracking();

        if (!string.IsNullOrWhiteSpace(query.FirstName))
            employees = employees.Where(e => e.FirstName != null && e.FirstName.Contains(query.FirstName.Trim()));

        if (!string.IsNullOrWhiteSpace(query.LastName))
            employees = employees.Where(e => e.LastName != null && e.LastName.Contains(query.LastName.Trim()));

        if (!string.IsNullOrWhiteSpace(query.JobTitle))
            employees = employees.Where(e => e.Title != null && e.Title.Contains(query.JobTitle.Trim()));

        if (!string.IsNullOrWhiteSpace(query.Department))
            employees = employees.Where(e => e.Department != null && e.Department == query.Department.Trim());

        if (!string.IsNullOrWhiteSpace(query.Location))
            employees = employees.Where(e => e.Location != null && e.Location == query.Location.Trim());

        return await employees
            .OrderBy(e => e.LastName)
            .ThenBy(e => e.FirstName)
            .ToListAsync(cancellationToken);
    }
}
