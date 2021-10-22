pragma solidity >= 0.7.0 < 0.9.0;


contract FallBack {
    
    event Log(uint gas);
    
    fallback () external payable {
      
        emit Log(gasleft());
    }
    
    function getBalance() public view returns(uint) {
        // return the stored balance of the contract 
        return address(this).balance;
    }
    
}


// new contract will send ether to Fallback contract which will triggger fallback functions 

contract SendToFallBack {
    
    function transferToFallBack(address payable _to) public payable {
        // send ether with the transfer method
        
        _to.transfer(msg.value);
    }
    
    
    function callFallBack(address payable _to) public payable {
        // send ether with the call method 
        (bool sent,) = _to.call{value:msg.value}('');
        require(sent, 'Failed to send!');
    }
    
}

