package address_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/address"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestCreate(t *testing.T) {
	testUser, err := tester.CreateUser()

	if !assert.Nil(t, err) {
		return
	}

	defer auth.DeleteUserByUserName(testUser.Username)

	context := controller.Context{Uid: testUser.Id}

	// 添加一个失败的地址
	r := address.Create(context, address.CreateAddressParams{
		Name: "123",
	})

	assert.Equal(t, schema.StatusFail, r.Status)

	// 添加一个合法的地址
	{
		var (
			Name         = "test"
			Phone        = "13888888888"
			ProvinceCode = "110000"
			CityCode     = "110100"
			AreaCode     = "110101"
			Address      = "中关村28号526"
		)

		r := address.Create(context, address.CreateAddressParams{
			Name:         Name,
			Phone:        Phone,
			ProvinceCode: ProvinceCode,
			CityCode:     CityCode,
			AreaCode:     AreaCode,
			Address:      Address,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		addressInfo := schema.Address{}

		assert.Nil(t, tester.Decode(r.Data, &addressInfo))

		defer address.DeleteAddressById(addressInfo.Id)

		assert.Equal(t, Name, addressInfo.Name)
		assert.Equal(t, Phone, addressInfo.Phone)
		assert.Equal(t, ProvinceCode, addressInfo.ProvinceCode)
		assert.Equal(t, CityCode, addressInfo.CityCode)
		assert.Equal(t, AreaCode, addressInfo.AreaCode)
		// 之前没有添加地址的话，就是默认地址
		assert.Equal(t, true, addressInfo.IsDefault)
	}
}

func TestCreateRouter(t *testing.T) {
	testUser, err := tester.CreateUser()

	if !assert.Nil(t, err) {
		return
	}

	defer auth.DeleteUserByUserName(testUser.Username)

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + testUser.Token,
	}

	body, _ := json.Marshal(&address.CreateAddressParams{
		Name:         "张三",
		Phone:        "18888888888",
		ProvinceCode: "110000",
		CityCode:     "110100",
		AreaCode:     "110101",
		Address:      "中关村28号526",
	})

	r := tester.HttpUser.Post("/v1/user/address", body, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, "", res.Message) {
		return
	}

	if !assert.Equal(t, schema.StatusSuccess, res.Status) {
		return
	}

	addressInfo := schema.Address{}

	assert.Nil(t, tester.Decode(res.Data, &addressInfo))

	defer address.DeleteAddressById(addressInfo.Id)
}
