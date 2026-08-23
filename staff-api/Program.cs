using Microsoft.EntityFrameworkCore;
using StaffSearch.Data;
using StaffSearch.Data.Repositories;
using StaffSearch.Interfaces;

var builder = WebApplication.CreateBuilder(args);

const string CorsPolicy = "StaffSearchClient";

builder.Services.AddControllers();
builder.Services.AddProblemDetails();

// Integration tests run the host in the "Test" environment and register their
// own in-memory provider; production/development always use SQL Server.
if (!builder.Environment.IsEnvironment("Test"))
{
    builder.Services.AddDbContext<StaffSearchDbContext>(options =>
        options.UseSqlServer(builder.Configuration.GetConnectionString("DefaultConnection")));
}

builder.Services.AddScoped<IEmployeeRepository, EmployeeRepository>();

builder.Services.AddHealthChecks()
    .AddDbContextCheck<StaffSearchDbContext>("database");

builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var allowedOrigins = builder.Configuration
    .GetSection("Cors:AllowedOrigins")
    .Get<string[]>() ?? [];

builder.Services.AddCors(options =>
    options.AddPolicy(CorsPolicy, policy =>
        policy.WithOrigins(allowedOrigins)
              .AllowAnyHeader()
              .WithMethods("GET")));

var app = builder.Build();

if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}
else
{
    app.UseExceptionHandler();
    app.UseHsts();
    app.UseHttpsRedirection();
}

app.UseStatusCodePages();

app.UseCors(CorsPolicy);

app.UseAuthorization();

app.MapControllers();
app.MapHealthChecks("/health");

app.Run();

// Exposes the implicit Program class to the integration-test project.
public partial class Program { }
