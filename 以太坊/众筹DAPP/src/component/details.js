import React,{Component} from 'react';
import web3 from '../libs/web3';
import Project from '../libs/project';
class Details extends Component {

  constructor(...args){
      super(...args)
      console.log(web3)
      this.state = {
          project:{
            address:"",
            description:"", 
            minInvest:"0", 
            maxInvest:"0", 
            goal:"0", 
            balance:"0", 
            investorsCount:"0", 
            paymentCount:"0", 
            owner:""   
          },
          errmsg:""
      }
  }


  async componentDidMount(){
    //console.log(this.props.match.params.address)
    const address = this.props.match.params.address
    const contract = Project(address);
    let summary = await contract.methods.getSummary().call();
    console.log(summary)
    let [description, minInvest, maxInvest, 
        goal, balance, investorsCount, paymentCount, owner] = Object.values(summary);
    this.setState({
        project:{address,description, minInvest, maxInvest, 
            goal, balance, investorsCount, paymentCount, owner}
    })
  }

  async contributeProject() {
    let amount = this.refs["amount"].value;
    let {minInvest,maxInvest} = this.state.project;
    let minInvestInETH = web3.utils.fromWei(minInvest, 'ether');
    let maxInvestInETH = web3.utils.fromWei(maxInvest, 'ether');
    if(amount <= 0){
        return this.setState({ errmsg: '项目最小投资金额必须大于0' });
    }
    if(parseInt(amount) < parseInt(minInvestInETH)){
        return this.setState({errmsg: '投资金额必须大于最小投资金额'});
    }
    if(parseInt(amount) > parseInt(maxInvestInETH)){
        return this.setState({errmsg: '投资金额必须小于最大投资金额'});
    }   
    

    let accounts = await web3.eth.getAccounts();
    let sender = accounts[0];

    const contract = Project(this.props.match.params.address);
    await contract.methods.contribute()
                    .send({from: sender, value: web3.utils.toWei(amount, 'ether'), gas: '5000000'});
    this.setState({errmsg: ""});
    this.setState({loading: true});
    setTimeout(()=>{
        window.location.reload();
    }, 1000);
  }


  renderProject(project){
    let progress = project.balance / project.goal * 100;

    //console.log(typeof project.goal)
    return (
        
     
          <div className="col-md-6">
              <div className="panel panel-default">
                  <div className="panel-heading">{project.description}</div>


                  <div className="panel-body">

                    <div className="progress">
                      <div  style={{width: `"${progress}%"`}} className="progress-bar progress-bar-striped active" role="progressbar" aria-valuenow={`"${progress}"`} aria-valuemin="0" aria-valuemax="100" >
                        <span className="sr-only">45% Complete</span>
                      </div>
                    </div>

                  </div>


                  <table className="table table-bordered">
                      <thead>
                          <tr>
                            <th>募资上限</th>
                            <th>最小投资金额</th>
                            <th>最大投资金额</th>
                            <th>参投人数</th>
                            <th>已募集资金数量</th>
                          </tr>
                      </thead>
                      <tbody>
                        <tr>
                            <td>{ web3.utils.fromWei(project.goal,"ether") } Eth</td>
                            <td>{ web3.utils.fromWei(project.minInvest,"ether")} Eth</td>
                            <td>{ web3.utils.fromWei(project.maxInvest,"ether")} Eth</td>
                            <td>{project.investorsCount} 人</td>
                            <td>{web3.utils.fromWei(project.balance,"ether")} Eth</td>
                        </tr>                      
                      </tbody>
                      <tfoot>
                        <tr>
                            <td colSpan="5">
                                    <div className="input-group">
                                    <input type="text" ref="amount" className="form-control" placeholder="单位:以太(ETH)" />
                                    <span className="input-group-btn">
                                    <button className="btn btn-primary" type="button" onClick={this.contributeProject.bind(this)}>立即投资</button>
                                    </span>
                                </div>                        
                            </td>
                        </tr>
                      </tfoot>
                  </table>
              </div>
          </div>
 
      
    

);
  }



  render() {
    console.log(this.state.project)
    return (
      <div className="container" style={{marginTop:"50px"}}>
        <div className="page-header">
          <h2>项目详情</h2>
        </div>
        {this.state.errmsg!=='' ? <div className="alert alert-danger" role="alert">{this.state.errmsg}</div> : ""}
        
        {this.state.loading ? <div className="alert alert-success" role="alert">投资成功</div> : ""}

        <div className="row" >
          {this.renderProject(this.state.project)}
        </div> 
      </div>
    )
  }
}

export default Details;
