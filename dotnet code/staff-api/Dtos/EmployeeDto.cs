namespace StaffSearch.Dtos;

/// <summary>
/// Public representation of an employee returned by the API.
/// </summary>
/// <param name="Id">Stable identifier of the employee record.</param>
public record EmployeeDto(
    string Id,
    string? FirstName,
    string? LastName,
    string? Title,
    string? Extension,
    string? Phone,
    string? Location,
    string? Department,
    string? Email,
    DateTime? HiredDate);
