using StaffSearch.Dtos;
using StaffSearch.Entities;

namespace StaffSearch.Extensions;

public static class EmployeeMappingExtensions
{
    public static EmployeeDto ToDto(this Employee employee) => new(
        Id: employee.Guid,
        FirstName: employee.FirstName,
        LastName: employee.LastName,
        Title: employee.Title,
        Extension: employee.Extension,
        Phone: employee.Phone,
        Location: employee.Location,
        Department: employee.Department,
        Email: employee.Email,
        HiredDate: employee.HiredDate);

    public static IReadOnlyList<EmployeeDto> ToDtos(this IEnumerable<Employee> employees)
        => employees.Select(ToDto).ToList();
}
