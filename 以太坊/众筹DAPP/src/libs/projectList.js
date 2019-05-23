import web3 from './web3';
//获取ABI
import projectList from '../build/ProjectList.json';
//获取合约的地址
import address from '../build/address.json';

let contract = new web3.eth.Contract(JSON.parse(projectList.interface) , address );

export default contract;