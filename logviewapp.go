package main

import (
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"strconv"
	"strings"
)

//Declare global variables
var flagClickA = true
var flagClickB = true
var flagParm = true
var inputSeq string
var inputSeq2 string
var modeRun = "mode-1 for plan1-2-3"
var defaultInputSeq = "ctgatgaccaaactcggcctgtccgggaccacccgcggcaaagcccgcaggaccacgatcgctgatccggccAcagcccgtcccgccgatctcgtccagcgccgcttcggaccaccagcacctaaccggctgtgggtagcagacctcacctatgtgtcgacctgggcagggttcgcctacgtggcctttgtcaccgacgcctacgtcgcaggatcctgggctggcgggtcgcttccacgatggcca"
var defaultInputSeq2 = "GCGATATCTGGTGGTCTGCACGGCGTCGGCGTGTCGGTGGTTAACGCGCTATCCACCCGGCTCGAAGTCGAGATCAAGCGCGACGGGTACGAGTGGTCTCAGGTTTATGAGAAGTCGGAACCCCTGGGCCTCAAGCAAGGGGCGCCGACCAAGAAGACGGGGTCAACGGTGCGGTTCTGGGCCGACCCCGCTGTTTTCGAAACCACGG"

//var defaultInputSeq = "GCGATATCTGGTGGTCTGCACGGCGTCGGCGTGTCGGTGGTTAACGCGCTATCCACCCGGCTCGAAGTCGAGATCAAGCGCGACGGGTACGAGTGGTCTCAGGTTTATGAGAAGTCGGAACCCCTGGGCCTCAAGCAAGGGGCGCCGACCAAGAAGACGGGGTCAACGGTGCGGTTCTGGGCCGACCCCGCTGTTTTCGAAACCACGG"

// wg is used to monitor go routines

var seq1, seq2, lenUS, selfPaired, numGCUS, circleThresh, circleThreshPlan3, maxOutput, showProcess, batchInput, threshold *walk.LineEdit
var modeRuns *walk.ComboBox

//var proBar *walk.ProgressBar

var (
	//Corresponding input variable
	seq1Name              string
	batchUS               int64   = 11
	threshGC              float64 = 0.7
	maxOutSeqNum          int64   = 10
	lengthUS              int64   = 22
	lengthS4                      = 8
	lowGCUS               int64   = 13
	highGCUS              int64   = 15
	lowCircleThresh       int64   = 6
	highCircleThresh      int64   = 8
	lowCircleThreshPlan3  int64   = 0
	highCircleThreshPlan3 int64   = 0
	selfPairedNum         int64   = 5
	//progressBar = 0
)

func initText() {
	//Initialize the input variable
	seq1Name = seq1.Text()
	if seq1Name != "" {
		inputSeq = seq1Name
	} else {
		inputSeq = defaultInputSeq
	}

	if seq2.Text() != "" {
		seq2Name := seq2.Text()
		inputSeq2 = seq2Name
	} else {
		inputSeq2 = defaultInputSeq2
	}

	if threshold.Text() != "" {
		thresh := threshold.Text()
		threshGC, _ = strconv.ParseFloat(thresh, 64)
	}

	if maxOutput.Text() != "" {
		maxNum := maxOutput.Text()
		maxOutSeqNum, _ = strconv.ParseInt(maxNum, 10, 64)
	}

	if lenUS.Text() != "" {
		lensUS := lenUS.Text()
		lengthUS, _ = strconv.ParseInt(lensUS, 10, 64)
	}

	if selfPaired.Text() != "" {
		selfPaireds := selfPaired.Text()
		selfPairedNum, _ = strconv.ParseInt(selfPaireds, 10, 64)
	}

	//if lenS4.Text() != "" {
	//	lensS4 := lenS4.Text()
	//	lengthS4s, _ := strconv.ParseInt(lensS4, 10, 64)
	//	strInt64 := strconv.FormatInt(lengthS4s, 10)
	//	lengthS4, _ = strconv.Atoi(strInt64)
	//}

	if batchInput.Text() != "" {
		maxNum := batchInput.Text()
		batchUS, _ = strconv.ParseInt(maxNum, 10, 64)
	}

	if numGCUS.Text() != "" {
		threshGCUS := numGCUS.Text()
		lowGCUS, _ = strconv.ParseInt(strings.Split(threshGCUS, ",")[0], 10, 64)
		highGCUS, _ = strconv.ParseInt(strings.Split(threshGCUS, ",")[1], 10, 64)
	}

	if modeRuns.Text() != "" {
		threshGCUS := modeRuns.Text()
		modeRun = threshGCUS
	}

	if circleThresh.Text() != "" {
		circleThreshs := circleThresh.Text()
		lowCircleThresh, _ = strconv.ParseInt(strings.Split(circleThreshs, ",")[0], 10, 64)
		highCircleThresh, _ = strconv.ParseInt(strings.Split(circleThreshs, ",")[1], 10, 64)
	}

	if circleThreshPlan3.Text() != "" {
		circleThreshs := circleThreshPlan3.Text()
		lowCircleThreshPlan3, _ = strconv.ParseInt(strings.Split(circleThreshs, ",")[0], 10, 64)
		highCircleThreshPlan3, _ = strconv.ParseInt(strings.Split(circleThreshs, ",")[1], 10, 64)
	}

	if lowGCUS < 0 || highGCUS > lengthUS || highGCUS < lowGCUS {
		log.Println("Error, please check your GC number of US and restart program!")
		flagParm = false
	}

	if lowCircleThresh < 0 || highCircleThresh > lengthUS || highCircleThresh < lowCircleThresh {
		log.Println("Error, please check your loop gap and restart program!")
		flagParm = false
	}

	if batchUS < 0 || batchUS > int64(numberBatchs) {
		log.Println("Error, please check your batchUS and restart program! (0 <= batchUS <= 2000)")
		flagParm = false
	}

	if threshGC < 0 || threshGC > 1 {
		log.Println("Error, please check your threshold and restart program!")
		flagParm = false
	}

	if maxOutSeqNum <= 0 {
		log.Println("Error, please check your maxSqeNum and restart program!")
		flagParm = false
	}

	if selfPairedNum <= 0 {
		log.Println("Error, please check your selfPairedNum and restart program!")
		flagParm = false
	}

	if lengthUS <= 0 || lengthS4 < 0 {
		log.Println(lengthS4)
		log.Println(lengthS4 <= 0)
		log.Println("Error, please check your length of US or 4S and restart program!")
		flagParm = false
	}

	//if threshGC > 0 && threshGC < 1 && maxOutSeqNum > 0 && lengthS4 >= 0 && lengthUS > 0 && lowGCUS > 0 && highGCUS < lengthUS && highGCUS > lowGCUS && lowCircleThresh > 0 && highCircleThresh < lengthUS && highCircleThresh > lowCircleThresh {
	//	flagParm = true
	//}

	inputSeq = transa2A(inputSeq)
	inputSeq2 = transa2A(inputSeq2)

	seq1.SetText(inputSeq)
	seq2.SetText(inputSeq2)

}

