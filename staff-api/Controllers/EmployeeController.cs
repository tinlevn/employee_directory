using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Mvc;
using Staff_Search.Entity;
using Staff_Search.Interface;

namespace Staff_Search.Controllers
{
    [Route("[controller]")]
    [ApiController]
    public class EmployeeController : ControllerBase
    {

        private readonly IEmployeeRepository _employeeRepository;
        public EmployeeController(IEmployeeRepository repo)
        {
            _employeeRepository = repo;
        }

        [HttpGet]
        public async Task<ActionResult<IEnumerable<Employee>>> GetEmployees()
        {
            var allEmployees = await _employeeRepository.GetEmployeesAsync();
            return Ok(allEmployees);
        }

        [HttpPost("filter")]
        public async Task<ActionResult<IEnumerable<Employee>>> GetEmployeesWithFilter([FromBody] Filter specs)
        {
            var result = await _employeeRepository.GetEmployeeResult(specs);
            return Ok(result);
        }
    }

}
