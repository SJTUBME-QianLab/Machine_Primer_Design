package main

var (
	lowNum1   int
	highNum1  int
	GCLowNum1 int
	lowNum2   int
	highNum2  int
	GCLowNum2 int
	input4S1  []string
	input4A2  []string
	input4S3  []string
	input4A4  []string
)

func analysisInputPlan4(inputSrc1 string, inputSrc2 string) {
	_, GCPercent1 := calGCPercent(inputSrc1)
	_, GCPercent2 := calGCPercent(inputSrc2)

	if GCPercent1 < threshGC {
		lowNum1 = 20
		highNum1 = 25
		GCLowNum1 = 9
	} else {
		lowNum1 = 18
		highNum1 = 22
		GCLowNum1 = 12
	}

	if GCPercent2 < threshGC {
		lowNum2 = 20
		highNum2 = 25
		GCLowNum2 = 9
	} else {
		lowNum2 = 18
		highNum2 = 22
		GCLowNum2 = 12
	}

	for i := 0; i < len(inputSrc1); i++ {
		for j := lowNum1 + i; j < highNum1+1+i; j++ {
			if j > len(inputSrc1) {
				break
			}
			tempSrc1 := inputSrc1[i:j]

			// judge GC content
			numGC, _ := calGCPercent(tempSrc1)
			if numGC < GCLowNum1 {
				continue
			}

			// delete the sequence whose self-pair is greater than or equal to 4
			if screeningFourSeqs(tempSrc1) {
				if screeningFourSeqs(tempSrc1[0 : len(tempSrc1)-1]) {
					break
				}
				continue
			}

			for gap := 0; gap <= 10; gap++ {
				for k := lowNum1 + j + gap; k < highNum1+gap+j+1; k++ {
					if k > len(inputSrc1) {
						break
					}
					tempSrc2 := pairSequenceDNA2DNA(inputSrc1[j+gap : k])

					numGC, _ := calGCPercent(tempSrc2)
					if numGC < GCLowNum1 {
						continue
					}

					if screeningFourSeqs(tempSrc2) {
						if screeningFourSeqs(tempSrc2[0 : len(tempSrc2)-1]) {
							break
						}
						continue
					}

					if pair5BetweenSeq(tempSrc1, tempSrc2) {
						continue
					}

					input4S1 = append(input4S1, tempSrc1)
					input4A2 = append(input4A2, tempSrc2)
				}
			}
		}
	}

	for i := 0; i < len(inputSrc2); i++ {
		for j := lowNum2 + i; j < highNum2+1+i; j++ {
			if j > len(inputSrc2) {
				break
			}
			tempSrc1 := inputSrc2[i:j]

			numGC, _ := calGCPercent(tempSrc1)
			if numGC < GCLowNum2 {
				continue
			}

			if screeningFourSeqs(tempSrc1) {
				if screeningFourSeqs(tempSrc1[0 : len(tempSrc1)-1]) {
					break
				}
				continue
			}

			for gap := 0; gap <= 10; gap++ {
				for k := lowNum2 + j + gap; k < highNum2+gap+j+1; k++ {
					if k > len(inputSrc2) {
						break
					}
					tempSrc2 := pairSequenceDNA2DNA(inputSrc2[j+gap : k])

					if screeningFourSeqs(tempSrc2) {
						if screeningFourSeqs(tempSrc2[0 : len(tempSrc2)-1]) {
							break
						}
						continue
					}

					if pair5BetweenSeq(tempSrc1, tempSrc2) {
						continue
					}

					input4S3 = append(input4S3, tempSrc1)
					input4A4 = append(input4A4, tempSrc2)
				}
			}
		}
	}
}
