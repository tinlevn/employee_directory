using System.ComponentModel.DataAnnotations;

namespace Staff_Search.Entity
{
    public class Employee
    {
        [Key]
        public required string Guid { get; set; }
        public string? FirstName { get; set; }
        public string? LastName { get; set; }    
        public string? Title { get; set; }
        public string? Extension { get; set; }
        public string? Phone { get; set; }
        public string? Location { get; set; }
        public string? Department { get; set; }
        public string? Email { get; set; }
        public DateTime? HiredDate { get; set; }
    }
}
