# Ubuntu安装相关开发环境

毫无疑问的一点,我们无论开发什么肯定在Linux下进行的(条件允许也可以使用Mac OS做开发),这里选用Ubuntu进行,所以我们需要先安装一些必备的开发工具

## 必备包

安装前,建议更新一下apt-get并且安装一些必备工具,这样在一定程度上可以避免一些很奇葩的问题发生,当然也不是100%的,当然如果你觉得你的ubuntu系统没有必要进行这一步,那你也可以跳过去.

```
sudo apt-get update
sudo apt-get install build-essential libssl-dev
```

## 安装git

git的安装是必须的,因为有可能我需要利用git进行clone一些仓库和管理我们开发的代码版本

```
sudo apt-get install git
```

## 安装nvm

nvm是Nodejs的版本管理器,所以安装Nodejs之前必须首先安装nvm

**使用git 安装**

```
cd ~/ - -切到主目录
git clone https://github.com/creationix/nvm.git .nvm 
cd ~/.nvm  --进入nvm代码目录
git checkout v0.33.11 --切换到v0.33.11版本
. nvm.sh --执行命令
```
> 在终端里面,执行 nvm --version,出来版本号,说明安装成功网上很多其他安装方法,这里推荐的git安装,你使用其他方法当然也可以

## 安装nodejs

这里安装的版本是Nodejs 8.11.3

```
nvm install v8.11.3
```

如果安装成功我们可以使用`nvm ls`查看到,安装nodejs成功会同时安装了npm


## 安装cnpm

由于npm下载很多包非常的慢,所以我们可以切换到国内到淘宝镜像,直接安装cnpm即可

```
npm install -g cnpm --registry=https://registry.npm.taobao.org
```

## 安装truffle

```
cnpm install -g truffle
```

## 安装truffle unbox react

```
mkdir /usr/local/src/react-demo
cd /usr/local/src/react-demo
truffle unbox react
cd client
npm install bignumber.js --save-dev
```

## 安装truffle unbox webpack

```
mkdir /usr/local/src/webpack-demo
cd /usr/local/src/webpack-demo
truffle unbox webpack
```


