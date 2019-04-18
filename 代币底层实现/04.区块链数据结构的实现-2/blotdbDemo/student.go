package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Student struct {
	Id     int
	Name   string
	School string
}

//创建一个学生
func NewStudent(id int, name, school string) *Student {
	return &Student{id, name, school}
}

func (stu *Student) getInfo() {
	fmt.Println(stu.Id, stu.Name, stu.School)
}

//对结构体进行序列化
func (stu *Student) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(stu)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func UnSerializeStudent(byteStudentData []byte) *Student {
	var student Student
	decoder := gob.NewDecoder(bytes.NewReader(byteStudentData))
	err := decoder.Decode(&student)
	if err != nil {
		log.Panic(err)
	}
	return &student

}
