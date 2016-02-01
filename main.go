package main

import (
//"bytes"
	"fmt"
	"strings"
//"go/types"
//"time"
	"time"
	"os/exec"
	"os"
//"math/big"
	"strconv"
)

var timeDelayed = false
var showBoardGlobal = true
var showCommentsGlobal = false
type sudoku [81] cell

type cell struct {
	horizontal rune
	vertical   rune
	stack      int
	value      string
	number     int
	options    string
}
var solveLog string

func main() {
	/*timeDelayed = true
	showBoard = false
	showComments = true
*/ //blablabal

	start:= time.Now()
	//lett board := setupBoard("53**7****6**195****98****6*8***6***34**8*3**17***2***6*6****28****419**5****8**79")
	//hardest board := setupBoard("8**********36******7**9*2***5***7*******457*****1***3***1****68**85***1**9****4**")
	//board := setupBoard("******5***8***4*6*3*4*6*7***1*2*3*****9***4*****7*6*5***5*8*9*2*6*1***8***2******") //hard
	//board := setupBoard("*4***6***9******41**8**9*5**9***7**2**3***8**4**8***1**8*3**9**16******7***5***8*") //impossible  klarte isje  2 solutions

	//board := setupBoard("*4***6***9*6****41**8**9*5**9***7**2**3***8**4**8***1**8*3**9**16******7***5***8*") //impossible  men lagt in 6 i b3
	//board := setupBoard("*4***6*989*6****41**8**9*5*89***7**26*3***87*4**8***1**8*3**9**16******7***5**18*") //impossible prøbd på
	board := setupBoard("**2*********8**291**6**7*43**4*5*9***3*9*4128*1**63*7**697*23***2**315***8***576*") //kim sin
	//board := setupBoard("**3*****1***76*4**5****2*9***7****48*5*3*9*2*49****6***3*8****4**6*27***8*****2**") //evil
	//board := setupBoard("*******1**189***6***6*814**8*3****245**3*8**627****5*8**251*7***8***324**9*******")
	//board := setupBoard("*********************************************************************************")

	//printBoard(board)

	//board,log, solved  := solve(board,false,false,false)
	board,log, solved  := guess(board)
	timeUsed := time.Since(start)


	if(true){
		fmt.Println(log)
	}

	if(solved){
		if(verifySudoku(board)){
			fmt.Println("verified")
			fmt.Println(" Solved it in: " +timeUsed.String() + "\n")

		}
	}else{
		fmt.Println("did not do it")
	}
	printBoard(board)
	//fmt.Println("Ran with timedelay of 100ms each solved cell")




}

func verifySudoku(board sudoku) bool{
	for y := 0; y < 81; y++ { // for each cell in the sudoku board

		for i := 0; i < 81; i++ {   // go through the other cells and check if it has the same value if its in the same groups.
			if(isInSameGroup(y,i,board) && y!=i){
				if(board[y].value == board[i].value){
					return false
				}
			}
		}
	}
	return true
}
func solve(board sudoku,slow bool, showBoard bool, showComments bool) (sudoku, string, bool){
	timeDelayed = slow
	showBoardGlobal = showBoard
	showCommentsGlobal = showComments

	var solved bool
	var solveCount int = 1
	var optionChangeCount int = 0
	var numberOfTriesToEasySolve int  =0

	var madeProgress bool = true
	for madeProgress && !solved {
		madeProgress = false
		solveCount = 1
		for (solveCount>0 && solved == false){
			board,solved,solveCount = easySolve(board)
			numberOfTriesToEasySolve++
			/*fmt.Printf("%v", "\n ")
			fmt.Printf("%v", solveCount)
			fmt.Printf("%v", " solved\n")*/

			if(solveCount>0){madeProgress = true}
		}

		if(!solved){
			board,solved,solveCount = solveByLookingAtGroupCellsOptions(board) // do this once
			/*fmt.Printf("%v", solveCount)
			fmt.Printf("%v", " solved\n")*/
			//fmt.Println("kjørte")
			//fmt.Println(solveCount)
			if(solveCount>0){madeProgress = true}
		}

		if(!solved){
			board,optionChangeCount = solveByLookingAtPairsWithTheSameTwoOptionsAndRemovingOptionOnOtherCells(board)
			if(optionChangeCount>0){madeProgress = true}

		}


	}
 return board, solveLog, solved
}

