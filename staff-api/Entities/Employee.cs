using System.ComponentModel.DataAnnotations;

namespace StaffSearch.Entities;

/// <summary>
/// Persistence model that maps to the Employees table in the staff database.
/// </summary>
public class Employee
{
    [Key]
    [MaxLength(64)]
    public string Guid { get; set; } = string.Empty;

    [MaxLength(100)]
    public string? FirstName { get; set; }

    [MaxLength(100)]
    public string? LastName { get; set; }

    [MaxLength(150)]
    public string? Title { get; set; }

    [MaxLength(25)]
    public string? Extension { get; set; }

    [MaxLength(25)]
    public string? Phone { get; set; }

    [MaxLength(100)]
    public string? Location { get; set; }

    [MaxLength(100)]
    public string? Department { get; set; }

    [MaxLength(256)]
    public string? Email { get; set; }

    public DateTime? HiredDate { get; set; }
}
