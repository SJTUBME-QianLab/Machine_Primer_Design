package main

func getOutOfPlan3() {

	for i := 0; i < len(input3S4); i++ {

		if numberTotalMatches3 > maxOutSeqNum {
			break
		}

		tempSrc1S := input3S1[i]
		tempSrc2A := input3A2[i]
		tempSrc3S := input3S3[i]
		tempSrc4S := input3S4[i]

		tempSrc1S2A := joinSeq(tempSrc1S, tempSrc2A)

		if screeningFourSeqs(tempSrc1S2A) {
			continue
		}

		flagCircleOut1, _ := screeningCircle8Seqs(tempSrc1S2A, lowCircleThreshPlan3, highCircleThreshPlan3)

		if flagCircleOut1 {
			continue
		}

		if pair5BetweenSeq(tempSrc1S, tempSrc1S2A) {
			continue
		}

		if len(tempSrc3S) >= 5 {
			if pair5BetweenSeq(tempSrc3S, tempSrc1S2A) {
				continue
			}
		}

		if len(tempSrc4S) >= 5 {
			if pair5BetweenSeq(tempSrc4S, tempSrc1S2A) {
				continue
			}
		}

		numberTotalMatches3++

		matchStrings31 = append(matchStrings31, tempSrc4S)
		matchStrings32 = append(matchStrings32, tempSrc3S)
		matchStrings33 = append(matchStrings33, tempSrc1S)
		matchStrings34 = append(matchStrings34, tempSrc1S2A)
	}
}