func main() {
	var mw *walk.MainWindow

	if err := (MainWindow{
		AssignTo: &mw,
		//Icon:     "sjtu2.ico",
		Title:   "Primer Designing by SJTU Biomedical Imaging Informatics Lab",
		MinSize: Size{850, 600},
		Layout:  VBox{},

		Children: []Widget{
			//VSplitter{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{

					Label{MaxSize: Size{150, 50},
						Text: "The 1st sequence",
					},
					LineEdit{AssignTo: &seq1, Text: defaultInputSeq},

					Label{MaxSize: Size{150, 50},
						Text: "The 2nd sequence",
					},
					LineEdit{AssignTo: &seq2, Text: defaultInputSeq2},

					Label{MaxSize: Size{150, 50},
						Text: "Progress display",
					},
					LineEdit{AssignTo: &showProcess, Text: "0", ReadOnly: true},
				},
			},

			//ProgressBar{
			//	AssignTo: &proBar,
			//	MaxValue: 10000,
			//	MinValue: 0,
			//	Value: progressBar,
			//},

			Composite{
				Layout: Grid{Columns: 6},
				Children: []Widget{

					Label{MaxSize: Size{150, 50},
						Text: "The value of the GC threshold",
					},
					LineEdit{AssignTo: &threshold, Text: fmt.Sprintf("%.2f", threshGC)},

					Label{MaxSize: Size{150, 50},
						Text: "The length of US",
					},
					LineEdit{AssignTo: &lenUS, Text: fmt.Sprintf("%d", lengthUS)},

					Label{MaxSize: Size{150, 50},
						Text: "The GC number section of US",
					},
					LineEdit{AssignTo: &numGCUS, Text: fmt.Sprintf("%d,%d", lowGCUS, highGCUS)},

					Label{MaxSize: Size{150, 50},
						Text: "The gap of loop in plan1-2",
					},
					LineEdit{AssignTo: &circleThresh, Text: fmt.Sprintf("%d,%d", lowCircleThresh, highCircleThresh)},

					Label{MaxSize: Size{150, 50},
						Text: "The gap of loop in plan3",
					},
					LineEdit{AssignTo: &circleThreshPlan3, Text: fmt.Sprintf("%d,%d", lowCircleThreshPlan3, highCircleThreshPlan3)},

					Label{MaxSize: Size{150, 50},
						Text: "The max number of the output",
					},
					LineEdit{AssignTo: &maxOutput, Text: fmt.Sprintf("%d", maxOutSeqNum)},

					Label{MaxSize: Size{150, 50},
						Text: "The batches of US (0-100)",
					},
					LineEdit{AssignTo: &batchInput, Text: fmt.Sprintf("%d", batchUS)},

					Label{MaxSize: Size{150, 50},
						Text: "The serial num of self-paired",
					},
					LineEdit{AssignTo: &selfPaired, Text: fmt.Sprintf("%d", selfPairedNum)},

					Label{
						Text: "The running mode for different plan",
					},
					ComboBox{
						Editable: true,
						AssignTo: &modeRuns,
						Model:    []string{"mode-1 for plan1-2-3", "mode-2 for plan4"},
					},
				},
			},


			PushButton{
				Text: "Analysis input sequence",
				OnClicked: func() {
					if flagClickA {

						log.Println("############################################################################################")
						log.Println("The program start...")
						initText()
						log.Println("The run mode:")
						log.Println(modeRun)
						log.Println("The input 1st sequence:")
						log.Println(inputSeq)
						log.Println("The input 2nd sequence:")
						log.Println(inputSeq2)
						log.Println("The input GC threshold:")
						log.Println(threshGC)
						log.Println("The input length of US:")
						log.Println(lengthUS)
						log.Println("The lowest and highest GC number of US:")
						log.Println(lowGCUS, highGCUS)
						log.Println("The lowest and highest gap of loop in plan1-2:")
						log.Println(lowCircleThresh, highCircleThresh)
						log.Println("The lowest and highest gap of loop in plan3:")
						log.Println(lowCircleThreshPlan3, highCircleThreshPlan3)
						log.Println("The input max number of the output sequence:")
						log.Println(maxOutSeqNum)
						log.Println("The batches of US:")
						log.Println(batchUS)
						log.Println("############################################################################################")

						if flagParm {
							flagClickA = false
							log.Println("############################################################################################")
							log.Println("Strat to analysis the input sequence:")
							if modeRun == "mode-1 for plan1-2-3" {
								//go func() {
								analysisInput(inputSeq)
								//}()
								log.Println("Find paired sequence between 1S and 2A in plan1-2:", len(input1S))
								log.Println("Find paired sequence between 1S, 2A, 3S and 4S in plan3:", len(input3S1))
								log.Println("############################################################################################")

							} else {
								analysisInputPlan4(inputSeq, inputSeq2)
								log.Println("Find paired sequence between 1S and 2A in plan4:", len(input4S1))
								log.Println("Find paired sequence between 3S and 4A in plan4:", len(input4S3))
								log.Println("############################################################################################")
							}

						} else {
							log.Println("The program detects parameter exceptions. Please check the input and restart the program!")
						}

					} else {
						log.Println("Please do not click while the program is running. If want to restart, please restart the program.", "\r")
					}
				},
			},

			PushButton{
				Text: "Run for making primers",
				OnClicked: func() {
					if flagClickB {
						if flagParm {
							flagClickB = false
							if flagClickA == true {
								log.Println("############################################################################################")
								log.Println("The program start...")
								initText()
								log.Println("The run mode:")
								log.Println(modeRun)
								log.Println("The input 1st sequence:")
								log.Println(inputSeq)
								log.Println("The input 2nd sequence:")
								log.Println(inputSeq2)
								log.Println("The input GC threshold:")
								log.Println(threshGC)
								log.Println("The input length of US:")
								log.Println(lengthUS)
								log.Println("The lowest and highest GC number of US:")
								log.Println(lowGCUS, highGCUS)
								log.Println("The lowest and highest gap of loop in plan1-2:")
								log.Println(lowCircleThresh, highCircleThresh)
								log.Println("The lowest and highest gap of loop in plan3:")
								log.Println(lowCircleThreshPlan3, highCircleThreshPlan3)
								log.Println("The input max number of the output sequence:")
								log.Println(maxOutSeqNum)
								log.Println("The batches of US:")
								log.Println(batchUS)
								log.Println("############################################################################################")
							}

							if modeRun == "mode-1 for plan1-2-3" {
								if flagClickA == true {
									analysisInput(inputSeq)
									log.Println("Find paired sequence between 1S and 2A in plan1-2:", len(input1S))
									log.Println("Find paired sequence between 1S, 2A, 3S and 4S in plan3:", len(input3S1))
									log.Println("############################################################################################")

								}
								getOutOfPlan3()
								go func() {
									makeUSOdd()
								}()
							} else {
								if flagClickA == true {
									analysisInputPlan4(inputSeq, inputSeq2)
									log.Println("Find paired sequence between 1S and 2A in plan4:", len(input4S1))
									log.Println("Find paired sequence between 3S and 4A in plan4:", len(input4S3))
									log.Println("############################################################################################")

								}
								go func() {
									getOutOfPlan4()
								}()
							}
							log.Println("############################################################################################")

						} else {
							log.Println("The program detects parameter exceptions. Please check the input and restart the program!")
						}

					} else {
						log.Println("Please do not click while the program is running. If want to restart, please restart the program.", "\r")
					}

					//flagClickA = true
					//flagClickB = true
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	lv, err := NewLogView(mw)
	if err != nil {
		log.Fatal(err)
	}

	// //Let the interface take over the log output
	log.SetOutput(lv)
	lv.PostAppendText("All components have been loaded, please set the parameters before click RUN to start the program.")
	mw.Run()

}
