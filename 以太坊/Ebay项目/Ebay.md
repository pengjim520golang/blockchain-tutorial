# 商店合约编写

```java
pragma solidity ^0.4.17;
//商店合约
contract EcommerceStore{
    /**
    产品拍卖状态
    Open : 开始拍卖
    Sold : 已完成拍卖
    Unsold : 未完成拍卖
     */
    enum ProductStatus { Open, Sold, Unsold }
    /**
    产品状态
    New : 全新产品
    Used : 非全新产品(相当与非一手货品)
     */
    enum ProductCondition{ New,Used }
    //产品描述
    struct Product{
        //产品id
        uint id;
        //产品名称
        string name;
        //产品类型
        string category;
        //产品在ipfs中的hash
        string imageLink;
        //产品的图片描述
        string descLink;
        //开始竞标时间
        uint auctionStartTime;
        //竞标结束时间
        uint auctionEndTime;
        //拍卖价格
        uint startPrice;
        //赢家的地址
        address highestBidder;
        //赢家的竞标的价格
        uint highestBid;
        //第2价高者的地址
        uint secondHighestBid;
        //参与竞标的总人数
        uint totalBids;
        //产品的拍卖状态，对应结构体的Open,Sold,Unsold
        ProductStatus status;
        //产品的属性状态，对应结构体的New,Used
        ProductCondition condition;
    }
    //产品的索引自增器 
    uint public productIndex;
    //通过地址，找到该产品发送者的产品
    mapping (address => mapping(uint => Product)) stores;
    //产品索引==>发送者地址
    mapping (uint => address) productIdInStore;

    //构造函数,初始化
    constructor() public {
        productIndex = 0; 
    }

    //添加产品
    function addProductToStore( string _name, 
                                string _category, 
                                string _imageLink, 
                                string _descLink, 
                                uint _auctionStartTime,
                                uint _auctionEndTime, 
                                uint _startPrice, 
                                uint _productCondition) public {
        //如果开始时间小于结束时间则无法添加产品???
        require (_auctionStartTime < _auctionEndTime);
        //产品id自动增长
        productIndex += 1;

        Product memory product = Product(productIndex, 
                                        _name, 
                                        _category, 
                                        _imageLink, 
                                        _descLink, 
                                        _auctionStartTime, 
                                        _auctionEndTime,
                                        _startPrice, 
                                        0, 
                                        0, 
                                        0, 
                                        0, 
                                        ProductStatus.Open, 
                                        ProductCondition(_productCondition));
        //产品发送者发送的产品列表
        stores[msg.sender][productIndex] = product;
        //通过产品索引找到发送者是谁
        productIdInStore[productIndex] = msg.sender;
    }
}
```