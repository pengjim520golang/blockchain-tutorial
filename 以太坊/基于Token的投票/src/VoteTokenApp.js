import React, { Component } from "react";
import VoteTokenContract from "./contracts/VoteToken.json";
import getWeb3 from "./utils/getWeb3";
import truffleContract from "truffle-contract";


class VoteTokenApp extends Component {
  state = { 
        web3: null, 
        accounts: null, 
        contract: null,
        tContract:null,
        candidates:[
          {
            id:1,
            name:"Alice",
            count:-1
          },
          {
            id:2,
            name:"Bob",
            count:-1
          },
          {
            id:3,
            name:"Tommy",
            count:-1
          },    
        ],
        tokens:0,
        tokensBalance:0,
        tokensSold:0,
        tokenPrice:0,
        voterInfo:{
          token:0,
          votings:[]
        },
        contractBalance:0,
        selfTokenBalance:0
  }
  
  componentDidMount = async () => {
    try {
      const web3 = await getWeb3();
      const accounts = await web3.eth.getAccounts();
      const Contract = truffleContract(VoteTokenContract);
      Contract.setProvider(web3.currentProvider);
      const instance = await Contract.deployed();
      this.setState({ web3, accounts, contract: instance},
        this.showBasicInfo);
    } catch (error) {
      alert(
        `Failed to load web3, accounts, or contract. Check console for details.`
      );
      console.log(error);
    }
  };
  //输出基本信息
  showBasicInfo = async ()=>{
    const {contract,web3} = this.state;
    /**
    await contract.buy({from:accounts[0],value:web3.utils.toWei("1")});
    let balance_tokens = await contract.balanceTokens();
    console.log(balance_tokens.toNumber())

    let voter_details = await contract.voterDetails("0x33E6A400a18530bf745012069a0fd555883Cf5e4");
    console.log("token:",voter_details[0].words[0])
    **/

    //获取候选人相关票数信息
    let candidateList = this.state.candidates;
    for(let index in  candidateList){
      let name = web3.utils.toHex(candidateList[index].name);
      let votingCount = await contract.totalVotesFor(name);
      candidateList[index].count = votingCount.toNumber();
    }
    
    //获取tokens发型总数
    let total_tokens = await contract.totalTokens();
    //获取tokens剩余数量
    let balance_tokens = await contract.balanceTokens();
    //获取售出的token数量
    let sold_tokens = await contract.tokensSold();
    //获取token的单价
    let price_token = await contract.tokenPrice();
    //获取合约中的余额
    let contract_balance = await web3.eth.getBalance(contract.address);
    //console.log(contract_balance.toString())
    //console.log(price_token.toString())
    this.setState({
      candidates:candidateList,
      tokens:total_tokens.toNumber(),
      tokensBalance:balance_tokens.toNumber(),
      tokensSold:sold_tokens.toNumber(),
      tokenPrice:web3.utils.fromWei(price_token.toString(),'ether'),
      contractBalance:web3.utils.fromWei(contract_balance.toString(),'ether')
    });

  };

  //购买token
  buyToken = async ()=>{
    const {contract,web3,accounts} = this.state;
    //获取输入的token数量
    let tokenToBuy = this.refs.tokenNum.value;
    console.log(tokenToBuy);
    //获取token的单价
    let tokenPriceEth = this.state.tokenPrice;
    console.log(tokenPriceEth);
    //需要花费的以太
    let spentEth = tokenToBuy * tokenPriceEth;
    //console.log(typeof spentEth.toString());
    let msgValue = web3.utils.toWei(spentEth.toString(),'ether');
    await contract.buy({from:accounts[0],value:msgValue});
    //获取tokens剩余数量
    let balance_tokens = await contract.balanceTokens();
    //获取售出的token数量
    let sold_tokens = await contract.tokensSold();
    //合约中余额
    let contract_balance = await web3.eth.getBalance(contract.address);   
    this.setState({
      tokensBalance:balance_tokens.toNumber(),
      tokensSold:sold_tokens.toNumber(),
      contractBalance:web3.utils.fromWei(contract_balance.toString(),'ether')
    },this.queryVoterDetails);
  };

