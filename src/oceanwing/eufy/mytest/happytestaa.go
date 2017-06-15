package mytest

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"oceanwing/eufy/protobuf.lib/pbexample"
)

// ================================  Marshal ==========================
func pbMarshal() []byte {

	m := make(map[string]string)
	m["key1"] = "value1"
	m["key2"] = "value2"

	p := &pbexample.Person{
		Name:   proto.String("Matt"),
		Id:     proto.Int32(611),
		Email:  proto.String("matt.ma@oceanwing.com"),
		Phone:  makePhone(),
		MyMap:  m,
		Choice: chooseOthers(),
	}

	data, err := proto.Marshal(p)
	if err != nil {
		fmt.Printf("Mashal pb message error: %s\n", err)
		return nil
	}
	return data
}

func makePhone() []*pbexample.Person_PhoneNumber {
	phone1 := &pbexample.Person_PhoneNumber{
		Number: proto.String("13812345678"),
		Type:   pbexample.Person_MOBILE.Enum(),
	}

	// phone2 的 Type 使用默认值：HOME
	phone2 := &pbexample.Person_PhoneNumber{
		Number: proto.String("0755-28888889"),
	}

	return []*pbexample.Person_PhoneNumber{phone1, phone2}
}

func chooseEat() *pbexample.Person_Eat {
	return &pbexample.Person_Eat{
		Eat: &pbexample.Food{
			Name:  proto.String("pig meat"),
			Price: proto.Int32(88),
		},
	}
}

func chooseLive() *pbexample.Person_Live {
	return &pbexample.Person_Live{
		Live: &pbexample.House{
			Addr:  proto.String("ShenZhen,Nanshan"),
			Owner: proto.String("yuanfang"),
		},
	}
}

func chooseOthers() *pbexample.Person_Others {
	return &pbexample.Person_Others{
		Others: "Free",
	}
}

// ===================== UnMarshal ==========================

func pbUnMarshal(payload []byte) {
	p := &pbexample.Person{}
	err := proto.Unmarshal(payload, p)
	if err != nil {
		fmt.Printf("UnMarshal pb message fail: %s\n", err)
		return
	}

	name := p.GetName()
	id := p.GetId()
	email := p.GetEmail()
	fmt.Printf("Base info, name: %s, id: %d, email: %s\n", name, id, email)

	phones := p.GetPhone()
	for _, phone := range phones {
		number := phone.GetNumber()
		ty := phone.GetType()
		fmt.Printf("Phone Number: %s, type: %s\n", number, ty)
	}

	m := p.GetMyMap()
	for k, v := range m {
		fmt.Printf("key: %s, value: %s\n", k, v)
	}

	eat := p.GetEat()
	if eat != nil {
		eatName := eat.GetName()
		eatPrice := eat.GetPrice()
		fmt.Printf("eat name: %s, eat price: %d\n", eatName, eatPrice)
	}

	live := p.GetLive()
	if live != nil {
		addr := live.GetAddr()
		owner := live.GetOwner()
		fmt.Printf("house address: %s, owner: %s\n", addr, owner)
	}

	others := p.GetOthers()
	if others != "" {
		fmt.Printf("Others: %s\n", others)
	}
}
