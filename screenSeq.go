package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	ell                       = int(lengthUS - 2)
	cgLowBound                = int(lowGCUS - 2)
	cgHighBound               = int(highGCUS - 2)
	ell1                      = 0
	ell2                      = 0
	numberBatchs              = 100
	flagOut1                  = true
	flagOut2                  = true
	numberTotalMatches1 int64 = 0
	numberTotalMatches2 int64 = 0
	numberTotalMatches3 int64 = 0
	allSeqsLooked       int   = 0
	patternInts               = []int{0x00, 0x15, 0x2a, 0x3f}
	patternStrings            = []string{"AAA", "TTT", "GGG", "CCC"}
	repeatPatterns            = []int{0x11, 0x44, 0xee, 0xbb}
	repeatString              = []string{"ATAT", "TATA", "CGCG", "GCGC"}
	matchStrings11            = make([]string, 0, maxOutSeqNum)
	matchStrings12            = make([]string, 0, maxOutSeqNum)
	matchStrings21            = make([]string, 0, maxOutSeqNum)
	matchStrings22            = make([]string, 0, maxOutSeqNum)
	matchStrings31            = make([]string, 0, maxOutSeqNum)
	matchStrings32            = make([]string, 0, maxOutSeqNum)
	matchStrings33            = make([]string, 0, maxOutSeqNum)
	matchStrings34            = make([]string, 0, maxOutSeqNum)
	//wg                  sync.WaitGroup
)

type Pair struct {
	str   string
	cgCnt int
}

func pair5BetweenSeq(u1 string, u2 string) bool {

	u1Reverse := pairSequenceReverseDNA2DNA(u1)
	for i := 0; i < len(u1)-4; i++ {
		if strings.Contains(u2, u1Reverse[i:i+5]) {
			return true
		}
	}
	return false
}

func screeningPiar2Seqs(s1 string) bool {
	for i := 0; i < len(s1)-3; i++ {
		tempSrc := s1[i : i+2]
		s1PairReverse := pairSequenceReverseDNA2DNA(tempSrc)
		if strings.Contains(s1[i+2:len(s1)], s1PairReverse) {
			return true
		}
	}
	return false
}

func screeningCircle8Seqs(s1 string, lowGapThresh int64, highGapThresh int64) (bool, string) {
	for i := 0; i < len(s1)-7; i++ {
		tempSrc := reverseString(s1[i : i+2])
		for j := i + 3; j < len(s1)-2; j++ {
			if pairString(tempSrc, s1[j:j+2]) {
				if screeningPiar2Seqs(s1[i+1 : j+2]) {
					break
				} else {
					tempGap := int64(j - i - 2)
					if tempGap <= highGapThresh && tempGap >= lowGapThresh {
						return true, s1[i : j+2]
					}
				}
			}
		}
	}
	return false, s1[:]
}

func screeningSeq(s1 string) bool {
	s1PairReverse := pairSequenceReverse(s1[len(s1)-3 : len(s1)])
	if strings.Contains(s1[0:len(s1)-4], s1PairReverse) {
		return true
	}
	return false
}

func screeningMiddleSeq(s1 string) bool {
	for i := 1; i < len(s1)-7; i++ {
		tempSrc := s1[i : i+3]
		s1PairReverse := pairSequenceReverse(tempSrc)
		if strings.Contains(s1[i+4:len(s1)-1], s1PairReverse) {
			return true
		}
	}
	return false
}

func countStr(s1 string, s2 byte) int {
	count := 0
	for i := 0; i < len(s1); i++ {
		if s1[i] == s2 {
			count++
		}
	}
	return count
}

func countStrings(s1 string, s2 string) int {
	count := 0
	for i := 0; i < len(s1)-1; i++ {
		if s1[i:i+2] == s2[:] {
			count++
		}
	}
	return count
}

func pairString(s1 string, s2 string) bool {
	seqCouple := true
	for i := 0; i < len(s1); i++ {
		if !((s1[i] == purines[0] && s2[i] == purines[1]) || (s1[i] == purines[1] && s2[i] == purines[0]) ||
			(s1[i] == purines[2] && s2[i] == purines[3]) || (s1[i] == purines[3] && s2[i] == purines[2])) {
			seqCouple = false
		}
	}
	return seqCouple
}

// get piar string, like from CTGA get GACT
func pairSequenceReverse(s1 string) string {
	s2 := []rune(s1)
	for i := 0; i < len(s1); i++ {
		if s1[i] == purines[0] {
			s2[i] = 'T'
		} else if s1[i] == purines[1] {
			s2[i] = 'A'
		} else if s1[i] == purines[2] {
			s2[i] = 'C'
		} else if s1[i] == purines[3] {
			s2[i] = 'G'
		}
	}
	return reverseString(string(s2))
}

func joinSeq(s1 string, s2 string) string {
	bufferSrc := new(bytes.Buffer)
	bufferSrc.WriteString(s1)
	bufferSrc.WriteString(s2)
	return bufferSrc.String()
}

