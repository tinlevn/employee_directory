using System.ComponentModel.DataAnnotations;

namespace StaffSearch.Dtos;

/// <summary>
/// Optional search criteria for querying employees. All values are matched
/// case-insensitively using "contains" semantics; omitted values are ignored.
/// </summary>
public class EmployeeQuery
{
    private const int MaxSearchLength = 100;

    [StringLength(MaxSearchLength)]
    public string? FirstName { get; set; }

    [StringLength(MaxSearchLength)]
    public string? LastName { get; set; }

    [StringLength(MaxSearchLength)]
    public string? JobTitle { get; set; }

    [StringLength(MaxSearchLength)]
    public string? Department { get; set; }

    [StringLength(MaxSearchLength)]
    public string? Location { get; set; }
}
