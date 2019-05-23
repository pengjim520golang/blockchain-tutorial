import React,{Component} from 'react';
import web3 from '../libs/web3';

class Accounts extends Component {

  constructor(...args){
      super(...args)
      this.state = {
          addresses:[]
      }
  }


  async componentDidMount(){
      const accounts = await web3.eth.getAccounts();
      const balances = await Promise.all( accounts.map(address=>web3.eth.getBalance(address))  ) 
      this.setState({
          addresses:accounts.map( (account,index)=>{
              return {address:account,balance:balances[index]}
          } )
      })
  }



  render() {
   
    let addresses = this.state.addresses.map((account,index)=>{
        return <li key={index}>{account.address}-->{ web3.utils.fromWei(account.balance,'ether')}ETH</li>
    })


    return (
        <ul>
            {addresses}
        </ul>
    );
  }
}

export default Accounts;
