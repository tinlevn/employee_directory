using Microsoft.EntityFrameworkCore;
using StaffSearch.Data;
using StaffSearch.Entities;

namespace StaffSearch.Tests.Helpers;

internal static class TestDbContextFactory
{
    /// <summary>
    /// Creates an isolated in-memory context (unique database per call) seeded
    /// with a predictable set of employees.
    /// </summary>
    public static StaffSearchDbContext CreateSeeded()
    {
        var options = new DbContextOptionsBuilder<StaffSearchDbContext>()
            .UseInMemoryDatabase(Guid.NewGuid().ToString())
            .Options;

        var context = new StaffSearchDbContext(options);
        context.Employees.AddRange(TestEmployees.All);
        context.SaveChanges();
        return context;
    }
}

internal static class TestEmployees
{
    public static readonly Employee[] All =
    [
        new Employee
        {
            Guid = "id-ada",
            FirstName = "Ada",
            LastName = "Lovelace",
            Title = "Principal Engineer",
            Extension = "1001",
            Phone = "617-555-0101",
            Location = "Chelsea",
            Department = "Engineering",
            Email = "ada.lovelace@example.com",
            HiredDate = new DateTime(2019, 4, 12)
        },
        new Employee
        {
            Guid = "id-grace",
            FirstName = "Grace",
            LastName = "Hopper",
            Title = "Engineering Manager",
            Extension = "1002",
            Phone = "617-555-0102",
            Location = "Chelsea",
            Department = "Engineering",
            Email = "grace.hopper@example.com",
            HiredDate = new DateTime(2017, 8, 3)
        },
        new Employee
        {
            Guid = "id-alan",
            FirstName = "Alan",
            LastName = "Turing",
            Title = "Data Scientist",
            Extension = "1003",
            Phone = "617-555-0103",
            Location = "Deer Island",
            Department = "Data Science",
            Email = "alan.turing@example.com",
            HiredDate = new DateTime(2021, 1, 25)
        },
        new Employee
        {
            Guid = "id-katherine",
            FirstName = "Katherine",
            LastName = "Johnson",
            Title = "Finance Analyst",
            Extension = "1005",
            Phone = "617-555-0105",
            Location = "Deer Island",
            Department = "Finance",
            Email = "katherine.johnson@example.com",
            HiredDate = new DateTime(2020, 6, 15)
        }
    ];
}
