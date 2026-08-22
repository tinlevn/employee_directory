using Microsoft.AspNetCore.Mvc;
using StaffSearch.Dtos;
using StaffSearch.Extensions;
using StaffSearch.Interfaces;

namespace StaffSearch.Controllers;

[ApiController]
[Route("api/[controller]")]
[Produces("application/json")]
public class EmployeesController(IEmployeeRepository employeeRepository) : ControllerBase
{
    /// <summary>
    /// Returns employees, optionally filtered by the supplied query-string criteria.
    /// GET /api/employees?firstName=mar&amp;department=Engineering
    /// </summary>
    [HttpGet]
    [ProducesResponseType<IReadOnlyList<EmployeeDto>>(StatusCodes.Status200OK)]
    [ProducesResponseType<ProblemDetails>(StatusCodes.Status400BadRequest)]
    public async Task<ActionResult<IReadOnlyList<EmployeeDto>>> GetEmployees(
        [FromQuery] EmployeeQuery query,
        CancellationToken cancellationToken)
    {
        var employees = await employeeRepository.SearchAsync(query, cancellationToken);
        return Ok(employees.ToDtos());
    }

    /// <summary>
    /// Returns a single employee by its identifier, or 404 when it does not exist.
    /// </summary>
    [HttpGet("{id}")]
    [ProducesResponseType<EmployeeDto>(StatusCodes.Status200OK)]
    [ProducesResponseType<ProblemDetails>(StatusCodes.Status404NotFound)]
    public async Task<ActionResult<EmployeeDto>> GetEmployeeById(string id, CancellationToken cancellationToken)
    {
        var employee = await employeeRepository.GetByIdAsync(id, cancellationToken);
        return employee is null ? NotFound() : Ok(employee.ToDto());
    }
}
