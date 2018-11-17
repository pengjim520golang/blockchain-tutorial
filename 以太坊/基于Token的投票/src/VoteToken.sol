pragma solidity ^0.4.18;
contract VoteToken {
	struct voter {  
		address voterAddress;  
		uint tokensBought;  
		uint[] tokensUsedPerCandidate; 
	}
    /*
    查询投票人信息的 mapping：
    给定一个投票人的账户地址，我们想要显示他的信息。
    我们会使用 voterInfo 字段来存储信息。
    */
	mapping (address => voter) public voterInfo;
	//候选人获得的票数
	mapping (bytes32 => uint) public votesReceived; 
	bytes32[] public candidateList;
    //发行 token 总量
	uint public totalTokens;
    //所有剩余的 token
	uint public balanceTokens; 
    //每个 token 的价格
	uint public tokenPrice;
	constructor(uint tokens, uint pricePerToken, bytes32[] candidateNames) public { 
	    //候选人
		candidateList = candidateNames;
		//发行 token 总量
		totalTokens = tokens;
		//所有剩余的 token
		balanceTokens = tokens; 
		//每个 token 的价格
		tokenPrice = pricePerToken; 
	}
	function buy() payable public returns (uint) {  
	    //购买的token数量
		uint tokensToBuy = msg.value / tokenPrice;
		require(tokensToBuy <= balanceTokens);
		voterInfo[msg.sender].voterAddress = msg.sender;
		voterInfo[msg.sender].tokensBought += tokensToBuy;
		balanceTokens -= tokensToBuy;  
		return tokensToBuy; 
	}
	function totalVotesFor(bytes32 candidate) view public returns (uint) {  
	    //候选人的获得token数量就是票数
		return votesReceived[candidate];
	}
	function voteForCandidate(bytes32 candidate, uint votesInTokens) public {
	    //获取候选人的索引
		uint index = indexOfCandidate(candidate); 
		require(index != uint(-1));
		if ( voterInfo[msg.sender].tokensUsedPerCandidate.length == 0) {
		    //初始化候选人在投票人中的token数据
			for(uint i = 0; i < candidateList.length ;i++) {    
				voterInfo[msg.sender].tokensUsedPerCandidate.push(0);
			}  
		}
		//当前投票人的可用的token=购买的token总数-投票使用的token总数
		uint availableTokens = voterInfo[msg.sender].tokensBought-totalTokensUsed(voterInfo[msg.sender].tokensUsedPerCandidate); 
		require (availableTokens >= votesInTokens);
		//候选人的获得token数量就是票数
        votesReceived[candidate] += votesInTokens; 
        //投票人在候选人中花费的token
		voterInfo[msg.sender].tokensUsedPerCandidate[index] += votesInTokens; 
	}
	//投票所花费掉的token总数
    function totalTokensUsed(uint[] _tokensUsedPerCandidate) private pure returns (uint) {  
		uint totalUsedTokens = 0;  
		for(uint i = 0; i < _tokensUsedPerCandidate.length; i++) {   
			totalUsedTokens += _tokensUsedPerCandidate[i];
		}
		return totalUsedTokens; 
	}
	//获取候选人的索引
    function indexOfCandidate(bytes32 candidate) view public returns (uint) {  
		for(uint i = 0; i < candidateList.length; i++) {   
			if (candidateList[i] == candidate) {    
				return i;   
			}  
		}  
		return uint(-1);
	}
	//已经卖出了多少token
    function tokensSold() view public returns (uint) {  
		return totalTokens - balanceTokens; 
	}
	//返回当前投票人所拥有的token和投票对应花费的token
	function voterDetails(address user) view public returns (uint, uint[]) {  
		return (voterInfo[user].tokensBought, voterInfo[user].tokensUsedPerCandidate); 
	}
	/*
	当一个用户调用 buy 方法发送以太来购买 token 时，所有的以太去了哪里？
	所有以太都在合约里。
	每个合约都有它自己的地址，这个地址里面存储了这些钱。可这些钱怎么拿出来呢？
	我们已经在这里定义了 transferTo 函数，它可以让你转移所有钱到指定的账户
	。该方法目前所定义的方式，任何人都可以调用，
	并向他们的账户转移以太，这并不是一个好的选择。
	*/
	function transferTo(address account) public {  
		account.transfer(address(this).balance); 
	}
	//返回所有候选人
	function allCandidates() view public returns (bytes32[]) {  
		return candidateList; 
	}
}
