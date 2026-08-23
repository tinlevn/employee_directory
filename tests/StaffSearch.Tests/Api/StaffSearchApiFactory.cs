using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Mvc.Testing;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.DependencyInjection;
using StaffSearch.Data;
using StaffSearch.Tests.Helpers;

namespace StaffSearch.Tests.Api;

/// <summary>
/// Boots the real API pipeline (routing, validation, CORS, health checks)
/// with SQL Server swapped out for an isolated in-memory database.
/// </summary>
public class StaffSearchApiFactory : WebApplicationFactory<Program>
{
    private const string DatabaseName = "staffsearch-api-tests";

    protected override void ConfigureWebHost(IWebHostBuilder builder)
    {
        // Program.cs skips the SQL Server registration in the Test environment,
        // so the in-memory provider below is the only one registered.
        builder.UseEnvironment("Test");

        builder.ConfigureServices(services =>
            services.AddDbContext<StaffSearchDbContext>(options =>
                options.UseInMemoryDatabase(DatabaseName)));
    }

    public void SeedDatabase()
    {
        using var scope = Services.CreateScope();
        var context = scope.ServiceProvider.GetRequiredService<StaffSearchDbContext>();

        context.Database.EnsureDeleted();

        if (!context.Employees.Any())
        {
            context.Employees.AddRange(TestEmployees.All);
            context.SaveChanges();
        }
    }
}