func screenCircle(strMergeds string) bool {
	// judgement-1234
	if pairString(strMergeds[0:3], reverseString(strMergeds[len(strMergeds)-3:len(strMergeds)])) {
		return true
	}

	// judgement-2
	if screeningSeq(strMergeds) {
		return true
	}

	// judgement-3
	if screeningSeq(reverseString(strMergeds)) {
		return true
	}

	// judgement-4
	if screeningMiddleSeq(strMergeds) {
		return true
	}
	return false
}

func checkSeries6(s1 string) bool {
	for i := 0; i < len(s1)-6; i++ {
		//if (s1[i] == purines[0] || s1[i] == purines[1]) && (s1[i+1] == purines[0] || s1[i+1] == purines[1]) && (s1[i+2] == purines[0] || s1[i+2] == purines[1]) && (s1[i+3] == purines[0] || s1[i+3] == purines[1]) && (s1[i+4] == purines[0] || s1[i+4] == purines[1]) && (s1[i+5] == purines[0] || s1[i+5] == purines[1]) && (s1[i+6] == purines[0] || s1[i+6] == purines[1]) {
		//	return true
		//}
		if (s1[i] == purines[2] || s1[i] == purines[3]) && (s1[i+1] == purines[2] || s1[i+1] == purines[3]) && (s1[i+2] == purines[2] || s1[i+2] == purines[3]) && (s1[i+3] == purines[2] || s1[i+3] == purines[3]) && (s1[i+4] == purines[2] || s1[i+4] == purines[3]) && (s1[i+5] == purines[2] || s1[i+5] == purines[3]) && (s1[i+6] == purines[2] || s1[i+6] == purines[3]) {
			return true
		}
	}
	return false
}

func checkSeries4GC(s1 string) bool {
	for i := 0; i < len(s1)-4; i++ {
		//if (s1[i] == purines[0] || s1[i] == purines[1]) && (s1[i+1] == purines[0] || s1[i+1] == purines[1]) && (s1[i+2] == purines[0] || s1[i+2] == purines[1]) && (s1[i+3] == purines[0] || s1[i+3] == purines[1]) && (s1[i+4] == purines[0] || s1[i+4] == purines[1]) && (s1[i+5] == purines[0] || s1[i+5] == purines[1]) && (s1[i+6] == purines[0] || s1[i+6] == purines[1]) {
		//	return true
		//}
		if (s1[i] == purines[2] || s1[i] == purines[3]) && (s1[i+1] == purines[2] || s1[i+1] == purines[3]) && (s1[i+2] == purines[2] || s1[i+2] == purines[3]) && (s1[i+3] == purines[2] || s1[i+3] == purines[3]) && (s1[i+4] == purines[2] || s1[i+4] == purines[3]) {
			return true
		}
	}
	return false
}

func checkSeries3GC(s1 string) bool {
	for i := 0; i < len(s1)-3; i++ {
		if (s1[i] == purines[2] || s1[i] == purines[3]) && (s1[i+1] == purines[2] || s1[i+1] == purines[3]) && (s1[i+2] == purines[2] || s1[i+2] == purines[3]) && (s1[i+3] == purines[2] || s1[i+3] == purines[3]) {
			return true
		}
	}
	return false
}

func checkSeries5GC(s1 string) bool {
	for i := 0; i < len(s1)-5; i++ {
		if (s1[i] == purines[2] || s1[i] == purines[3]) && (s1[i+1] == purines[2] || s1[i+1] == purines[3]) && (s1[i+2] == purines[2] || s1[i+2] == purines[3]) && (s1[i+3] == purines[2] || s1[i+3] == purines[3]) && (s1[i+4] == purines[2] || s1[i+4] == purines[3]) && (s1[i+5] == purines[2] || s1[i+5] == purines[3]) {
			return true
		}
	}
	return false
}

func checkSeries3AT(s1 string) bool {
	for i := 0; i < len(s1)-2; i++ {
		if (s1[i] == purines[0] || s1[i] == purines[1]) && (s1[i+1] == purines[0] || s1[i+1] == purines[1]) && (s1[i+2] == purines[0] || s1[i+2] == purines[1]) {
			return true
		}
	}
	return false
}

// subsets is responsible for  finding all valid strings of length 10 with the following constraints:
// 1) it contains no "AAA", "TTT", "GGG", "CCC" substrings
// 2) the number of character 'C' and 'G' should be above 2
// results are stored into a Pair slice, where each item stores a valid string and the number of 'C'+'G'
func subsets() []Pair {
	res := make([]Pair, 0)
	length := ell / 2
	N := 1 << uint(2*length)
	//matchStrings := make([]string, 0, N)
	//resultFileName := fmt.Sprintf("result.txt")
	//resultFile, err := os.OpenFile(resultFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	//if err != nil {
	//	log.Println("open result file failed!", err.Error())
	//	os.Exit(1)
	//}
	//defer resultFile.Close()

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

			// 'A' + 'T' check
			//if num == 0 {
			//	aCnt++
			//}else if num == 1{
			//	tCnt++
			//}

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
		if flag && !checkSeries3AT(string(str)) && !checkSeries4GC(string(str)) {
			//if flag && !checkSeries3AT(string(str)) && !checkSeries6(string(str)) {
			res = append(res, Pair{string(str), cgCnt})
			//matchStrings = append(matchStrings, string(str))
		}
	}
	log.Println("get src done!!")
	log.Println(len(res))

	//var b strings.Builder
	//for _, str := range matchStrings {
	//	b.WriteString(str)
	//	b.WriteString("\n")
	//}
	//resultFile.WriteString(b.String())

	return res
}