  //实现token投票
  votingFor = async ()=>{
    const {contract,accounts,web3} = this.state;
    //获取候选人的姓名
    let candidate = this.refs.candidateName.value;
    //对姓名进行bytes32转换
    let candidateNameToHex = web3.utils.toHex(candidate);
    //获取投票数
    let voteNum = this.refs.tokenSpentValue.value;
    //给候选人投票
    await contract.voteForCandidate(candidateNameToHex,voteNum,{from:accounts[0]});
    //获取候选人相关票数信息
    let candidateList = this.state.candidates;
    //刷新投票
    for(let index in  candidateList){
      //console.log(candidateList[index].name,candidate)
       if(candidateList[index].name === candidate)
       {
          //console.log("for",candidateList[index].name,candidate)
          let votingCount = await contract.totalVotesFor(candidateNameToHex);
          candidateList[index].count = votingCount.toNumber();
          break;
       }
       
    }
    this.setState({
      candidates:candidateList
    },this.queryVoterDetails);      
  };
  //查询投票人的详细信息
  queryVoterDetails = async ()=>{
    const {contract,accounts,web3} = this.state;
    //获取投票人的地址
    let voterAdr = this.refs.voterAddress.value;
    //console.log(typeof voterAdr)
    if(voterAdr!==""){
      let voterInfo = this.state.voterInfo;
      let voter_details = await contract.voterDetails(voterAdr);
      
      voterInfo.token = voter_details[0].words[0];
      voterInfo.votings = voter_details[1];
      //计算当前投票者所消费的余额
      let selfTokenBalance = 0;
      let spentToken = 0;
      for(let rs of voter_details[1]){
        //console.log(rs)
        //console.log(rs.words[0])
        spentToken += rs.words[0]
      }
      console.log( voterInfo.token)
      console.log(spentToken)
      selfTokenBalance = voterInfo.token-spentToken;
      //console.log(selfTokenBalance)
      this.setState({voterInfo:voterInfo,selfTokenBalance:selfTokenBalance});
    }
    
  };
  render() {
    if (!this.state.web3) {
      return <div>正在加载web3,contract....如果无法显示Dapp请刷新浏览器.</div>;
    }
    return (

    
      <div className="row">
          <div className="col-md-12 text-center">
               <h2>基于Token的投票</h2>
                <hr />
          </div>
          <div className="col-md-5">
          <table className="table table-bordered table-striped">
              <thead>
                <tr className="info">
                  <th className="text-center">候选人ID</th>
                  <th className="text-center">候选人姓名</th>
                  <th className="text-center">得票数</th>
                </tr>
              </thead>

              <tbody>
                {
                  this.state.candidates.map(candidate=>{
                    return (
                      <tr key={candidate.id}>
                        <td className="text-center">
                          {candidate.id}
                        </td>
                        <td className="text-center">
                          {candidate.name}
                        </td>
                        <td className="text-center">
                          {candidate.count}
                        </td>             
                      </tr>
                    )
                  })
                }
              </tbody>


          </table>
          </div>
          <div className="col-md-5 col-md-offset-2">

          <table className="table table-bordered table-striped">

              <tbody>
                <tr>
                  <td className="text-right info" style={{width:'50%'}}>
                      Token发行数量:
                  </td>
                
                  <td className="text-left" style={{width:'50%'}}>
                     {this.state.tokens}
                  </td>
                </tr>

                <tr>
                  <td className="text-right info">
                     剩余Token数量:
                  </td>
                  <td className="text-left">
                  {this.state.tokensBalance}
                  </td>                 
                </tr>

                <tr>
                  <td className="text-right info">
                     售出Token数量:
                  </td>
                  <td className="text-left">
                  {this.state.tokensSold}
                  </td>                 
                </tr>

 
                 <tr>
                  <td className="text-right info">
                     Token单价成本:
                  </td>
                  <td className="text-left">
                  {this.state.tokenPrice}<strong>Eth</strong>
                </td>               

                </tr>


                 <tr>
                  <td className="text-right info">
                     合约余额:
                  </td>
                  <td className="text-left">
                    {this.state.contractBalance}<strong>Eth</strong>
                  </td>               

                </tr>



              </tbody>

          </table>

          </div>

          <div className="col-md-5">
            <table className="table table-bordered table-striped">

                <thead>
                    <tr className="info">
                      <th colSpan="2">投票候选人</th>
                    </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>花费token票数:</td>
                    <td><input type="text" className="form-control" ref="tokenSpentValue" id="tokenSpentValue" /></td>
                  </tr>
                  <tr>
                    <td>候选人姓名:</td>
                    <td><input type="text" className="form-control" ref="candidateName" id="candidateName" /></td>
                  </tr>
                </tbody>
                <tfoot>
                    <tr>
                      <td colSpan="2" className="text-center">
                        <button type="button" className="btn btn-primary" onClick={this.votingFor.bind(this)}>确认投票</button>
                      </td>
                    </tr>                 
                </tfoot>
            </table>
          </div>



          <div className="col-md-5 col-md-offset-2">
            <table className="table table-bordered table-striped">

                <thead>
                    <tr className="info">
                      <th colSpan="2">购买Token</th>
                    </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>购买数量:</td>
                    <td><input type="text" className="form-control" ref="tokenNum" id="tokenNum" /></td>
                  </tr>
                </tbody>
                <tfoot>
                    <tr>
                      <td colSpan="2" className="text-center">
                        <button type="button" className="btn btn-primary" onClick={this.buyToken.bind(this)}>确认购买</button>
                      </td>
                    </tr>                 
                </tfoot>
            </table>
          </div>

          <div className="col-md-12 text-center">
          <h2>投票人查询</h2>
          <hr />
          </div>
          
          <div className="col-md-5">
            <table className="table table-bordered table-striped">

                <thead>
                    <tr className="info">
                      <th colSpan="2">投票人信息查询</th>
                    </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>投票人:</td>
                    <td><input type="text" className="form-control" ref="voterAddress" id="voterAddress" /></td>
                  </tr>
                </tbody>
                <tfoot>
                    <tr>
                      <td colSpan="2" className="text-center">
                        <button type="button" className="btn btn-primary" onClick={this.queryVoterDetails.bind(this)}>查询</button>
                      </td>
                    </tr>                 
                </tfoot>
            </table>
          </div>
           

          <div className="col-md-5 col-md-offset-2">
            <table className="table table-bordered table-striped">

                <thead>
                    <tr className="info">
                      <th colSpan="2">投票人信息明细</th>
                    </tr>
                </thead>
                <tbody>
                  <tr>
                    <td className="text-right" style={{width:"50%"}}>拥有的token数量:</td>
                    <td>
                        {this.state.voterInfo.token}
                    </td>
                  </tr>

                  <tr>
                    <td className="text-right" style={{width:"50%"}}>Token余额:</td>
                    <td>
                        {this.state.selfTokenBalance}
                    </td>
                  </tr>
                  <tr>
                      <td colSpan="2">投票明细</td>
                  </tr>

                  {
                    this.state.voterInfo.votings.map((details,index)=>{
                      
                      return (

                          
                          <tr key={"id-"+index}>
  
                          <td className="text-right" style={{width:"50%"}}>
                              {this.state.candidates[index].name}
                          </td>
                          <td>
                              {details.words[0]}
                          </td>
                        </tr>
                      )
                    })
                  }


                </tbody>

            </table>
          </div>


      </div> 
    );
  }
}

export default VoteTokenApp;
