```javascript
pragma solidity ^0.4.22;
contract VotingByToken{
    /**
     投票人
     - 投票人必须具有股权(token)[tokenBought]
     - 投票人可以为候选人投票,需要消耗手上token(选举人索引=>token)[ tokenUsedCandidate ]
     - 投票人拥有自己的以太坊账号[account]
     - 投票人的信息(voterInfo)
    */
    struct voter {
        address account;
        uint tokenBought;
        uint[] tokenUsedCandidate; 
    }
    mapping(address=>voter) voterInfo;
    
    /**
     *token
    - token的发行总量(totalTokens)
    - token剩余量(balanceTokens)
    - token的购买单价(tokenPrice/wei)
    */
    uint public totalTokens;
    uint public balanceTokens;
    uint public tokenPrice;
    /**
    候选人
    - 候选人投票信息( 候选人的名字=>投票数 )
    - 候选人的列表
    */
    bytes32[] public candidatesList;
    mapping(bytes32=>uint) votesReceived;
    //初始化合约部署者的管理员是谁
    address admin;
    constructor(bytes32[] candidates,uint tokensTotalSupply,uint tokenPrePrice) public {
        candidatesList = candidates;
        totalTokens = tokensTotalSupply;
        balanceTokens = tokensTotalSupply;
        tokenPrice = tokenPrePrice;
        admin = msg.sender;
    }
    /**
    购买token
    - 购买量(tokenNum)=以太的数量(wei)/token单价
    - 验证发行的token发行余额是否大于等于购买量
    - 更新投票人的相关信息
    - 更新token发行余额
    ps: payable 扣除真正的以太给合约账号
    */
    function buy() payable public returns( uint ){
        uint tokenNum = msg.value / tokenPrice;
        require( tokenNum<=balanceTokens,"must be less than token of balances" );
        voterInfo[msg.sender].account = msg.sender;
        voterInfo[msg.sender].tokenBought += tokenNum;
        balanceTokens -= tokenNum;
        return tokenNum;
    }
    //把合约中的以太进行提取
    function transferToAccount(address _to) public {
        require(admin == msg.sender,"transferToAccount control must be admin");
        _to.transfer( address(this).balance );
    }
    //ps:在新版本中必须加上以下代码才能使得transfer函数交易真正的以太币
    function() payable{}
    
    /**
    为候选人投票
    - 使用token给候选人投票
       a.初始化候选人的票数
       b.求出一共消耗的token有多少
       c.可用的token是否大于等于投票的token
       d.为候选人加上票数,同时更新自己为谁投过票
    - 获取指定候选人的票数
    - 编写一个获取候选人索引的方法,-1 not found
    */
    function VotingFor(bytes32 candidateName,uint voteInTokens) public{
         uint index = indexOfCandidate(candidateName);
         require( index !=uint(-1),"must be in candidates list" ) ;
         if( voterInfo[msg.sender].tokenUsedCandidate.length == 0 ){
             for(uint i=0;i<candidatesList.length;i++){
                voterInfo[msg.sender].tokenUsedCandidate.push(0);
             }            
         }
         uint ableTokenCount = voterInfo[msg.sender].tokenBought - totalTokenUsed(voterInfo[msg.sender].tokenUsedCandidate);
         require(ableTokenCount>=voteInTokens,"ableTokenCount must gt voteInTokens");
         votesReceived[candidateName] += voteInTokens;
         voterInfo[msg.sender].tokenUsedCandidate[index] += voteInTokens;
    }
    
    function totalTokenUsed(uint[] tokenUsedCandidates) pure public returns(uint){
        uint totalUsed = 0;
        for(uint i=0;i<tokenUsedCandidates.length;i++){
            totalUsed += tokenUsedCandidates[i];
        }
        return totalUsed;
    }
    
    function getVoteCount(bytes32 candidateName) view public returns (uint){
        require( indexOfCandidate(candidateName)!=uint(-1),"must be in candidates list" ) ;
        return votesReceived[candidateName];
        
    }
    function indexOfCandidate(bytes32 candidateName) public returns (uint){
        for(uint i=0;i<candidatesList.length;i++){
            if(candidateName==candidatesList[i]){
                return i;
            }
        }
        return uint(-1);
    }
    
    function tokensSold() view public returns (uint) {  
		return totalTokens - balanceTokens;
	}
    
    function voterDetails(address user) view public returns (uint, uint[]) {  
		return (voterInfo[user].tokenBought,voterInfo[user].tokenUsedCandidate);
	}
}

```