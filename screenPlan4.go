package main

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

var (
	matchStrings41 = make([]string, 0, maxOutSeqNum)
	matchStrings42 = make([]string, 0, maxOutSeqNum)
	matchStrings43 = make([]string, 0, maxOutSeqNum)
	matchStrings44 = make([]string, 0, maxOutSeqNum)
	numberTotalMatches4 int64 = 0
)

func getOutOfPlan4() {

	timeBatchBegin := time.Now()
	allSeqsNeeded := len(input4S1) * len(input4S3)

	for i := 0; i < len(input4S1); i++ {
		for j := 0; j < len(input4S3); j++ {

			allSeqsLooked = i * len(input4S3) + j

			if allSeqsLooked % 1 == 0 {
				timeBatchElapsed := time.Since(timeBatchBegin).Seconds()
				valuePercent := float64(allSeqsLooked) / float64(allSeqsNeeded)
				timeLeft := (timeBatchElapsed / valuePercent) / 3600
				valuePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", valuePercent*100), 64)
				showProcessSen := fmt.Sprintf("Checked %d in %d allSeqs %v%% in %.1f s, find %d plan-4 targets, time left: %.1f h \n", allSeqsLooked, allSeqsNeeded, valuePercent, timeBatchElapsed, numberTotalMatches4, timeLeft)
				showProcess.SetText(showProcessSen)
				//proBar.SetValue(int(valuePercent*100))
				// log.Println("checked!!!")
			}

			if numberTotalMatches4 > maxOutSeqNum {
				goto Loop
			}

			tempSrc1S := input4S1[i]
			tempSrc2A := input4A2[i]
			tempSrc3S := input4S3[j]
			tempSrc4A := input4A4[j]

			tempSrc1 := joinSeq(joinSeq(tempSrc3S, tempSrc2A),tempSrc1S)
			tempSrc2 := joinSeq(joinSeq(tempSrc1S, tempSrc4A),tempSrc2A)
			tempSrc3 := joinSeq(joinSeq(tempSrc4A, tempSrc2A),tempSrc3S)
			tempSrc4 := joinSeq(joinSeq(tempSrc3S, tempSrc1S),tempSrc4A)

			//if screeningFourSeqs(tempSrc1) || screeningFourSeqs(tempSrc2) || screeningFourSeqs(tempSrc3) || screeningFourSeqs(tempSrc4) {
			//	continue
			//}

			if screeningContinueSeqs(tempSrc1, int(selfPairedNum)) || screeningContinueSeqs(tempSrc2, int(selfPairedNum)) || screeningContinueSeqs(tempSrc3, int(selfPairedNum)) || screeningContinueSeqs(tempSrc4, int(selfPairedNum)) {
				continue
			}

			if pair5BetweenSeq(tempSrc1S, tempSrc3S) || pair5BetweenSeq(tempSrc1S, tempSrc4A) || pair5BetweenSeq(tempSrc3S, tempSrc2A) || pair5BetweenSeq(tempSrc4A, tempSrc2A) {
				continue
			}

			numberTotalMatches4++

			matchStrings41 = append(matchStrings41, tempSrc1)
			matchStrings42 = append(matchStrings42, tempSrc2)
			matchStrings43 = append(matchStrings43, tempSrc3)
			matchStrings44 = append(matchStrings44, tempSrc4)
		}
	}
Loop:
	timeBatchElapsed := time.Since(timeBatchBegin).Seconds()
	showProcessSen := fmt.Sprintf("Done in %.1f s \n", timeBatchElapsed)
	showProcess.SetText(showProcessSen)
	log.Println("############################################################################################")
	printPairedStringArrayForPlan4(matchStrings41, matchStrings42, matchStrings43, matchStrings44)
}
