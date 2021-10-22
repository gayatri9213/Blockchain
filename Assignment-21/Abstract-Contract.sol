pragma solidity >=0.5.6;

abstract contract Calculator
{
    function getResult()public virtual view returns(uint);
}
contract Test is Calculator
{
    function getResult()public override pure returns(uint)
    {
        uint a=2;
        uint b=1;
        uint result=a+b;
        return result;
    }
}