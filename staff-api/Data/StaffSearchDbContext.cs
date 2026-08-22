using Microsoft.EntityFrameworkCore;
using StaffSearch.Entities;

namespace StaffSearch.Data;

public class StaffSearchDbContext(DbContextOptions<StaffSearchDbContext> options) : DbContext(options)
{
    public DbSet<Employee> Employees => Set<Employee>();
}
