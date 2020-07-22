##DEMO

	//初始化短信
	initialization.InitSms()
	var smsParam = initialization.RegParam{
		Realname: "熊建波",
		Username: "熊建波",
		Password: "123456",
	}
	//测试短信
	res,err:=initialization.Sms.
		SetSignName("云立方网").
		SetPhone("17111111111").
		SetTemplateCode("xxxx").
		Send(smsParam)
