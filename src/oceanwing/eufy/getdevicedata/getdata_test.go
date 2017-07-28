package getdevicedata

import (
	"fmt"
	"testing"
)

func Test_mydata(t *testing.T) {
	//getDataFromDB()
	// u := time.Now().Unix()
	// fmt.Printf("unit time: %d\n", u)
	u := registerUser()
	fmt.Printf("uid: %s\n", u.uid)
	fmt.Printf("token: %s\n", u.token)
	fmt.Printf("email: %s\n", u.email)
	b := bindDevice("T10121707000C994", u)
	fmt.Printf("bind device result: %t\n", b)
}
