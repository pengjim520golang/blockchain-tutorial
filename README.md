# 区块链技术学习前言导读
非常开心能来到火链区块链学院为大家分享相关的我对区块链技术学习心得的文章。正当我准备这些教材的时候,我感觉非常有必要跟大家谈谈时下中国我们学习区块链的一些资料,目前在国内区块链技术如火如荼,各种教程也遍布互联网,这非常符合国情。然而非常负责任地告诉大家,如果你是以程序开发人员对区块链进行学习和理解的话。恭喜你！这才是真正理解区块链技术最正确的学习方向,因为从技术层面出发去理解你会遇不到太多花哨的概念,并能从真正意义上明白区块链技术的原理和应用的方向,对市面上一些说法做到真正的去伪存真。

## 关于比特币的教材
比特币是由中本聪使用c++编写的一个区块链技术的数字货币系统,俗称比特币系统。很多技术人员学习区块链都是先从比特币系统学习进行入门的。为了便于理解,不同的程序语言方向的技术人员使用不同的语言根据中本聪比特币系统的源码机制从原理上实现了不同语言版本的比特币系统。所以你看到有Java,Python,Golang等各种版本的比特币系统。

<font color=red>注意：你需要明白的一点是只有中本聪的c++版本的比特币系统才是真正的比特币系统,其余任何语言的版本都是在模拟实现比特币系统,目的是使用自己擅长的语言去理解比特币的机制和原理</font>

我们需要学习比特币系统的机制和原理选择使用的是Golang技术,模拟实现比特币机制的过程又被称为"公链开发"。最为权威的参考教程都来自一个外国开发者Ivan Kuznetsov所编写的7篇文章。

[Ivan Kuznetsov的Blog](https://jeiwan.cc/)中的7篇文章如下所示:

* Building Blockchain in Go. Part 1: Basic Prototype [区块链的基础原型]
* Building Blockchain in Go. Part 2: Proof-of-Work [工作量证明]
* Building Blockchain in Go. Part 3: Persistence and CLI [持久化存储和客户端命令行]]
* Building Blockchain in Go. Part 4: Transactions1 [交易1-实现UTXO交易的机制]
* Building Blockchain in Go. Part 5: Addresses [比特币地址和数字签名]
* Building Blockchain in Go. Part 6: Transactions2 [交易2-优化UTXO交易]
* Building Blockchain in Go. Part 7: NetWork  [模拟比特币网络]