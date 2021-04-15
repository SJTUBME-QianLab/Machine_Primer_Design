package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// subsets is responsible for  finding all valid strings of length 10 with the following constraints:
// 1) it contains no "AAA", "TTT", "GGG", "CCC" substrings
// 2) the number of character 'C' and 'G' should be above 2
// results are stored into a Pair slice, where each item stores a valid string and the number of 'C'+'G'
func subsetsOdd(ells int) []Pair {
	res := make([]Pair, 0)
	length := ells / 2
	N := 1 << uint(2*length)

	for i := 0; i < N; i++ {
		// to store encoded results of i
		str := make([]byte, length)
		// count the numbers of 'C' and 'G'
		cgCnt := 0
		// match = 0B11, used to extract the character ('A', 'T', 'G', 'C'), for counting purpose
		matchChar := 0x3
		// matchPattern = 0B11111111, used to extract consequent binary of length 8, e.g., "0101010101",
		// for further patterns match (0x00, 0x55, 0xaa, 0xff)
		matchPattern := 0x3f
		repeatMatchPattern := 0xff
		// to indicate whether i is valid or not
		flag := true

		for j := 0; j < length; j++ {
			//encode i into a quaternary number represented by 'A', 'T', 'G', and 'C' (0, 1, 2, 3 respectively)
			num := (matchChar & i) >> uint(2*j)
			str[length-1-j] = purines[num]

			// 'C'+'G' check
			if num == 2 || num == 3 {
				cgCnt++
			}
			if (j + 1 - cgCnt) > (ell - cgLowBound) {
				flag = false
				break
			}

			// "AAA", "TTT", "GGG", "CCC" check
			if j >= 2 {
				for _, patternInt := range patternInts {
					// extract 8-bit number and use XOR to check whether it's a pattern
					if ((i&matchPattern)>>uint((2*j-4)))^patternInt == 0 {
						flag = false
						break
					}
				}
				if !flag {
					break
				}

				// shift left to match the next character
				matchPattern <<= 2
			}

			// AT/TA/GC/CG check
			if j >= 3 {
				for _, pattern := range repeatPatterns {
					if ((i&repeatMatchPattern)>>uint((2*j-6)))^pattern == 0 {
						flag = false
						break
					}
				}
				if !flag {
					break
				}
				repeatMatchPattern <<= 2
			}

			// shift left to match the next pattern
			matchChar <<= 2
		}

		// find one valid string, store it into the result slice
		if flag {
			res = append(res, Pair{string(str), cgCnt})
		}
	}
	log.Println("get src done!!")
	return res
}

