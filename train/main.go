package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type users struct {
	Root  xml.Name `xml:"root"`
	Users []user   `xml:"row"`
}

type user struct {
	ID        int    `xml:"id"`
	FirstName string `xml:"first_name"`
	LastName  string `xml:"last_name"`
	Age       int    `xml:"age"`
	About     string `xml:"about"`
	Gender    string `xml:"gender"`
}

func findUsersTest(u *users) {
	file, err := os.Open("./dataset.xml")
	if err != nil {
		fmt.Println("error")
	}
	defer file.Close()
	dataBytes, _ := ioutil.ReadAll(file)
	if err := xml.Unmarshal(dataBytes, &u); err != nil {
		fmt.Printf("error: %v", err)
	}
}

func main() {
	var u users
	findUsersTest(&u)
	for _, v := range u.Users {
		fmt.Printf("%v \n", v.LastName)
	}
}
