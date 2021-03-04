pragma solidity ^0.5;

contract EventEmitter {
    // indexed puts it in topic
    event ManyTypes(
        bytes32 indexed direction,
        bool trueism,
        string german ,
        int64 indexed newDepth,
        int bignum,
        string indexed hash);

    event ManyTypes2(
        bytes32 indexed direction,
        bool trueism,
        string german ,
        int128 indexed newDepth,
        int8 bignum,
        string indexed hash);

    function emitOne() public {
        emit ManyTypes("Downsie!", true, "Donaudampfschifffahrtselektrizitätenhauptbetriebswerkbauunterbeamtengesellschaft", 102, 42, "hash");
    }

    function emitTwo() public {
        emit ManyTypes2("Downsie!", true, "Donaudampfschifffahrtselektrizitätenhauptbetriebswerkbauunterbeamtengesellschaft", 102, 42, "hash");
    }
}
