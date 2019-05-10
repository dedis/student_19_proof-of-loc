package latencyprotocol

import (
	"errors"
	"log"

	"strconv"
	"time"

	sigAlg "golang.org/x/crypto/ed25519"
)

const delta = time.Duration(1000 * time.Millisecond)

func (A *Block) to(B *Block) (time.Duration, bool) {
	aToB, aToBKnown := A.getLatency(B)
	bToA, bToAKnown := B.getLatency(A)

	if aToBKnown && bToAKnown {
		return (aToB + bToA) / 2, true
	}

	if aToBKnown {
		return aToB, true
	}

	if bToAKnown {
		return bToA, true
	}

	return time.Duration(0), false

}

/*ApproximateDistance is a function to approximate the distance between two given nodes, e.g.,
node A wants to approximate the distance between nodes B and C. Node A relies on the information
in the blockchain about distances to B, C, between B and C and its own estimations to B and C,
applies triangularization and computes an estimate of the distance. */
func (A *Block) ApproximateDistance(B *Block, C *Block, delta time.Duration) (time.Duration, bool, error) {

	aToB, aToBKnown := A.getLatency(B)
	bToA, bToAKnown := B.getLatency(A)

	aToC, aToCKnown := A.getLatency(C)
	cToA, cToAKnown := C.getLatency(A)

	bToC, bToCKnown := B.getLatency(C)
	cToB, cToBKnown := C.getLatency(B)

	//the nodes know each other
	if cToBKnown && bToCKnown {

		//they say different things
		if timesContradictory(bToC, cToB, delta) {
			return time.Duration(0), false, errors.New("Distances contradictory: " + strconv.Itoa(int(time.Duration(bToC-cToB))))
		}

		lat := time.Duration((cToB + bToC) / 2)

		if aToCKnown && aToBKnown {
			if aToC+aToB < lat {
				return time.Duration(0), false, errors.New("Distances contradictory: " + strconv.Itoa(int(time.Duration(bToC-cToB))))
			}
		}
		return lat, true, nil
	}

	if aToBKnown && bToAKnown {
		if timesContradictory(aToB, bToA, delta) {
			return time.Duration(0), false, errors.New("Distances contradictory: " + strconv.Itoa(int(time.Duration(aToB-bToA))))
		}

		avgAB := (aToB + bToA) / 2

		if aToCKnown && cToAKnown {
			if timesContradictory(aToC, cToA, delta) {
				return time.Duration(0), false, errors.New("Distances contradictory: " + strconv.Itoa(int(time.Duration(aToC-cToA))))
			}
			avgAC := (cToA + aToC) / 2

			return Pythagoras(avgAB, avgAC), true, nil

		}

		if aToCKnown && !cToAKnown {
			return Pythagoras(avgAB, aToC), true, nil

		}

		if !aToCKnown && cToAKnown {
			return Pythagoras(avgAB, cToA), true, nil

		}

	}

	if bToAKnown && !aToBKnown {
		if aToCKnown && cToAKnown {
			if timesContradictory(aToC, cToA, delta) {
				return time.Duration(0), false, errors.New("Distances contradictory: " + strconv.Itoa(int(time.Duration(aToC-cToA))))
			}

			avgAC := (cToA + aToC) / 2

			return Pythagoras(bToA, avgAC), true, nil

		}

		if aToCKnown && !cToAKnown {
			return Pythagoras(bToA, aToC), true, nil

		}

		if !aToCKnown && cToAKnown {
			return Pythagoras(bToA, cToA), true, nil

		}

	}

	if !bToAKnown && aToBKnown {

		if aToCKnown && cToAKnown {
			if timesContradictory(aToC, cToA, delta) {
				return time.Duration(0), false, errors.New("Distances contradictory: " + strconv.Itoa(int(time.Duration(aToC-cToA))))
			}

			avgAC := (cToA + aToC) / 2

			return Pythagoras(aToB, avgAC), true, nil

		}

		if aToCKnown && !cToAKnown {
			return Pythagoras(aToB, aToC), true, nil

		}

		if !aToCKnown && cToAKnown {
			return Pythagoras(aToB, cToA), true, nil

		}

	}

	return time.Duration(0), false, errors.New("Not enough information")

}

//Pythagoras estimates the distance between two points with known distances to a common third point b using the Pythagorean theorem
//Since the angle between the three points is between 0 and 180 degrees, the function assumes an average angle of 90 degreess
func Pythagoras(p1 time.Duration, p2 time.Duration) time.Duration {
	return ((p1 ^ 2) + (p2 ^ 2)) ^ (1 / 2)
}

func (A *Block) getLatency(B *Block) (time.Duration, bool) {

	key := string(B.ID.PublicKey)
	latencyStruct, isPresent := A.Latencies[key]
	if !isPresent {
		return 0, false
	}
	return latencyStruct.Latency, true
}

func timesContradictory(time1 time.Duration, time2 time.Duration, delta time.Duration) bool {
	return time.Duration(time1-time2) > delta || time.Duration(time2-time1) > delta
}

