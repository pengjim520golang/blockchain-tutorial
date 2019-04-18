package main

import (
	_ "fmt"
	"log"

	"github.com/boltdb/bolt"
)

//定义为bolt数据库的文件路径
const boltdbFile = "./database/blotTest.db"

//定义表的名称,接受类型是byte
const users = "userTable"

func main() {
	//打开已经存在的数据库或者创建一个数据库
	db, err := bolt.Open(boltdbFile, 0666, nil)
	if err != nil {
		log.Panic(err)
	}

	//这种叫bolt重建索引操作
	err = db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(users))
		//如果一张表不存在我们删除就会出错，对于当前逻辑来说我们应该先忽略表不存在的错误
		if err != nil && err != bolt.ErrBucketNotFound {
			log.Panic(err)
		}
		return nil
	})

	if err != nil {
		//该err是bolt.Update的异常
		log.Panic(err)
	}

	//利用db实例中方法：Update , View,创建表使用Update的回调方法
	//回调参数的第1个参数是acid类型,事务增删查改的操作
	err = db.Update(func(tx *bolt.Tx) error {

		bucket, err := tx.CreateBucket([]byte(users))
		if err != nil {
			return err
		}
		////构建学生的value值
		zs := NewStudent(100, "zhangsan", "清华大学")
		//序列化张三
		zsValue := zs.Serialize()
		ls := NewStudent(200, "lisi", "北京大学")
		//序列化lisi
		lsValue := ls.Serialize()

		//把学生的数据录入到数据库当中
		bucket.Put([]byte("zs"), zsValue)
		bucket.Put([]byte("ls"), lsValue)

		return nil
	})

	if err != nil {
		//该err是bolt.Update的异常
		log.Panic(err)
	}
	//读取数据

	db.View(func(tx *bolt.Tx) error {
		//定义一个keys的集合才能遍历
		keys := make(map[string][]byte)

		keys["zs"] = []byte("zs")
		keys["ls"] = []byte("ls")

		//获取表的名称
		bucket := tx.Bucket([]byte(users))

		for _, key := range keys {
			data := bucket.Get([]byte(key))
			stu := UnSerializeStudent(data)
			stu.getInfo()
		}

		return nil
	})
}
