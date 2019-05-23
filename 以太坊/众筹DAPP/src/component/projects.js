import React,{Component} from 'react';
import web3 from '../libs/web3';
import ProjectList from '../libs/projectList';
import Project from '../libs/project';
class Projects extends Component {

  constructor(...args){
      super(...args)

      //初始化项目列表状态
      this.state = {
        projects:[]
      }
  }


  async componentDidMount(){
    const addresses = await ProjectList.methods.getProjects().call();
    const summaryList = await Promise.all( addresses.map( address=> Project(address).methods.getSummary().call()  ) ) 
    //console.log( Object.values(summaryList[0]) )
    const projects = addresses.map( (address,index)=>{
        let [description,minInvest,maxInvest,goal,balance,investorsCount,paymentCount,owner] = Object.values(summaryList[index])

        return {address,description,minInvest,maxInvest,goal,balance,investorsCount,paymentCount,owner}
    } )

    this.setState({projects})
  
  }



  renderProject(project,index){
    let progress = project.balance / project.goal * 100;
    return (
        
     
          <div className="col-md-6" key={index}>
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
                            <td colSpan="5" className="text-center">
                                <a className="btn  btn-sm btn-primary" href={`/details/${project.address}`}>立即投资</a>&nbsp;&nbsp;
                                <a className="btn  btn-sm btn-info" href={`/details/${project.address}`}>查看详情</a>
                            </td>
                          </tr>
                      </tfoot>
                  </table>
              </div>
          </div>

      
    

);
  }



  render() {
    let _projects = this.state.projects.map( (project,index)=>this.renderProject(project,index) )
    return (
      <div className="container" style={{marginTop:"50px"}}>
        <div className="page-header">
          <h2>项目列表</h2>
        </div>
        <div className="row" >
          {_projects}
        </div> 
      </div>
    )
  }
}

export default Projects;