//ApproximateOverChain approximates a distance between two nodes over a chain
func (chain *Chain) ApproximateOverChain(B *Node, C *Node) (time.Duration, error) {

	collectedDistances := make([]time.Duration, 0)

	blocks := chain.Blocks

	var latestBlockB *Block
	var latestBlockC *Block

	bFound := false
	cFound := false

	for i := len(blocks) - 1; i >= 0 && !bFound && !cFound; i-- {
		if blocks[i].ID == B.ID && !bFound {
			latestBlockB = blocks[i]
			bFound = true
		}
		if blocks[i].ID == C.ID && !cFound {
			latestBlockC = blocks[i]
			cFound = true
		}
	}

	if !bFound && !cFound {
		return time.Duration(100000), errors.New("Nodes not part of chain")
	}

	for _, block := range blocks {
		if block.ID != B.ID && block.ID != C.ID {
			distance, isValid, err := block.ApproximateDistance(latestBlockB, latestBlockC, delta)
			if err != nil {
				return time.Duration(0), err
			}
			if isValid {
				collectedDistances = append(collectedDistances, distance)
			}
		}
	}

	if len(collectedDistances) == 0 {
		return time.Duration(0), errors.New("No information available")
	}

	//TODO compare distances among each other

	averageDistance := time.Duration(0)
	for _, dist := range collectedDistances {
		averageDistance += dist
	}

	return averageDistance / time.Duration(len(collectedDistances)), nil

}

type nodeTuple struct {
	A *sigAlg.PublicKey
	B *sigAlg.PublicKey
}

type triangle struct {
	A string
	B string
	C string
}

//CreateBlacklist iterates through a chain and for each block checks if the latencies qiven by its node make sense
// If they do not, the node is added to a blacklist of nodes not to be trusted
//Warning: for now, the following assumptions are made:
//* all nodes give latencies to all other nodes (except themselves)
//* all latencies are symmetric (A -> B == B -> A)
func CreateBlacklist(chain *Chain, delta time.Duration, verbose bool, threshGiven bool, threshold int) (Blacklistset, error) {

	N := len(chain.Blocks)

	if !threshGiven {
		threshold = UpperThreshold(N)
	}

	if verbose {
		log.Print("Threshold: " + strconv.Itoa(threshold))
	}
	blockMapper := make(map[string]*Block)

	blacklist := NewBlacklistset()

	for _, block := range chain.Blocks {
		blockMapper[string(block.ID.PublicKey)] = block
	}

	//for each node B
	//for each node C
	//for each node D
	/* Check B -> C, B -> D, C -> D
	* if triangle of lengths does not result in realistic angles (rule of 3 for triangles),
	B, C or D needs to be blacklisted -> add (B,C, D) to a "suspicious" list and keep checking B
	*/

	for _, BBlock := range chain.Blocks {

		Bstring := string(BBlock.ID.PublicKey)

		for Cstring := range BBlock.Latencies {
			if Cstring != Bstring {
				CBlock, CHere := blockMapper[Cstring]

				if CHere {

					for Dstring := range BBlock.Latencies {
						if Dstring != Cstring && Dstring != Bstring {
							DBlock, DHere := blockMapper[Dstring]

							if DHere {

								BtoD, BtoDHere := BBlock.getLatency(DBlock)
								BtoC, BtoCHere := BBlock.getLatency(CBlock)
								CtoD, CtoDHere := CBlock.getLatency(DBlock)

								if BtoDHere && BtoCHere && CtoDHere && !TriangleInequalitySatisfiedInt(int(BtoD), int(BtoC), int(CtoD)) {

									blacklist.Add(sigAlg.PublicKey([]byte(Bstring)))
									blacklist.Add(sigAlg.PublicKey([]byte(Cstring)))
									blacklist.Add(sigAlg.PublicKey([]byte(Dstring)))
								}

							}
						}

					}
				}
			}
		}
	}

	if verbose {
		log.Print("Before Thresholding: ")
		log.Print(blacklist.ToStringFake())
	}

	threshBlacklist := blacklist.GetBlacklistWithThreshold(threshold)

	if verbose {
		log.Print("After Thresholding: ")
		log.Print(threshBlacklist.ToStringFake())
	}

	return threshBlacklist, nil

}

//TriangleInequalitySatisfied returns whether the triangle inequality theorem is satisfied
func TriangleInequalitySatisfied(latAB time.Duration, latBC time.Duration, latCA time.Duration) bool {
	return latAB+latBC >= latCA && latAB+latCA >= latBC && latBC+latCA >= latAB
}

//TriangleInequalitySatisfiedInt returns whether the triangle inequality theorem is satisfied
func TriangleInequalitySatisfiedInt(latAB int, latBC int, latCA int) bool {
	return latAB+latBC >= latCA && latAB+latCA >= latBC && latBC+latCA >= latAB
}

func acceptableDifference(time1 time.Duration, time2 time.Duration, delta time.Duration) bool {
	return time1-time2 < delta && time2-time1 < delta
}

func triangleAlreadyEvaluated(A string, B string, C string, triangles []triangle) bool {
	for _, triangle := range triangles {
		angles := []string{triangle.A, triangle.B, triangle.C}
		if listsEquivalent(angles, []string{A, B, C}) {
			return true
		}

	}

	return false
}

func listsEquivalent(a, b []string) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !containsString(b, a[i]) {
			return false
		}
	}

	return true
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
