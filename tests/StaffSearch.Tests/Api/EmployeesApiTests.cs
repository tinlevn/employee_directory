using System.Net;
using System.Net.Http.Json;
using StaffSearch.Dtos;

namespace StaffSearch.Tests.Api;

public class EmployeesApiTests : IClassFixture<StaffSearchApiFactory>
{
    private readonly HttpClient _client;

    public EmployeesApiTests(StaffSearchApiFactory factory)
    {
        factory.SeedDatabase();
        _client = factory.CreateClient();
    }

    [Fact]
    public async Task GetEmployees_ReturnsAllEmployees()
    {
        var response = await _client.GetAsync("/api/employees");

        response.EnsureSuccessStatusCode();
        var employees = await response.Content.ReadFromJsonAsync<List<EmployeeDto>>();

        Assert.NotNull(employees);
        Assert.Equal(4, employees.Count);
        Assert.Contains(employees, e => e.Id == "id-ada");
    }

    [Fact]
    public async Task GetEmployees_AppliesQueryStringFilters()
    {
        var employees = await _client.GetFromJsonAsync<List<EmployeeDto>>(
            "/api/employees?firstName=ra&department=Engineering");

        Assert.NotNull(employees);
        var employee = Assert.Single(employees);
        Assert.Equal("Grace", employee.FirstName);
    }

    [Fact]
    public async Task GetEmployees_IgnoresBlankQueryParameters()
    {
        var employees = await _client.GetFromJsonAsync<List<EmployeeDto>>(
            "/api/employees?firstName=&department=");

        Assert.NotNull(employees);
        Assert.Equal(4, employees.Count);
    }

    [Fact]
    public async Task GetEmployees_Returns400_WhenQueryParameterExceedsMaxLength()
    {
        var tooLong = new string('x', 101);

        var response = await _client.GetAsync($"/api/employees?firstName={tooLong}");

        Assert.Equal(HttpStatusCode.BadRequest, response.StatusCode);
    }

    [Fact]
    public async Task GetEmployeeById_ReturnsEmployee_WhenItExists()
    {
        var response = await _client.GetAsync("/api/employees/id-katherine");

        response.EnsureSuccessStatusCode();
        var employee = await response.Content.ReadFromJsonAsync<EmployeeDto>();

        Assert.NotNull(employee);
        Assert.Equal("Katherine", employee.FirstName);
        Assert.Equal("Finance", employee.Department);
    }

    [Fact]
    public async Task GetEmployeeById_Returns404_WhenMissing()
    {
        var response = await _client.GetAsync("/api/employees/does-not-exist");

        Assert.Equal(HttpStatusCode.NotFound, response.StatusCode);
    }

    [Fact]
    public async Task HealthCheck_ReturnsHealthy()
    {
        var response = await _client.GetAsync("/health");

        response.EnsureSuccessStatusCode();
        Assert.Equal("Healthy", await response.Content.ReadAsStringAsync());
    }

    [Fact]
    public async Task GetEmployees_SetsCorsHeader_ForConfiguredOrigin()
    {
        using var request = new HttpRequestMessage(HttpMethod.Get, "/api/employees");
        request.Headers.Add("Origin", "http://localhost:4200");

        var response = await _client.SendAsync(request);

        Assert.True(response.Headers.TryGetValues("Access-Control-Allow-Origin", out var origins));
        Assert.Contains("http://localhost:4200", origins);
    }
}
