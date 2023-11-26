package sdkInit

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *Application) Set(args []string) (string, error) {
	var tempArgs [][]byte
	for i := 1; i < len(args); i++ {
		tempArgs = append(tempArgs, []byte(args[i]))
	}
	
	fmt.Println("进入了set函数")
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2])}}
	response, err := t.SdkEnvInfo.ChClient.Execute(request)
	if err != nil {
		// 资产转移失败
		return "", fmt.Errorf("资产转移失败")
	}

	//fmt.Println("============== response:",response)

	return string(response.TransactionID), nil
}

func (t *Application) Path(args []string) (string, error) {
	var tempArgs [][]byte
	for i := 1; i < len(args); i++ {
		tempArgs = append(tempArgs, []byte(args[i]))
	}
	
	fmt.Println("进入了path函数")
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: tempArgs}
	response, err := t.SdkEnvInfo.ChClient.Execute(request)
	if err != nil {
		// 资产转移失败
		fmt.Println("执行execute函数失败")
		return "", err
	}

	//fmt.Println("============== response:",response)

	return string(response.Payload), nil
}


func (t *Application) Route(args []string) (string, error) {
	var tempArgs [][]byte
	for i := 1; i < len(args); i++ {
		tempArgs = append(tempArgs, []byte(args[i]))
	}
	
	fmt.Println("进入了Route函数")
	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: tempArgs}
	response, err := t.SdkEnvInfo.ChClient.Execute(request)
	if err != nil {
		// 资产转移失败
		fmt.Println("执行execute函数失败")
		return "", err
	}

	//fmt.Println("============== response:",response)

	return string(response.Payload), nil
}


// func (t *Application) Path(args []string) (string, error) {
// 	var tempArgs [][]byte
// 	for i := 1; i < len(args); i++ {
// 		tempArgs = append(tempArgs, []byte(args[i]))
// 	}
	
// 	fmt.Println("进入了path函数")
// 	request := channel.Request{ChaincodeID: t.SdkEnvInfo.ChaincodeID, Fcn: args[0], Args: tempArgs}
// 	response, err := t.SdkEnvInfo.ChClient.Execute(request)
// 	if err != nil {
// 		// 资产转移失败
// 		fmt.Println("执行execute函数失败")
// 		return "", err
// 	}

// 	//fmt.Println("============== response:",response)

// 	return string(response.Payload), nil
// }