// merge is responsible for find all valid combinations
// it merge two strings of length 10 in to a string of length 20
// and find all valid ones fulfilling the following constraints:
// 1) it contains no "AAA", "TTT", "GGG", "CCC" substrings
// 2) it contains no "ATAT", "TATA", "GCGC", "CGCG" substrings
// 3) the number of character 'C' and 'G' should be in the range [11, 13]
// final results are written in to a given file
func merge(src []Pair, start, end int) {

	// concanated string
	var strMerged string
	var strMergeds string

	allSeqsNeeded := (end - start) * len(src) * 4 * len(input1S)
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

		for j := 0; j < len(src); j++ {
			// to indicate a valid string or not
			flag := true

			if numberTotalMatches1 > maxOutSeqNum && numberTotalMatches2 > maxOutSeqNum {
				goto Loop
			}

			// 'C'+'G' check
			cgTotal := src[i].cgCnt + src[j].cgCnt
			if cgTotal < cgLowBound || cgTotal > cgHighBound {
				continue
			}

			// "AAA", "TTT", "GGG", "CCC" check
			strMerged = src[i].str + src[j].str

			if countStr(strMerged, 'A') < 3 || countStr(strMerged, 'T') < 3 {
				continue
			}

			if countStrings(strMerged, "AA") > 0 || countStrings(strMerged, "TT") > 1 {
				continue
			}

			//if checkSeries3GC(strMerged[0 : 3]) || checkSeries3GC(strMerged[len(strMerged)-3 : ]) || checkSeries6(strMerged[(ell/2)-3 : (ell/2)+3]) {
			//	continue
			//}

			if checkSeries3GC(strMerged[0:5]) || checkSeries3GC(strMerged[len(strMerged)-5:]) {
				continue
			}

			if checkSeries3AT(strMerged[(ell/2)-2 : (ell/2)+2]) {
				continue
			}

			//if checkSeries6(strMerged[(ell/2)-6 : (ell/2)+6]) {
			//	continue
			//}

			if checkSeries4GC(strMerged[(ell/2)-6 : (ell/2)+6]) {
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

			if !flag && checkSeries6(strMerged) {
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

						allSeqsLooked = (i-start)*len(src)*4*len(input1S) + j*4*len(input1S) + k*2*len(input1S) + z*len(input1S) + n + 1

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

						if checkSeries4GC(out1Us1S[(len(out1Us1S)/2)-5:(len(out1Us1S)/2)+5]) || checkSeries5GC(out1Us2A[(len(out1Us2A)/2)-5:(len(out1Us2A)/2)+5]) {
							flagOut1 = false
						}

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

func printPairedStringArray(s1 []string, s2 []string, plan int) {
	for i := 0; i < len(s1); i++ {
		if plan == 1 {
			log.Printf("Plan-One US+1S: %v \n", s1[i])
			log.Printf("Plan-One US+2A: %v \n", s2[i])
			log.Println("-------------------------------------------------------------------------------------------------------------------------")
		} else {
			log.Printf("Plan-Two 1S: %v \n", s1[i])
			log.Printf("Plan-Two 1S+US+2A: %v \n", s2[i])
			log.Println("-------------------------------------------------------------------------------------------------------------------------")
		}
	}

	if len(s1) == 0 {
		log.Printf("Sorry, we don't find any targets for plan-%d! \n", plan)
	}
}

func printPairedStringArrayForPlan3(s1 []string, s2 []string, s3 []string, s4 []string) {
	for i := 0; i < len(s1); i++ {

		log.Printf("Plan-Three 4S: %v \n", s1[i])
		log.Printf("Plan-Three 3S: %v \n", s2[i])
		log.Printf("Plan-Three 1S: %v \n", s3[i])
		log.Printf("Plan-Three 1S+2A: %v \n", s4[i])
		log.Println("-------------------------------------------------------------------------------------------------------------------------")

	}

	if len(s1) == 0 {
		log.Printf("Sorry, we don't find any targets for plan-3. \n")
	}
}

func printPairedStringArrayForPlan4(s1 []string, s2 []string, s3 []string, s4 []string) {
	for i := 0; i < len(s1); i++ {
		log.Printf("Plan-Four 3S+2A+1S: %v \n", s1[i])
		log.Printf("Plan-Four 1S+4A+2A: %v \n", s2[i])
		log.Printf("Plan-Four 4A+2A+3S: %v \n", s3[i])
		log.Printf("Plan-Four 3S+1S+4A: %v \n", s4[i])
		log.Println("-------------------------------------------------------------------------------------------------------------------------")
	}

	if len(s1) == 0 {
		log.Printf("Sorry, we don't find any targets for plan-4. \n")
	}
}
