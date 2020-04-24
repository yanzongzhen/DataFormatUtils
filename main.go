package main

import (
	"github.com/yanzongzhen/DataFormatUtils/xml"
	"github.com/yanzongzhen/Logger/logger"
	"log"
	"sync"
	"time"
)

var once sync.Once

func oneFunc(index int) {
	log.Println("one Func start ===", index)
	once.Do(func() {
		log.Println("once start", index)
		time.Sleep(time.Second * 3)
		log.Println("once end", index)
	})
	log.Println("one Func end ===", index)
}

func main() {

	logger.InitLogConfig(logger.DEBUG, true)

	/*	data := `{
		   "paginator": {
			   "page": 1,
			   "rows": 10,
			   "pages": 1,
			   "totalRows": 2,
			   "hasNextPage": false,
			   "hasPrevPage": false,
			   "hasOtherPage": false,
			   "startIndex": 0,
			   "endIndex": 10,
			   "nextPage": 1,
			   "prevPage": 1,
			   "currentPages": [
				   1
			   ]
		   },
		   "loanList":
		}`

			err := json.TravelJsonData([]byte(data), func(path string, value interface{}) bool {
				logger.Info(path, value)
				return false
			})*/

	//res, _ := xml.ChangeXmlValueNew("<Request><YYID>$.yyid</YYID></Request>", func(path string) string {
	//	logger.Debug(path)
	//	return "1111"
	//})

	//var text template.HTML = `<?xml version='1.0' encoding='UTF-8'?><S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"><S:Body><getyyxx xmlns="http://tempuri.org/" xmlns:ns2="http://schemas.microsoft.com/2003/10/Serialization/"><param_xml>{{index . 0}}</param_xml></getyyxx></S:Body></S:Envelope>`

	//t, _ := template.New("asda").Parse(`<?xml version='1.0' encoding='UTF-8'?><S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/"><S:Body><getyyxx xmlns="http://tempuri.org/" xmlns:ns2="http://schemas.microsoft.com/2003/10/Serialization/"><param_xml>{{index . 0}}</param_xml></getyyxx></S:Body></S:Envelope>`)
	//buffer := bytes.Buffer{}
	//array := make([]interface{}, 0, 10)
	//text := "<Request><YYID>1231</YYID></Request>"
	////strings.TrimSpace()
	//
	//array = append(array, template.HTMLEscapeString(text))
	//_ = t.Execute(&buffer, array)
	//logger.Debug(buffer.String())
	//logger.Debug(template.HTMLEscapeString(buffer.String()))

	//res, _ := xml.ChangeXmlValueNew("<p><s sfzhm=\"$.yyid\"/></p>", func(path string) string {
	//	logger.Debug(path)
	//	return "1111"
	//})
	//
	//logger.Debug(res)
	//	newData, err := json.ChangeValueTravel([]byte(`{
	//    "header": {
	//        "authCode": "b8389a0f37bd7439b49318348b9445b2",
	//        "senderID": "3413-0050"
	//    },
	//    "request": {
	//        "data": [
	//            {
	//                "name": "condition",
	//                "type": "GsbString",
	//                "value": ""
	//            },
	//            {
	//                "name": "requiredItems",
	//                "type": "GsbString",
	//                "value": "<requiredItems><item><DQYE>当前余额</DQYE><DWBL>单位比例</DWBL><DWMC>单位名称</DWMC><DWYHZH>单位开户行</DWYHZH><DWYJE>单位月缴额</DWYJE><DWZH>单位账号</DWZH><GMSFHM>公民身份号码(入参)</GMSFHM><GRBL>个人比例</GRBL><GRYJE>个人月缴额</GRYJE><JCJS>缴存基数</JCJS><KHRQ>开户日期</KHRQ><XM>姓名</XM><ZHHJY>最后汇缴月</ZHHJY><ZHZT>账户状态</ZHZT></item></requiredItems>"
	//            },
	//            {
	//                "name": "clientInfo",
	//                "type": "GsbString",
	//                "value": "<clientInfo><loginName>sadmin</loginName></clientInfo>"
	//            }
	//        ],
	//        "method": "requestQueryV1",
	//        "serviceCode": "3413-0050-1-00000006",
	//        "version": 1
	//    }
	//}`), func(path string, oldValue interface{}) (i interface{}, b bool) {
	//		if path == "$.request.data[0].value" {
	//			return "<condition><item><GMSFHM>342221195103010048</GMSFHM></item></condition>", true
	//		}
	//		return nil, false
	//	})
	//
	//	ioutil.WriteFile("./res.txt")
	//	logger.Debug(string(newData))

	//data := "<?xml version=\"1.0\" encoding=\"GBK\"?>\n" +
	//	"<p>\n" +
	//	"    <s typelistforcbxx_ds=\"qsrq:s,zzrq:s,cbzt:s,cbztmc:s,xzbz:s,xzbzmc:s,jfrylb:s,jfrylbmc:s,zzny:s,dwbh:s,cbdwmc:s,cbrq:s\">1</s>\n" +
	//	"    <s _lesb__errcode_=\"0\">1</s>\n" +
	//	"    <s ryid=\"06372321198307270022\">1</s>\n" +
	//	"    <s errflag=\"0\">1</s>\n" +
	//	"</p>"
	//err := xml.TraverseXmlNew([]byte(data), func(name string, path string, value string) bool {
	//	logger.Debug(name)
	//	logger.Debug(path)
	//	logger.Debug(value)
	//	return false
	//})

	//res, err := mixformat.XmlToJsonNew([]byte(data), "{\n"+
	//	"    \"result\":[\n"+
	//	"        {\n"+
	//	"            \"a\":\"$.p.d.r(cbrq).string\"\n"+
	//	"        }\n"+
	//	"    ]\n"+
	//	"}")
	//if err != nil {
	//logger.Error(err)
	//} else {
	//	logger.Debug(string(res))
	//}

	//	busData := `{
	//    "code": "1",
	//    "msg": "",
	//    "zgzzylgrjcxxList": [
	//        [
	//            {
	//                "jfjs": "4410.0",
	//                "cbzt": "正常参保",
	//                "snmljgrjfbx": "11427.29",
	//                "DWMC": "达伟（重庆）物流有限公司",
	//                "sysjjfje": "352.8",
	//                "GRBH": "2045723122",
	//                "CBXZ": "职工养老保险",
	//                "ljys": "54",
	//                "CBD": "沙坪坝区社保局"
	//            }
	//        ],
	//        [
	//            {
	//                "cbzt": "暂停参保",
	//                "snmljgrjfbx": "679.53",
	//                "DWMC": "重庆市众业人力资源管理顾问有限公司",
	//                "GRBH": "2056642460",
	//                "CBXZ": "职工养老保险",
	//                "ljys": "4",
	//                "CBD": "江北区社会保险局"
	//            }
	//        ],
	//        [
	//            {
	//                "cbzt": "暂停参保",
	//                "snmljgrjfbx": "2026.72",
	//                "DWMC": "达伟（重庆）物流有限公司",
	//                "GRBH": "2059501562",
	//                "CBXZ": "职工养老保险",
	//                "ljys": "10",
	//                "CBD": "沙坪坝区社保局"
	//            }
	//        ]
	//    ],
	//    "zglzylgrjcxxList": [],
	//    "jmzzylgrjcxxList": [],
	//    "jmlzylgrjcxxList": []
	//}`
	//
	//	template := `{
	//    "code": "0000",
	//    "data": {
	//        "base_info": {
	//            "work":
	//            [
	//                {
	//                    "personalNumber": "$.zgzzylgrjcxxList[0].[0].GRBH.string",
	//                    "insuredState": "$.zgzzylgrjcxxList[0].[0].cbzt.string"
	//                }
	//            ],
	//            "offwork":  [
	//                {
	//                    "personalNumber": "$.zglzylgrjcxxList[0].[0].GRBH.string",
	//                    "insuredState": "$.zglzylgrjcxxList[0].[0].cbzt.string"
	//                }
	//            ]
	//        }
	//    },
	//    "message": "成功"
	//}`
	//
	//	res := json.TransJsonFormatNew([]byte(busData), template)
	//
	//	fmt.Println(string(res))

	//testList := make([]string, 1, 10)
	//testList[0] = "aaaa"
	//testList = append(testList, "aaa")
	//log.Println(testList)
	//t := time.Unix(1530608641325/1000, 0)
	//log.Println(t.String())
	//var score interface{}
	//score = 12.34
	//log.Println(fmt.Sprintf("%v", score))
	//go oneFunc(1)
	//oneFunc(2)

	//matcher, err := regexp.Compile(`个人\S?\s*([0-9,]+\s*(.\s*[0-9]+)?%)`)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//src := "个人:8% 单位:8%"
	//if matcher.MatchString(src) {
	//	res := matcher.FindAllStringSubmatch(src, -1)
	//	log.Println(res[0][1])
	//}

	//params := make(map[string]string)
	//params["act"] = "getSearchCall"
	////
	//params["sq"] = `{'keyword':"test",'searchField':1,'kindId':3}`
	//
	//_, _, res, err := utils.Get(`http://np.fj12345.gov.cn/callcenter/api/fzCallServlet`, nil, params)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Println(string(res))

	//	log.SetFlags(log.Lshortfile | log.LstdFlags)
	//source := `{"status":"200","message":"数据获取成功","result":{"data":[{"workerId":"20703173","name":"王井龙","year":"2018","companyPayMonths":"0","personalPayMonths":"3","salaryBase":"17400","amount":"1392","dealList":[{"workerId":"20703173","companyId":"0245515","amount":"464","name":"王井龙","dealDate":"201803"},{"workerId":"20703173","companyId":"0245515","amount":"464","name":"王井龙","dealDate":"201802"},{"workerId":"20703173","companyId":"0245515","amount":"464","name":"王井龙","dealDate":"201801"}]}],"personnalInfo":{"name":"王井龙","workerId":"20703173","companyId":"0245515","company":"辽宁前锦众程人力资源有限公司","birthday":"1986-11-04","status":"正常"}}}`
	//
	//res, err := json.ChangeValueTravel([]byte(source), func(path string, value interface{}) (interface{}, bool) {
	//	return false, true
	//})
	//if err != nil {
	//	log.Println(err)
	//} else {
	//	log.Println(string(res))
	//}

	//	template := `
	//{
	//    "status": "200",
	//    "message": "数据获取成功",
	//    "result": {
	//        "data": [
	//            {
	//				"year": "2018",
	//                "dealList": [
	//                    {
	//                        "amount": "464",
	//                        "dealDate": "201803"
	//                    }
	//                ]
	//            }
	//        ]
	//    }
	//}`
	//	transFormJson := `
	//	{
	//    "code": "0000",
	//    "message": "$.message.string",
	//    "data": {
	//        "trans_detail": [
	//            {
	//                "deal_record": [
	//                    {
	//                        "year": "$.result.data[0].year.int.e%x",
	//                        "money_amount": "$.result.data[0].dealList[0].amount.float",
	//                        "deal_date": "$.result.data[0].dealList[0].dealDate.string"
	//                    }
	//                ]
	//            }
	//        ]
	//    }
	//}`
	//template := make(map[string]interface{})
	//json2.Unmarshal([]byte(transFormJson), &template)
	//res := json.TransJsonFormat([]byte(source), template)
	//if res != nil {
	//	data, _ := json2.Marshal(res)
	//	log.Println(string(data))
	//}
	//
	//	insertValues, err := json.GetJsonValueByTemplate([]byte(source), []byte(template))
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	data, err := j.Marshal(insertValues)
	//	log.Println(string(data))
	//
	//
	//	transFormMap := make(map[string]interface{})
	//	err = j.Unmarshal([]byte(transFormJson), &transFormMap)
	//	//log.Println(err)
	//
	//	finalRes := json.TransJsonFormat(insertValues, transFormMap)
	//	dataFinal, err := j.Marshal(finalRes)
	//	log.Println("final=============")
	//	log.Println(string(dataFinal))
	//
	//
	//	xmlToJsonRes, err := mixformat.XmlToJson([]byte(`<?xml version="1.0" encoding="GBK"?>
	//	<p>
	//	   <s jzrq="20171231000000" />
	//	   <s errflag="0" />
	//	   <s ye="1002.03" />
	//	   <s typelistforr_result="czsj:d,czdd:s,czlx:s,je:n" />
	//	   <s errtext="" />
	//	   <d k="r_result">
	//	       <r czsj="20171231000000" czlx="利息收入" xh="1" czdd="" je="4.85" />
	//	       <r czsj="20171212000000" czlx="医保收入" xh="2" czdd="" je="144.51" />
	//	       <r czsj="20171113000000" czlx="医保收入" xh="3" czdd="" je="144.51" />
	//	       <r czsj="20171018000000" czlx="医保收入" xh="4" czdd="" je="144.51" />
	//	       <r czsj="20170912000000" czlx="医保收入" xh="5" czdd="" je="133.2" />
	//	       <r czsj="20170811000000" czlx="医保收入" xh="6" czdd="" je="133.2" />
	//	       <r czsj="20170711000000" czlx="医保收入" xh="7" czdd="" je="133.2" />
	//	       <r czsj="20170612000000" czlx="医保收入" xh="8" czdd="" je="133.2" />
	//	       <r czsj="20170512000000" czlx="医保收入" xh="9" czdd="" je="133.2" />
	//	       <r czsj="20170412000000" czlx="医保收入" xh="10" czdd="" je="133.2" />
	//	       <r czsj="20170409142823" czlx="医保支出" xh="11" czdd="山东大学口腔医院" je="-360.0" />
	//	       <r czsj="20170313000000" czlx="医保收入" xh="12" czdd="" je="133.2" />
	//	       <r czsj="20170226155048" czlx="医保支出" xh="13" czdd="山东大学口腔医院" je="-730.36" />
	//	       <r czsj="20170214000000" czlx="医保收入" xh="14" czdd="" je="133.2" />
	//	       <r czsj="20170212153456" czlx="医保支出" xh="15" czdd="山东大学口腔医院" je="-912.44" />
	//	       <r czsj="20170205153404" czlx="医保支出" xh="16" czdd="山东大学口腔医院" je="-331.4" />
	//	       <r czsj="20170111000000" czlx="医保收入" xh="17" czdd="" je="133.2" />
	//	       <r czsj="20161231000000" czlx="利息收入" xh="18" czdd="" je="9.83" />
	//	       <r czsj="20161227084014" czlx="医保支出" xh="19" czdd="济南平嘉大药房有限公司山影店" je="-15.0" />
	//	       <r czsj="20161212000000" czlx="医保收入" xh="20" czdd="" je="133.2" />
	//	       <r czsj="20161210212127" czlx="医保支出" xh="21" czdd="济南平嘉大药房有限公司解放东路店" je="-46.3" />
	//	       <r czsj="20161203204941" czlx="医保支出" xh="22" czdd="济南平嘉大药房有限公司解放东路店" je="-3.9" />
	//	       <r czsj="20161111000000" czlx="医保收入" xh="23" czdd="" je="133.2" />
	//	       <r czsj="20161017000000" czlx="医保收入" xh="24" czdd="" je="133.2" />
	//	       <r czsj="20160913000000" czlx="医保收入" xh="25" czdd="" je="108.0" />
	//	       <r czsj="20160810000000" czlx="医保收入" xh="26" czdd="" je="108.0" />
	//	       <r czsj="20160712000000" czlx="医保收入" xh="27" czdd="" je="108.0" />
	//	       <r czsj="20160613000000" czlx="医保收入" xh="28" czdd="" je="108.0" />
	//	       <r czsj="20160512000000" czlx="医保收入" xh="29" czdd="" je="108.0" />
	//	       <r czsj="20160412000000" czlx="医保收入" xh="30" czdd="" je="108.0" />
	//	       <r czsj="20160310000000" czlx="医保收入" xh="31" czdd="" je="108.0" />
	//	       <r czsj="20160217000000" czlx="医保收入" xh="32" czdd="" je="108.0" />
	//	       <r czsj="20160112000000" czlx="医保收入" xh="33" czdd="" je="108.0" />
	//	   </d>
	//	</p>`), `{
	//	   "code": "0000",
	//	   "message": "成功",
	//	   "data": {
	//	       "account_info": {
	//	           "balance": "$.p.s[ye].string"
	//	       },
	//	       "trans_detail":[
	//	           {
	//"balance": "$.p.s[ye].string.e%x",
	//	               "type":"$.p.d.r[czlx].string",
	//	               "num":"$.p.d.r[xh].string",
	//	               "money":"$.p.d.r[je].string"
	//	           }
	//	       ]
	//	   }
	//	}`)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	log.Println(string(xmlToJsonRes))
	//
	//
	//	json.Test(insertValues)

	//scanner := json.Scanner{}

	//json.TestScanner()

	//data, _ := ioutil.ReadFile("/home/fan/workspace/go_git/src/DataFormatUtils/test.json")
	//
	//res := make(map[string]interface{})
	//e := json2.Unmarshal(data, &res)
	//logger.Debug(e)
	//
	//err := json.TravelJsonData(data, func(path string, value interface{}) bool {
	//
	//	logger.Debug(path)
	//	logger.Debug(value)
	//	return false
	//})
	//6XvCD5VYvvWvJg8pLkvuSh5Onnfi/5yuRLkjw5V8MP2ZxvtxQ+krwrObZ/B7ZFR4ODKkk+cuXO173hq171ysPz57/6Los6MclRspN8Q/JRW+4lT1zcvW/3Vbvbs5cXZm3JTYxJBsnY3o/Mjr0zpnxiPL9aPFpsFfdcf+qpZ59KXvbsQy/7zLP3XjXXH1wtd19sVhef/7z7HniJuf/SHcbdPfiaLM1ek2bpYDBPi2rnox8cj4fD8fi9OwP7r2u/tO0/mvx8ciX5kWtX/mRm7pjesvimrb+19b1bP7b144O/O/3ZrZ+a/fzsieWvzv632e/PqvWlC/Ni+5ccfcZcO3/uYnV47vCew5cdvv7wLYfvOvyOwx84/B8Of+7wVw6L5w7N4ZkLT5qdT6wvrS4eJ0+ayUdXF9MnzfRaVh1WF6cXr17cvnjP9U9/5tnp0898ZjZ/4Qttd64/9sKnX/bMp92T61Vy/cH7L128MLSdWy2HFy9cev79DzzvPnc/mpiR6+DzRhef/1Lz4EvN9em0mI0Gg+Ho7Gy9lxdH57ZfnBZbw+H1wd5iNn3/fFFmw/HOYDAar3bX82Kdn7k4nEzuKNMzi9lZy4mqlq2fTx5I/ptrl/6O+TvL//roR67+9NHHr35297ee/wfLP7szu3vn7ruO73lk8cjZb9n6W2dH9+8/adbX0uO1Wa/v3r5y7lcsV+62Y5Za3qyTbTP4xP3rdTI880kzSe5NVrb/yfDWT5oquSMpzPQT0+HV4dbwnuvPfOYZx4BnbZ/vsSx4xvb/aXd//RmwoNtzywvb8edbvjy/uQTZc9Y8emZ/fjDeGQ52qnQ+XxwsF9V0NN3Z
	//logger.Error(err)
	//<s html="this.LODOP.SET_PRINT_STYLE(&quot;FontSize&quot;,14);this.LODOP.SET_PRINT_STYLE(&quot;FontName&quot;,&quot;黑体&quot;)" />
	xmlData := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
    <soap:Body>
        <ns2:getCurrentLoanBooksResponse xmlns:ns2="http://webservice.opac.interlib.com/">
            <return>
                <biblios>
                    <author>孔祥涛，孙先伟，刘翔宇著</author>
                    <bookSign>0</bookSign>
                    <bookrecno>45878</bookrecno>
                    <booktype>1</booktype>
                    <classNo>A752</classNo>
                    <comnum>0</comnum>
                    <isbn>978-7-5034-4301-5</isbn>
                    <publisher>中国文史出版社</publisher>
                    <title>毛泽东家风</title>
                </biblios>
                <holding>
                    <barcode>HZ0043005</barcode>
                    <callno>A752/26</callno>
                    <cirtype>HZT_001</cirtype>
                    <curlocal>HZT_009</curlocal>
                    <orglib>HZT</orglib>
                    <orglocal>HZT_009</orglocal>
                    <totalRenewNum>8</totalRenewNum>
                    <volnum>1</volnum>
                </holding>
                <loanCount>1</loanCount>
                <loanDate>2019-09-26T09:42:52+08:00</loanDate>
                <loanDateInStr>2019-09-26</loanDateInStr>
                <rdid>0000066</rdid>
                <returnDate>2019-11-25T09:42:52+08:00</returnDate>
                <returnDateInStr>2019-11-25</returnDateInStr>
                <underlease>false</underlease>
            </return>
        </ns2:getCurrentLoanBooksResponse>
    </soap:Body>
</soap:Envelope>`
	_ = xml.TraverseXmlIterative([]byte(xmlData), []*xml.Array{{4, "return"}}, func(name string, path string, value string) bool {
		logger.Debugf("name:%v", name)
		logger.Debugf("path:%v", path)
		logger.Debugf("value:%v", value)
		return false
	})

}
