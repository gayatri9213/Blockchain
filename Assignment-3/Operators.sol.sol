pragma solidity >= 0.7.0 < 0.9.0;


contract AssignmentOperator {
    uint b = 4;
    uint a = 3;
 
    function assign() public view returns(uint) {
            uint c = 2;
 
            
            return c += c + b; 
    }
    
}

contract FinalExercise {
    uint a = 300;
    uint b = 12;
    uint f = 47;
    
    function finalize() public view returns(uint) {
        uint d = 23;

        if(a > b && b < f) { 
        d *= d;
        return d - b;
        } else {
            return d;
        }
    }
}


contract LogicalOperators {
    uint a = 17;
    uint b = 32;
    
    function logic() public view returns(uint) {
        uint result = 0;
        if(b > a && a != b) {
            result = a * b;
        }
        return result;
    }
    
   
    
}

contract ComparisonOperators {
    
    uint a = 323;
    uint b = 54;
  
    function compare() view public {
        
        require(a <= b, 'This comparison is false!' );
        
    }
    
}

contract ArithmeticOperators {
    
    function calculator() public pure returns (uint) {
        uint a = 5;
        uint b = 7;
        uint result;
        result = a--;
        return result;
    }
    
    
    
}