func exportSudoko(board sudoku) string {
	export :=""
	for y := 0; y < 81; y++ { // for each cell in the sudoku board
	export+= board[y].value
	}
	return export
}

func guess(originBoard sudoku) (sudoku, string, bool){
	guessCounter :=0
	//make an attempt to solve
	originBoard,originLog,solved :=solve(originBoard,false,false,false)
	if(!solved){//if not solved start guessing

		board:=originBoard
		logBeforeNewLog:=""

		for y := 0; y < 81; y++ { // for each cell in the sudoku board
			if ( board[y].value == "*" ) {//if the cell is not solved
				for i:=0; i <len(board[y].options); i++{
					solveLog=""
					board=originBoard
					logBeforeNewLog:=originLog+"Not so easy to solve.  Probably has multiple solutions. We need to guess\n"
					guessCounter++
					//fmt.Println("i er:"+ strconv.Itoa(i))
					//fmt.Println("y er:")
					//fmt.Println(y)
					//fmt.Println("options er:"+string(board[y].options))
					board[y].value=string(board[y].options[i])
					board[y].options=""

					//board = removeFromCellOptionsInAllFriendlyGroups(string(board[y].options[i]),y,board)
					logBeforeNewLog = logBeforeNewLog+"Guessing "+ getCellName(board[y])+ " is "+ board[y].value + "\n"
					//fmt.Println(getCellName(board[y]) + " is maybe:"+board[y].value+" \n")

					var log string

					board,log,solved =solve(board,false,false,false)

					if(solved){
						if(verifySudoku(board)){
							log = logBeforeNewLog+log
							log = log+"I guessed "+strconv.Itoa(guessCounter)
							//log = log+string(guessCounter)  doesnt work
							log = log+" times\n"
							log = log+"But when i guessed that "+ getCellName(board[y])+ " was "+ board[y].value + ", well then the rest was easy"
							return board,log,solved
						}else{
							board[y].value="*"
							board=originBoard
						}

					} else{
						board[y].value="*"
						board=originBoard
					}
				}
			}
		}
		logBeforeNewLog = logBeforeNewLog+ "Failed finding the solution after guessing ? times"
		return originBoard, logBeforeNewLog , false
	}else{
		return originBoard, originLog , true
	}
}


type bin time.Duration

func (b bin) String() string {
	return fmt.Sprintf("%b", b)
}

func getCellName(cell cell) (string){
	return string(cell.horizontal) + string(cell.vertical)
}
func printOptionChange(board sudoku,i int,y int, o int, value string){
			printBoardAndDelay(board)
			printComment( "Removed "+ value + " from options on " +getCellName(board[o]) +". Reason: "+getCellName(board[i])+" and "+ getCellName(board[y]))
			printComment(" are paired with options: "+string(board[i].options[0])+""+string(board[i].options[1])+". \n" )
			printComment(""+  getCellName(board[o]) +" have now options: "+string(board[o].options)+" \n" )
}

//help function.
func removeTwoValuesFromCellOptionsIfItCountainsTheseValuesAndPrint(board sudoku,i int,y int, o int, optionChangeCount int) (sudoku, int){
	var changed bool = false
	board[o] , changed = removeFromCellOptionsIfItCountainsThisValue(board[o],string(board[i].options[0]))
	if(changed){
		printOptionChange(board,i,y,o,string(board[i].options[0]))
		optionChangeCount++
	}

	board[o] , changed = removeFromCellOptionsIfItCountainsThisValue(board[o],string(board[i].options[1]))
	if(changed){
		printOptionChange(board,i,y,o, string(board[i].options[1]))
		optionChangeCount++
	}
	return board, optionChangeCount
}

func removeFromCellOptionsIfItCountainsThisValue(cell cell, value string) (cell, bool){
	if(strings.Contains(cell.options, value)){
		cell.options = strings.Replace(cell.options, value, "", -1)
		return cell, true
	}	else {
		return cell,false
	}

}


