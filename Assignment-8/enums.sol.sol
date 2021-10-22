pragma solidity >= 0.7.0 < 0.9.0;


contract enumsLearn {
    enum frenchFriesSize {LARGE, MEDIUM, SMALL}
   	frenchFriesSize choice;
     	frenchFriesSize constant defaultChoice = frenchFriesSize.MEDIUM;
    
    function setSmall() public {
        choice = frenchFriesSize.SMALL;
    }
    
    function getChoice() public view returns(frenchFriesSize) {
            return choice;
    }
    
    function getDefaultChoice() public view returns (uint) {
           return uint(defaultChoice);
    }
    
    enum shirtColor {RED, WHITE, BLUE}
    shirtColor choice;
    shirtColor constant defaultChoice = shirtColor.BLUE;
    
    function setWhite() public {
        choice = shirtColor.WHITE;
    }
    
    function getChoice() public view returns(shirtColor) {
        return choice;
    }
    
        function getDefaultChoice() public view returns(uint) {
        return uint(defaultChoice);
    }
}

