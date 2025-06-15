package bs

import (
	state "CommonDevP1/Server/SimpleState"
	"fmt"
	"strings"

	shell "github.com/c-bata/go-prompt"
)

var s_main = &S_Main{cmdsLst: []ICmd{
	cmd_1Cylinder, cmd_2Cylinder, cmd_1Test,
}}
var s_1Cylinder = &S_1Cylinder{}
var s_2Cylinder = &S_2Cylinder{}
var s_1Test = &S_1Test{}

var _StateMgr = state.NewStateMgr()

func RunState() {
	_StateMgr.Run(s_main)
}

func dupStrAry2InterfaceAry(trg []interface{}, src []string) {
	for i, v := range src {
		trg[i] = v
	}
}

type S_Main struct {
	cmdsLst []ICmd
}

func (this *S_Main) completer(d shell.Document) []shell.Suggest {

	var s []shell.Suggest
	for _, v := range this.cmdsLst {
		s = append(s, shell.Suggest{
			Text: v.Text(), Description: v.Desc(),
		})
	}

	// s := []shell.Suggest{
	// 	{Text: cmd_1Cylinder.Text(), Description: cmd_1Cylinder.Desc()},
	// 	{Text: cmd_2Cylinder.Text(), Description: cmd_2Cylinder.Desc()},
	// 	{Text: cmd_1Test.Text(), Description: cmd_1Test.Desc()},
	// }
	return shell.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
func (this *S_Main) SwitchIn(context []interface{}) {
	for {
		cmdWithArgs := shell.Input("input task command> ", this.completer)
		cmdAry := strings.Split(cmdWithArgs, " ")

		/*
			match := false
			for _, v := range this.cmdsLst {
				if cmdAry[0] == v.Text() {
					if v.ArgsCheck(cmdAry[1:]) {
						_StateMgr.Switch(v, cmdAry[1:])
					} else {
						fmt.Println("cmd params error")
					}
					match = true
				}
			}

			if !match {
				fmt.Printf("cmd input error, try again")
			}
		*/
		switch cmdAry[0] {
		case cmd_1Cylinder.Text():
			trgC := make([]interface{}, len(cmdAry)-1)
			dupStrAry2InterfaceAry(trgC, cmdAry[1:])

			fmt.Println("main state 1")

			_StateMgr.Switch(s_1Cylinder, trgC)
			return
		case cmd_2Cylinder.Text():
			trgC := make([]interface{}, len(cmdAry)-1)
			dupStrAry2InterfaceAry(trgC, cmdAry[1:])
			fmt.Println("main state 2")
			_StateMgr.Switch(s_2Cylinder, trgC)
			return
		case cmd_1Test.Text():
			_StateMgr.Switch(s_1Test, nil)
			return
		case "":
		default:
			fmt.Printf("%+v \r\n", cmdAry)
			fmt.Printf("input error, try again")
		}
	}
}
func (this *S_Main) SwitchOut() {
	fmt.Println("s_main switch out")
}

//------------------
type S_1Cylinder struct {
}

func (this *S_1Cylinder) completer(d shell.Document) []shell.Suggest {
	s := []shell.Suggest{}
	return shell.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
func (this *S_1Cylinder) SwitchIn(context []interface{}) {
	cmdWithArgs := shell.Input("input force,times > ", this.completer)
	//cmdAry := strings.Split(cmdWithArgs, " ")
	if err := cmd_1Cylinder.Resolve(cmdWithArgs); err != nil {
		fmt.Printf("cmd_1Cylinder.Resolve failed: %s \r\n", err.Error())
	}

	_StateMgr.Switch(s_main, nil)
}
func (this *S_1Cylinder) SwitchOut() {
	fmt.Println("S_1Cylinder switch out")
}

//------------------
type S_1Test struct {
}

func (this *S_1Test) completer(d shell.Document) []shell.Suggest {
	s := []shell.Suggest{}
	return shell.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
func (this *S_1Test) SwitchIn(context []interface{}) {
	//cmdWithArgs := shell.Input("input force,times > ", this.completer)
	//cmdAry := strings.Split(cmdWithArgs, " ")
	if err := cmd_1Test.Resolve(""); err != nil {
		fmt.Printf("cmd_1Test.Resolve failed: %s \r\n", err.Error())
	}

	_StateMgr.Switch(s_main, nil)
}
func (this *S_1Test) SwitchOut() {
	fmt.Println("S_1Test switch out")
}

//------------------
type S_2Cylinder struct {
}

func (this *S_2Cylinder) completer(d shell.Document) []shell.Suggest {
	s := []shell.Suggest{
		{Text: cmd_1Cylinder.Text(), Description: cmd_1Cylinder.Desc()},
		{Text: cmd_2Cylinder.Text(), Description: cmd_2Cylinder.Desc()},
	}
	return shell.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}
func (this *S_2Cylinder) SwitchIn(context []interface{}) {
	//t := shell.Input("> ", this.completer)
	cmdWithArgs := shell.Input("input force1,force2,times > ", this.completer)
	cmd_2Cylinder.Resolve(cmdWithArgs)

	_StateMgr.Switch(s_main, nil)
}
func (this *S_2Cylinder) SwitchOut() {
	fmt.Println("S_2Cylinder switch out")

}
