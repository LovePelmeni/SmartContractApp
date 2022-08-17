pragma solidity ^0.4.11; 

contract NFTToken {
    string public name = "token";
    string public symbol = "TRX";
    uint256 public decimal = "";

    mapping(address => uint256) public balanceOf;
    mapping(address => (address => uint256)) public allowance;
    uint256 public stopped = false; 
}