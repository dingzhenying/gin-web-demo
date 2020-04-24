package rpc

import (
	"errors"
	"fmt"
)

/**
服务器端需要注册结构体对象(传输体，传入、传出对象载体)，然后通过对象所属的方法暴露给调用者，从而提供服务
go语言官方给出的对外暴露的服务方法的定义标准，其中包含了主要的几条规则，分别是：
1、对外暴露的方法有且只能有两个参数，这个两个参数只能是输出类型或内建类型，两种类型中的一种。
2、方法的第二个参数必须是指针类型。！！！！！
3、方法的返回类型为error。
4、方法的类型是可输出的。
5、方法本身也是可输出的。
*/

func (mu MathUtil) RpcFunction(request int, response *int) error {
	*response = request * request //圆形的面积 s = π * r * r
	return nil                    //返回类型
}

/**
1、Calculate方法是服务对象MathUtil向外提供的服务方法，该方法用于接收传入的圆形半径数据，计算圆形面积并返回。
2、第一个参数request代表的是调用者（client）传递提供的参数。
3、第二个参数response代表要返回给调用者的计算结果，必须是指针类型。
4、正常情况下，方法的返回值为是error，为nil。如果遇到异常或特殊情况，则error将作为一个字符串返回给调用者，此时，resp参数就不会再返回给调用者。
*/

type MathUtil struct {
	out int
}

//函数必须是导出的(首字母大写)
//必须有两个导出类型的参数，
//第一个参数是接收的参数，第二个参数是返回给客户端的参数，第二个参数必须是指针类型的
//函数还要有一个返回值error
//标准方法体1（取乘积）
func (mu *MathUtil) RpcFunction1(args *Args, reply *int) error {
	//计算面积
	*reply = args.A * args.B
	fmt.Println("RpcFunction1 方法执行: Multiply=== ", reply)

	return nil
}

//标准方法体2（取整）
func (t *MathUtil) RpcFunction2(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	//取整数
	quo.Quo = args.A / args.B
	//取余
	quo.Rem = args.A % args.B
	fmt.Println("RpcFunction2 方法执行 quo==", quo)

	return nil
}

//入参结构体
type Args struct {
	A, B int
}

//出参结构体
type Quotient struct {
	Quo, Rem int
}
