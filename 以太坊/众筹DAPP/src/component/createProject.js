import React,{Component} from 'react';
import ProjectList from '../libs/projectList';
import web3 from '../libs/web3';

class createProject extends Component {

    
    constructor(...args) {
        super(...args);
        
        this.state = {
            addresses:[],
            description: '',
            minInvest: 0,
            maxInvest: 0,
            goal: 0,
            errmsg: '',
            loading: false
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

    async createICO(){

        //console.log(this.refs["description"].value) 
      
                                   
        let description = this.refs["description"].value
        let minInvest = this.refs["minInvest"].value
        let maxInvest = this.refs["maxInvest"].value
        let goal = this.refs["goal"].value
        
        //console.log(this.state)
        //let {description, minInvest, maxInvest, goal} = this.state;

        if(!description){
            return this.setState({errmsg: "项目名称不能为空!"});
        }
        if(minInvest <= 0){
            return this.setState({ errmsg: '项目最小投资金额必须大于0' });
        }
        if(maxInvest <= 0){
            return this.setState({ errmsg: '项目最大投资金额必须大于0' });
        }
        if(parseInt(maxInvest) < parseInt(minInvest)){
            return this.setState({ errmsg: '项目最小投资金额必须小于最大投资金额' });
        }
        if(goal <= 0){
            return this.setState({ errmsg: '项目募资上限必须大于0' });
        }

        let minInvestInWei = web3.utils.toWei(minInvest, "ether");
        let maxInvestInWei = web3.utils.toWei(maxInvest, "ether");
        let goalInWei = web3.utils.toWei(goal, "ether");
        
        
        let accounts = await web3.eth.getAccounts();
        let owner = accounts[0];
        let result = await ProjectList.methods.createProject(description, minInvestInWei, maxInvestInWei, goalInWei)
                    .send({from: owner, gas: '5000000'});
        console.log(result);
        this.setState({loading: true});
        this.setState({errmsg: ""});
        
        
    }


  render() {
    let addresses = this.state.addresses.map((account,index)=>{
        return <span className="label label-info" key={index}>您的账号余额: {web3.utils.fromWei(account.balance,'ether')} ETH</span>
    })
    return (
        <div className="container" style={{marginTop:"50px"}}>
            <div className="page-header">
                <h2>发起项目</h2>
                {addresses}
            </div>
            {this.state.errmsg!=='' ? <div className="alert alert-danger" role="alert">{this.state.errmsg}</div> : ""}
            <div className="col-md-6">
                <div className="form-group">
                    <label htmlFor="projectDesc">项目描述:</label>
                    <input type="text" className="form-control" id="projectDesc" ref="description" placeholder="例如:项目名称" />
                    
                    <label htmlFor="ProjectminInvest">最小投资金额:</label>
                    <div className="input-group">
                        <input type="text" id="ProjectminInvest" ref="minInvest" className="form-control" placeholder="单位:以太（ETH）" aria-describedby="basic-addon2" />
                        <span className="input-group-addon" id="basic-addon2">以太(ETH)</span>
                    </div>


                    <label htmlFor="ProjectmaxInvest">最大投资金额:</label>
                    <div className="input-group">
                        <input type="text" id="ProjectmaxInvest" ref="maxInvest" className="form-control" placeholder="单位:以太（ETH）" aria-describedby="basic-addon2" />
                        <span className="input-group-addon" id="basic-addon2">以太(ETH)</span>
                    </div>



                    <label htmlFor="ProjectmaxInvest">募资上限:</label>
                    <div className="input-group">
                        <input type="text" id="Projectgoal" ref="goal" className="form-control" placeholder="单位:以太（ETH）" aria-describedby="basic-addon2" />
                        <span className="input-group-addon" id="basic-addon2">以太(ETH)</span>
                    </div>                    
                
                </div>
    
                
                {this.state.loading ? <div class="alert alert-success" role="alert">项目创建成功<a href="/" className="btn btn-sm btn-success"><span className="glyphicon glyphicon-pushpin"></span>&nbsp;&nbsp;返回项目列表</a></div> : <button type="button" className="btn btn-lg btn-primary" onClick={this.createICO.bind(this)}>创建项目</button>}
                
            </div>


        </div>
    );
  }
}

export default createProject;