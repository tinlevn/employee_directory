using StaffSearch.Extensions;
using StaffSearch.Tests.Helpers;

namespace StaffSearch.Tests.Extensions;

public class EmployeeMappingTests
{
    [Fact]
    public void ToDto_MapsEveryField_AndPromotesGuidToId()
    {
        var employee = TestEmployees.All[0];

        var dto = employee.ToDto();

        Assert.Equal(employee.Guid, dto.Id);
        Assert.Equal(employee.FirstName, dto.FirstName);
        Assert.Equal(employee.LastName, dto.LastName);
        Assert.Equal(employee.Title, dto.Title);
        Assert.Equal(employee.Extension, dto.Extension);
        Assert.Equal(employee.Phone, dto.Phone);
        Assert.Equal(employee.Location, dto.Location);
        Assert.Equal(employee.Department, dto.Department);
        Assert.Equal(employee.Email, dto.Email);
        Assert.Equal(employee.HiredDate, dto.HiredDate);
    }

    [Fact]
    public void ToDtos_MapsCollections()
    {
        var dtos = TestEmployees.All.ToDtos();

        Assert.Equal(TestEmployees.All.Length, dtos.Count);
        Assert.Equal(TestEmployees.All.Select(e => e.Guid), dtos.Select(d => d.Id));
    }
}