func solveByLookingAtPairsWithTheSameTwoOptionsAndRemovingOptionOnOtherCells(board sudoku) (sudoku, int) {

	optionChangeCount := 0
	for y := 0; y < 81; y++ { // for each cell in the sudoku board
		if ( board[y].value == "*" ) {//if the cell is not solved

			if(len(board[y].options) == 2){
				for i := 0; i < 81; i++ {   // go through the other cells to see their possibilities

					if (board[y].horizontal == board[i].horizontal && board[y].number != board[i].number) {  //compare with cells in same horizontal
						if (board[i].options==board[y].options) {
							//then you can remove these two options from all horizontal
							for o := 0; o < 81; o++ {
								if((o!=i && o !=y) && board[i].horizontal == board[o].horizontal){ //dont remove from the pair
									board, optionChangeCount = removeTwoValuesFromCellOptionsIfItCountainsTheseValuesAndPrint(board,i,y,o,optionChangeCount)
								}
							}
						}
					}
					if (board[y].vertical == board[i].vertical && board[y].number != board[i].number) {  //compare with cells in same vertical
						if (board[i].options==board[y].options) {
							for o := 0; o < 81; o++ {
								if((o!=i && o !=y) && board[i].vertical == board[o].vertical){

									board, optionChangeCount = removeTwoValuesFromCellOptionsIfItCountainsTheseValuesAndPrint(board,i,y,o,optionChangeCount)
								}
							}
						}
					}
					if (board[y].stack == board[i].stack && board[y].number != board[i].number) {   //compare with cells in same stack
						if (board[i].options==board[y].options) {
							for o := 0; o < 81; o++ {
								if((o!=i && o !=y) && board[i].stack == board[o].stack){
									board, optionChangeCount = removeTwoValuesFromCellOptionsIfItCountainsTheseValuesAndPrint(board,i,y,o,optionChangeCount)
								}
							}
						}
					}
				}
			}

		}
	}
	return board , optionChangeCount
}

func solveByLookingAtGroupCellsOptions(board sudoku) (sudoku, bool, int) {
	figuredOutCell := false
	solveCount :=0
	sudokuComplete := true
	for y := 0; y < 81; y++ { // for each cell in the sudoku board
		figuredOutCell = false
		if ( board[y].value == "*" ) {//if the cell is not solved
			horizontalCellsOptions := ""
			verticalCellsOptions:= ""
			stackCellsOptions := ""
			for i := 0; i < 81; i++ {   // go through the other cells to see their possibilities


				if (board[y].horizontal == board[i].horizontal && board[y].number != board[i].number) {  //compare with cells in same horizontal
					horizontalCellsOptions = horizontalCellsOptions + board[i].options
				}
				if (board[y].vertical == board[i].vertical && board[y].number != board[i].number) {  //compare with cells in same vertical
					verticalCellsOptions = verticalCellsOptions +board[i].options
				}
				if (board[y].stack == board[i].stack && board[y].number != board[i].number) {   //compare with cells in same stack
					stackCellsOptions = stackCellsOptions +board[i].options
				}


			}

			//look at each option you have left and compare with groups, if is is missing from a group it is the correct one.
			for _, option := range board[y].options {
				if(board[y].value =="*"){
					if (!strings.ContainsRune(horizontalCellsOptions, option) && len(horizontalCellsOptions) > 1) {

						printBoardAndDelay(board)
						printComment(getCellName(board[y]) + " is " + string(option) + ". Reason: Is horisontally only possible here. \n")

						board[y].value = string(option)
						board[y].options = ""
						board = removeFromCellOptionsInAllFriendlyGroups(board[y].value,y,board)
						figuredOutCell = true
						solveCount++
					}else if (!strings.ContainsRune(verticalCellsOptions, option) && len(verticalCellsOptions) > 1) {

						printBoardAndDelay(board)
						printComment(getCellName(board[y]) + " is " + string(option) + ". Reason: Is vertically only possible here. \n")

						board[y].value = string(option)
						board[y].options = ""
						board = removeFromCellOptionsInAllFriendlyGroups(board[y].value,y,board)
						figuredOutCell = true
						solveCount++
					}else if (!strings.ContainsRune(stackCellsOptions, option) && len(stackCellsOptions) > 1) {

						printBoardAndDelay(board)
						printComment(getCellName(board[y]) + " is " + string(option) + ". Reason: Is only possible here in this stack  \n")
						board[y].value = string(option)
						board[y].options = ""
						board = removeFromCellOptionsInAllFriendlyGroups(board[y].value,y,board)
						figuredOutCell = true
						solveCount++
					}
				}
			}
			if(!figuredOutCell){
				//there is still an unsolved cell here
				sudokuComplete= false

			}
		}//if it was not solved

	}
	return board , sudokuComplete, solveCount
}

