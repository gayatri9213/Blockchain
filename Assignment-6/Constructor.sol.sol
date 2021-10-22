pragma solidity >= 0.7.0 < 0.9.0;



contract Member {
    string name;
    uint age;
    
    constructor(string memory _name, uint _age) public {
        name = _name;
        age = _age;
        return name;
    }
}

contract Teacher is Member // Inherited form member contract

{
    constructor(string memory n, uint a) Member(n, a) public{}
    
    function getName() public view returns(string memory) {
        return name;
    }
}

contract Base {
   uint data;
   constructor(uint _data) public {
      data = _data;   
   }
}

contract Derived is Base (5) {
   function getData() public view returns(uint) {
       return data;
   }   
}