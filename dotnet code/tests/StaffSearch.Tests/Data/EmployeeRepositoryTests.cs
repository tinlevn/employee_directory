using StaffSearch.Data.Repositories;
using StaffSearch.Dtos;
using StaffSearch.Tests.Helpers;

namespace StaffSearch.Tests.Data;

// NOTE: The in-memory provider evaluates string.Contains with .NET (case-sensitive)
// semantics, while SQL Server uses a case-insensitive collation in production.
// These tests therefore assert filtering structure with matching-case values.
public class EmployeeRepositoryTests
{
    [Fact]
    public async Task GetAllAsync_ReturnsAllEmployees_OrderedByLastName()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.GetAllAsync();

        Assert.Equal(4, result.Count);
        Assert.Equal(["Hopper", "Johnson", "Lovelace", "Turing"], result.Select(e => e.LastName).ToArray());
    }

    [Fact]
    public async Task GetByIdAsync_ReturnsEmployee_WhenItExists()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.GetByIdAsync("id-grace");

        Assert.NotNull(result);
        Assert.Equal("Grace", result.FirstName);
    }

    [Fact]
    public async Task GetByIdAsync_ReturnsNull_WhenMissing()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.GetByIdAsync("does-not-exist");

        Assert.Null(result);
    }

    [Fact]
    public async Task SearchAsync_WithNoCriteria_ReturnsAll()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery());

        Assert.Equal(4, result.Count);
    }

    [Fact]
    public async Task SearchAsync_FiltersFirstNameByContains()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery { FirstName = "da" });

        var employee = Assert.Single(result);
        Assert.Equal("Ada", employee.FirstName);
    }

    [Fact]
    public async Task SearchAsync_FiltersLastNameByContains()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery { LastName = "ove" });

        var employee = Assert.Single(result);
        Assert.Equal("Lovelace", employee.LastName);
    }

    [Fact]
    public async Task SearchAsync_FiltersJobTitleAgainstTitleColumn()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery { JobTitle = "Engineer" });

        Assert.Equal(2, result.Count);
        Assert.All(result, e => Assert.Contains("Engineer", e.Title!));
    }

    [Fact]
    public async Task SearchAsync_FiltersDepartmentByExactMatch()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery { Department = "Engineering" });

        Assert.Equal(2, result.Count);
        Assert.All(result, e => Assert.Equal("Engineering", e.Department));
    }

    [Fact]
    public async Task SearchAsync_DepartmentDoesNotMatchSubstring()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        // "Engineer" is a substring of "Engineering" but must NOT match: dropdown
        // criteria use exact-match semantics.
        var result = await repository.SearchAsync(new EmployeeQuery { Department = "Engineer" });

        Assert.Empty(result);
    }

    [Fact]
    public async Task SearchAsync_FiltersLocationByExactMatch()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery { Location = "Deer Island" });

        Assert.Equal(2, result.Count);
        Assert.All(result, e => Assert.Equal("Deer Island", e.Location));
    }

    [Fact]
    public async Task SearchAsync_CombinesCriteriaWithAndSemantics()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery
        {
            FirstName = "ra",
            Department = "Engineering"
        });

        var employee = Assert.Single(result);
        Assert.Equal("Grace", employee.FirstName);
    }

    [Theory]
    [InlineData(null)]
    [InlineData("")]
    [InlineData("   ")]
    public async Task SearchAsync_IgnoresBlankCriteria(string? blank)
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery
        {
            FirstName = blank,
            LastName = blank,
            JobTitle = blank,
            Department = blank,
            Location = blank
        });

        Assert.Equal(4, result.Count);
    }

    [Fact]
    public async Task SearchAsync_TrimsCriteria()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery { Department = "  Finance  " });

        var employee = Assert.Single(result);
        Assert.Equal("Finance", employee.Department);
    }

    [Fact]
    public async Task SearchAsync_ReturnsEmployeesSortedByLastThenFirstName()
    {
        await using var context = TestDbContextFactory.CreateSeeded();
        var repository = new EmployeeRepository(context);

        var result = await repository.SearchAsync(new EmployeeQuery { Location = "Deer Island" });

        Assert.Equal(["Johnson", "Turing"], result.Select(e => e.LastName).ToArray());
    }
}
