import "interface.sol";

contract ewasm is E {
	function get_vm() public override returns (string memory) {
		return "ewasm";
	}

	function get_number() public override returns (int) {
		return 54321;
	}
	
	function call_get_vm(E e) public returns (string memory) {
		// solc can't do this
		return "ewasm called " + e.get_vm();
	}

	function call_get_number(E e) public returns (int) {
		return e.get_number();
	}
}
