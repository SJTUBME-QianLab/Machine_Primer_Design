package main

import (
	"strings"
)

const purines = "ATGC"

var (
	lowNum   int
	highNum  int
	GCLowNum int
	input1S  []string
	input2A  []string
	input3S1 []string
	input3A2 []string
	input3S3 []string
	input3S4 []string
)

func transa2A(s1 string) string {
	s2 := []rune(s1)

	for i := 0; i < len(s1); i++ {
		if s1[i] == 'a' {
			s2[i] = 'A'
		} else if s1[i] == 't' {
			s2[i] = 'T'
		} else if s1[i] == 'g' {
			s2[i] = 'G'
		} else if s1[i] == 'c' {
			s2[i] = 'C'
		}
	}
	return string(s2)
}

func calGCPercent(s1 string) (int, float64) {
	num := 0
	lenS1 := len(s1)
	for i := 0; i < len(s1); i++ {
		if s1[i] == 'G' || s1[i] == 'C' {
			num++
		}
	}
	numPercent := float64(num) / float64(lenS1)
	return num, numPercent
}

func reverseString(s1 string) string {
	runes := []rune(s1)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func pairSequenceReverseDNA2DNA(s1 string) string {
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

func pair4BetweenSeq(u1 string, u2 string) bool {

	u1Reverse := pairSequenceReverseDNA2DNA(u1)
	for i := 0; i < len(u1)-3; i++ {
		if strings.Contains(u2, u1Reverse[i:i+4]) {
			return true
		}
	}
	return false
}

func screenInputSrc(tempSrc1 string) bool {

	// judge GC content
	numGC, _ := calGCPercent(tempSrc1)
	if numGC < GCLowNum {
		return true
	}

	// the GC content of the last six digits
	numGCLast3, _ := calGCPercent(tempSrc1[len(tempSrc1)-6:])
	if numGCLast3 < 3 {
		return true
	}

	// the last base must be G or C
	if tempSrc1[len(tempSrc1)-1] == 'G' || tempSrc1[len(tempSrc1)-1] == 'C' {
		return true
	}

	return false
}

func pairSequenceDNA2DNA(s1 string) string {
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
	return string(s2)
}

func screeningFourSeqs(s1 string) bool {
	for i := 0; i < len(s1)-7; i++ {
		tempSrc := s1[i : i+4]
		s1PairReverse := pairSequenceReverseDNA2DNA(tempSrc)
		if strings.Contains(s1[i+4:], s1PairReverse) {
			return true
		}
	}
	return false
}

func screeningContinueSeqs(s1 string, num int) bool {
	for i := 0; i < len(s1)-(num*2-1); i++ {
		tempSrc := s1[i : i+num]
		s1PairReverse := pairSequenceReverseDNA2DNA(tempSrc)
		if strings.Contains(s1[i+num:], s1PairReverse) {
			return true
		}
	}
	return false
}

func analysisInput(inputSrc1 string) {
	_, GCPercent1 := calGCPercent(inputSrc1)

	if GCPercent1 < threshGC {
		lowNum = 20
		highNum = 25
		GCLowNum = 9
	} else {
		lowNum = 18
		highNum = 22
		GCLowNum = 12
	}

	for i := 0; i < len(inputSrc1); i++ {
		for j := lowNum + i; j < highNum+1+i; j++ {

			if j > len(inputSrc1) {
				break
			}

			tempSrc1 := inputSrc1[i:j]

			if screenInputSrc(tempSrc1) {
				continue
			}

			// delete the sequence whose self-pair is greater than or equal to 4
			if screeningFourSeqs(tempSrc1) {
				if screeningFourSeqs(tempSrc1[0 : len(tempSrc1)-1]) {
					break
				}
				continue
			}

			for gap := 0; gap <= 50; gap++ {
				for k := lowNum + j + gap; k < highNum+gap+j+1; k++ {
					if k > len(inputSrc1) {
						break
					}
					tempSrc2 := pairSequenceDNA2DNA(inputSrc1[j+gap : k])

					if screenInputSrc(tempSrc2) {
						continue
					}

					// delete the sequence whose self-pair is greater than or equal to 4
					if screeningFourSeqs(tempSrc2) {
						if screeningFourSeqs(tempSrc2[0 : len(tempSrc2)-1]) {
							break
						}
						continue
					}

					if pair4BetweenSeq(tempSrc1, tempSrc2) {
						continue
					}

					input1S = append(input1S, tempSrc1)
					input2A = append(input2A, tempSrc2)

					if i < lowNum {
						continue
					}

					if gap < lowNum {
						continue
					}

					for lenS4 := lowNum; lenS4 <= highNum; lenS4++ {
						for bGap := 0; bGap <= i-lowNum; bGap++ {

							//if i-lenS4-bGap < 0 || bGap > 50 {
							//	break
							//}

							if i-lenS4-bGap < 0 {
								break
							}

							tempSrc4S := inputSrc1[i-lenS4-bGap : i-bGap]

							if screeningFourSeqs(tempSrc4S) {
								continue
							}

							if pair4BetweenSeq(tempSrc1, tempSrc4S) {
								continue
							}

							if pair4BetweenSeq(tempSrc2, tempSrc4S) {
								continue
							}

							for lenS3 := lowNum; lenS3 <= highNum; lenS3++ {
								for fGap := 0; fGap <= i-lowNum; fGap++ {

									if j+lenS3+fGap > j+gap {
										break
									}

									tempSrc3S := inputSrc1[j+fGap : j+lenS3+fGap]

									if screeningFourSeqs(tempSrc3S) {
										continue
									}

									if pair4BetweenSeq(tempSrc1, tempSrc3S) {
										continue
									}

									if pair4BetweenSeq(tempSrc2, tempSrc3S) {
										continue
									}

									if pair4BetweenSeq(tempSrc4S, tempSrc3S) {
										continue
									}

									input3S1 = append(input3S1, tempSrc1)
									input3A2 = append(input3A2, tempSrc2)
									input3S3 = append(input3S3, tempSrc3S)
									input3S4 = append(input3S4, tempSrc4S)
								}
							}
						}
					}
				}
			}
		}
	}
}