// merge is responsible for find all valid combinations
// it merge two strings of length 10 in to a string of length 20
// and find all valid ones fulfilling the following constraints:
// 1) it contains no "AAA", "TTT", "GGG", "CCC" substrings
// 2) it contains no "ATAT", "TATA", "GCGC", "CGCG" substrings
// 3) the number of character 'C' and 'G' should be in the range [11, 13]
// final results are written in to a given file
func mergeOdd(src1 []Pair, src2 []Pair, start, end int) {

	// concanated string
	var strMerged string
	var strMergeds string

	allSeqsNeeded := (end - start) * len(src2) * 4 * len(input1S)
	timeBatchBegin := time.Now()

	if numberTotalMatches1 > maxOutSeqNum && numberTotalMatches2 > maxOutSeqNum {
		goto Loop
	}
	for i := start; i < end; i++ {
		// store all valid strings
		// apply for a maximum possible storage to avoid reallocation during growth

		if numberTotalMatches1 > maxOutSeqNum && numberTotalMatches2 > maxOutSeqNum {
			goto Loop
		}

		for j := 0; j < len(src2); j++ {
			// to indicate a valid string or not
			flag := true

			if numberTotalMatches1 > maxOutSeqNum && numberTotalMatches2 > maxOutSeqNum {
				goto Loop
			}

			// 'C'+'G' check
			cgTotal := src1[i].cgCnt + src2[j].cgCnt
			if cgTotal < cgLowBound || cgTotal > cgHighBound {
				continue
			}

			// "AAA", "TTT", "GGG", "CCC" check
			strMerged = src1[i].str + src2[j].str

			if checkSeries6(strMerged[(ell/2)-6 : (ell/2)+6]) {
				continue
			}

			for _, substr := range patternStrings {
				// check only whether pattern strings exist in the inner substring of length 6
				if strings.Contains(strMerged[(ell/2)-2:(ell/2)+2], substr) {
					flag = false
					break
				}
			}
			for _, substr := range repeatString {
				// check only whether pattern strings exist in the inner substring of length 6
				if strings.Contains(strMerged[(ell/2)-3:(ell/2)+3], substr) {
					flag = false
					break
				}
			}

			if !flag {
				continue
			}

			for k := 0; k < 2; k++ {

				if numberTotalMatches1 > maxOutSeqNum && numberTotalMatches2 > maxOutSeqNum {
					goto Loop
				}

				for z := 0; z < 2; z++ {

					if numberTotalMatches1 > maxOutSeqNum && numberTotalMatches2 > maxOutSeqNum {
						goto Loop
					}

					bufferSrc := new(bytes.Buffer)
					bufferSrc.WriteString(string(purines[2+k]))
					bufferSrc.WriteString(strMerged)
					bufferSrc.WriteString(string(purines[2+z]))
					strMergeds = bufferSrc.String()

					for _, substr := range patternStrings {
						// check only whether pattern strings exist in the inner substring of length 6
						if strings.Contains(strMergeds[0:4], substr) || strings.Contains(strMergeds[len(strMergeds)-3:len(strMergeds)], substr) {
							flag = false
							break
						}
					}

					if !flag {
						continue
					}

					for _, substr := range repeatString {
						// check only whether pattern strings exist in the inner substring of length 6
						if strings.Contains(strMergeds[0:5], substr) || strings.Contains(strMergeds[len(strMergeds)-4:len(strMergeds)], substr) {
							flag = false
							break
						}
					}

					if !flag {
						continue
					}

					// judgement-1234
					if screenCircle(strMergeds) {
						continue
					}

					for n := 0; n < len(input1S); n++ {

						allSeqsLooked = (i-start)*len(src2)*4*len(input1S) + j*4*len(input1S) + k*2*len(input1S) + z*len(input1S) + n + 1

						if allSeqsLooked%1 == 0 {
							timeBatchElapsed := time.Since(timeBatchBegin).Seconds()
							valuePercent := float64(allSeqsLooked) / float64(allSeqsNeeded)
							timeLeft := (timeBatchElapsed / valuePercent) / 3600
							valuePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.6f", valuePercent*100), 64)
							showProcessSen := fmt.Sprintf("Checked %d in %d allSeqs %v%% in %.1f s, find %d plan-1, %d plan-2 and %d plan-3 targets, time left: %.1f h \n", allSeqsLooked, allSeqsNeeded, valuePercent, timeBatchElapsed, numberTotalMatches1, numberTotalMatches2, numberTotalMatches3, timeLeft)
							showProcess.SetText(showProcessSen)
							//proBar.SetValue(int(valuePercent*100))
							// log.Println("checked!!!")
						}

						flagOut1 = true
						flagOut2 = true

						if numberTotalMatches1 > maxOutSeqNum {
							flagOut1 = false
						}
						if numberTotalMatches2 > maxOutSeqNum {
							flagOut2 = false
						}
						if numberTotalMatches1 > maxOutSeqNum && numberTotalMatches2 > maxOutSeqNum {
							goto Loop
						}

						out1Us1S := joinSeq(strMergeds, input1S[n])
						out1Us2A := joinSeq(strMergeds, input2A[n])
						out21S := input1S[n]
						out21SUs2A := joinSeq(input1S[n], joinSeq(strMergeds, input2A[n]))

						if flagOut1 {
							flagpair5BetweenSeq := pair5BetweenSeq(out1Us1S, out1Us2A)
							if flagpair5BetweenSeq {
								flagOut1 = false
							}
						}

						if flagOut1 && (screeningFourSeqs(out1Us1S) || screeningFourSeqs(out1Us2A)) {
							flagOut1 = false
						}

						if flagOut1 {
							flagCircleOut1, _ := screeningCircle8Seqs(out1Us1S, lowCircleThresh, highCircleThresh)
							if flagCircleOut1 {
								flagOut1 = false
							}
						}

						if flagOut1 {
							flagCircleOut1, _ := screeningCircle8Seqs(out1Us2A, lowCircleThresh, highCircleThresh)
							if flagCircleOut1 {
								flagOut1 = false
							}
						}

						if flagOut2 {
							flagpair5BetweenSeq := pair5BetweenSeq(out21S, out21SUs2A)
							if flagpair5BetweenSeq {
								flagOut2 = false
							}
						}

						if flagOut2 && screeningFourSeqs(out21SUs2A) {
							flagOut2 = false
						}

						if flagOut2 {
							flagCircleOut2, _ := screeningCircle8Seqs(out21SUs2A, lowCircleThresh, highCircleThresh)
							if flagCircleOut2 {
								flagOut2 = false
							}
						}

						if flagOut1 {
							numberTotalMatches1++
							matchStrings11 = append(matchStrings11, out1Us1S)
							matchStrings12 = append(matchStrings12, out1Us2A)
						}
						if flagOut2 {
							numberTotalMatches2++
							matchStrings21 = append(matchStrings21, out21S)
							matchStrings22 = append(matchStrings22, out21SUs2A)
						}

					}
				}
			}
		}
	}
Loop:
	timeBatchElapsed := time.Since(timeBatchBegin).Seconds()
	showProcessSen := fmt.Sprintf("Done in %.1f s \n", timeBatchElapsed)
	showProcess.SetText(showProcessSen)
}

