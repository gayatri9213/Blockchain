pragma solidity >= 0.7.0 < 0.9.0;


contract learnStrings {   
    
    
    string favoriteColor = 'blue';
    
    function getColor() public view returns (string memory) {

        return favoriteColor;
    }
    
    
    function changeColor(string memory _color) public {

        favoriteColor = _color;
    }
    
    function getColorLength() public view returns(uint) {

        bytes memory stringToBytes = bytes(favoriteColor);

        return stringToBytes.length;

    }
    
}