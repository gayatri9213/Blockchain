pragma solidity >= 0.7.0 < 0.9.0;

 
contract learnArrays {
    
  uint[] public myArray; /// add 1 2 3 4 
  uint[] public myArray2;
  uint[200] public myFixedSizedArray;
  
  
    
    function push(uint number) public {
        myArray.push(number);
    }
    
  
    function pop() public {
        myArray.pop();
    }
    
   
    function getlength() public view returns (uint) {
        return myArray.length;
    }
    
    function remove(uint i) public {
        delete myArray[i];
       
    }
    
    
    uint[] public changeArray;
    
    function removeElement(uint i) public {
        changeArray[i] = changeArray[changeArray.length - 1];
        changeArray.pop();
    } 
    
    function test() public {
        changeArray.push(1);
        changeArray.push(2);
        changeArray.push(3);
        changeArray.push(4);
    }
    
    function getArray() public view returns (uint) {
        return changeArray.length;
}
}