// This program is used to find all strings of length 20 consisting of 'A', 'T', 'G', 'C'
// each valid one should fulfill the following requirements:
// 1) the number of 'C' and 'G' is no less than 12 and no more than 14
// 2) it contains no substrings like "AAAA", "TTTT", "GGGG", and "CCCC"

// Algorithm: divide and conquer
// 1. find all valid strings of length 10 (using encoding techniques)
// 2. find all valid strings of length 20 by combinations of strings of length 10
// 3. multiple goroutines are launched in this program to take the advantage of  multi-core processor
// 4. results are written into multiple files (199114775296 pairs, about 4TB)
func makeUSOdd() {

	if ell%2 != 0 {
		ell1 = ell / 2
		ell2 = ell - ell/2

		src1 := subsetsOdd(ell1)
		src2 := subsetsOdd(ell2)

		batchSize := len(src1) / numberBatchs

		i := int(batchUS)

		batchStart := i * batchSize
		batchEnd := (i + 1) * batchSize
		if i == numberBatchs {
			batchEnd = len(src1)
		}

		log.Printf("\nbatch %d begins, batch start: %d, batch end: %d\n", i, batchStart, batchEnd)

		timeBatchBegin := time.Now()
		mergeOdd(src1, src2, batchStart, batchEnd)
		timeBatchElapsed := time.Since(timeBatchBegin)

		log.Printf("batch %d finished.\n", i)
		log.Printf("Time used for this batch: %v, found matches: %d and %d\n", timeBatchElapsed, numberTotalMatches1, numberTotalMatches2)

		// log.Printf("\n%d batchs processed, program finished!\n", numberBatchs)
		log.Println("############################################################################################")
		printPairedStringArray(matchStrings11, matchStrings12, 1)
		log.Println("############################################################################################")
		printPairedStringArray(matchStrings21, matchStrings22, 2)
		log.Println("############################################################################################")
		printPairedStringArrayForPlan3(matchStrings31, matchStrings32, matchStrings33, matchStrings34)
	} else {
		src := subsets()
		batchSize := len(src) / numberBatchs

		i := int(batchUS)

		batchStart := i * batchSize
		batchEnd := (i + 1) * batchSize
		if i == numberBatchs {
			batchEnd = len(src)
		}

		log.Printf("\nbatch %d begins, batch start: %d, batch end: %d\n", i, batchStart, batchEnd)

		timeBatchBegin := time.Now()
		merge(src, batchStart, batchEnd)
		timeBatchElapsed := time.Since(timeBatchBegin)

		log.Printf("batch %d finished.\n", i)
		log.Printf("Time used for this batch: %v, found matches: %d and %d\n", timeBatchElapsed, numberTotalMatches1, numberTotalMatches2)

		// log.Printf("\n%d batchs processed, program finished!\n", numberBatchs)
		log.Println("############################################################################################")
		printPairedStringArray(matchStrings11, matchStrings12, 1)
		log.Println("############################################################################################")
		printPairedStringArray(matchStrings21, matchStrings22, 2)
		log.Println("############################################################################################")
		printPairedStringArrayForPlan3(matchStrings31, matchStrings32, matchStrings33, matchStrings34)
		log.Println("############################################################################################")
	}
}
