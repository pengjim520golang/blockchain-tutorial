import Web3 from 'web3';


let web3;
const rpcUrl = "http://localhost:8545";
if( typeof window!=='undefined' && typeof window.web3 !== 'undefined' ){
    web3 = new Web3( window.web3.currentProvider )
}else{
    web3 = new Web3( new Web3.providers.HttpProvider(rpcUrl) )
}

export default web3;