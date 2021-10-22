pragma solidity >=0.5.6;

contract A
{
    uint innerVal=100;
    
    function innerAddTen(uint val)public pure returns(uint)
    {
        return val+10;
    }
}

contract B is A
{
    function outerAddTen(uint val)public pure returns(uint)
    {
        return A.innerAddTen(val);
    }
    
      function getInnerVal()public view returns(uint)
    {
        return A.innerVal;
    }
}