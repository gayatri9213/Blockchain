pragma solidity >= 0.7.0 < 0.9.0;



contract learnFunctions {
 
   
   function addValues(uint a) public pure returns (uint) {
       
       uint b = 3;
       uint result = a + b;
       return result;
   }
   
   
   function addNewValues() public pure returns (uint){
       uint a = 1;
       uint b = 5;
       uint result = a + b;
       return result;
   }
   
   
   
   
   uint b = 4; 
   
    function multiplyCalculatorByFour(uint a) public view returns (uint) {
        uint result = a * b;
        return result;
    }
    
        function divideCalculatorByFour(uint a) public view returns (uint) {
        uint result = a / b;
        return result;
    }
   
}