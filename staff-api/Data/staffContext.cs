using Microsoft.EntityFrameworkCore;
using Staff_Search.Entity;

namespace Staff_Search.Data
{
    public class staffContext : DbContext
    {
        protected readonly IConfiguration Configuration;
        public staffContext(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        protected override void OnConfiguring(DbContextOptionsBuilder options)
        {
            // connect to sql server with connection string from app settings
            options.UseSqlServer(Configuration.GetConnectionString("DefaultConnection"));

        }

        //DBSet for corresponding tables in database
        public DbSet<Employee> Employees { get; set; }

    }
}