func removeFromCellOptionsInAllFriendlyGroups(option string, cellNumber int,board sudoku) sudoku{
	for y := 0; y < 81; y++ { // for each cell in the sudoku board
		if(isInSameGroup(y,cellNumber, board)){
			var changed bool
			board[y], changed = removeFromCellOptionsIfItCountainsThisValue(board[y],option)
			if(changed){
				//do absolutely nothing
			}
		}
	}
	return board
}

func isInSameGroup(number1 int, number2 int, board sudoku) bool{
	if(board[number1].horizontal == board[number2].horizontal || board[number1].vertical == board[number2].vertical || board[number1].stack == board[number2].stack){
		return true
	}else{
		return false
	}
}

func easySolve(board sudoku) (sudoku, bool, int) {
	var solveCount int = 0
	var sudokuComplete bool = true
	figuredOutCell := false

	for y := 0; y < 81; y++ { // for each cell in the sudoku board
		figuredOutCell = false

		if ( board[y].value == "*" ) {//if the cell is not solved
			for i := 0; i < 81; i++ {   // go through the other cells to compare
				if ((board[y].horizontal == board[i].horizontal || board[y].vertical == board[i].vertical || board[y].stack == board[i].stack) &&board[y].number != board[i].number) {  //only compare with cells that are related , and not the same
					if (board[i].value != "*") {//if this cell is solved, remove the value from options.
						board[y].options = strings.Replace(board[y].options, board[i].value, "", -1)
						if (len(board[y].options) == 1) { //if options now are only one,  set that as value and break

							board[y].value = board[y].options
							board[y].options = ""
							figuredOutCell = true
							solveCount++
							printBoardAndDelay(board)
							printComment(  getCellName(board[y]) +" is "+ string(board[y].value)+". Reason: "+ string(board[y].value)+" is the only option left when comparing its groups\n" )

							break
						}
					}
				}

			}
			//if it was not solved
			if(!figuredOutCell){
				//there is still an unsolved cell here
				sudokuComplete= false

			}
		}else{
			board[y].options = ""
		}

	}

	return board , sudokuComplete, solveCount
}


func setupBoard(initials string) sudoku {
	var board sudoku

	counter := 1
	for i := 'A'; i <= 'I'; i++ {
		for y := '1'; y <= '9'; y++ {
			stack := getStackFromNumber(counter)
			value := string([]rune(initials)[counter - 1])
			var options string
			if (value == "*" ) { options = "123456789" }
			board [counter - 1] = cell{i, y, stack, value, counter, options}
			counter++
		}
	}
	return board
}

func getStackFromNumber(number int) int {
	stack1 := []int{1, 2, 3, 10, 11, 12, 19, 20, 21}
	stack2 := []int{4, 5, 6, 13, 14, 15, 22, 23, 24}
	stack3 := []int{7, 8, 9, 16, 17, 18, 25, 26, 27}
	stack4 := []int{28, 29, 30, 37, 38, 39, 46, 47, 48}
	stack5 := []int{31, 32, 33, 40, 41, 42, 49, 50, 51}
	stack6 := []int{34, 35, 36, 43, 44, 45, 52, 53, 54}
	stack7 := []int{55, 56, 57, 64, 65, 66, 73, 74, 75}
	stack8 := []int{58, 59, 60, 67, 68, 69, 76, 77, 78}
	stack9 := []int{61, 62, 63, 70, 71, 72, 79, 80, 81}
	if contains(stack1, number) {return 1}
	if contains(stack2, number) {return 2}
	if contains(stack3, number) {return 3}
	if contains(stack4, number) {return 4}
	if contains(stack5, number) {return 5}
	if contains(stack6, number) {return 6}
	if contains(stack7, number) {return 7}
	if contains(stack8, number) {return 8}
	if contains(stack9, number) {return 9}
	return 99
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func printComment(comment string){
	solveLog = solveLog+comment
	if(showCommentsGlobal){
		fmt.Printf("%v", comment )
	}
}

func printBoardAndDelay(board sudoku) {
	if(timeDelayed){
		time.Sleep(time.Millisecond * 100)
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	if(showBoardGlobal){
	printBoard(board)
	}
}
func printBoard(board sudoku) {
	newline := 1

	for y := 0; y < 81; y++ {

		fmt.Printf("%v", board[y].value + " ")
		if (newline == 9) {
			newline = 0
			fmt.Printf("%v", "\n")
		}
		newline++
	}
	fmt.Println("", "\n\n\n\n")
}