package latencyprotocol

import (
	"go.dedis.ch/kyber/v3/pairing"
	"go.dedis.ch/onet/v3"
	sigAlg "golang.org/x/crypto/ed25519"
	"math/rand"
	"time"
)

type sourceType int

const (
	random sourceType = iota
	accurate
	inaccurate
	variant
)

var distanceSuite = pairing.NewSuiteBn256()

/*Test ApproximateDistance initially by assuming
all nodes are honest, each node adds in the blockchain x distances from itself to other x nodes,
where these x nodes are randomly chosen. You can assume for now that there’s a publicly known source
of randomness that nodes use. Check the results by varying the number x and the total number of nodes N.*/

func initChain(N int, x int, src sourceType, nbLiars int, nbVictims int) (*Chain, []*NodeID) {
	local := onet.NewTCPTest(distanceSuite)
	local.Check = onet.CheckNone
	_, el, _ := local.GenTree(N, false)
	defer local.CloseAll()

	nodeIDs := make([]*NodeID, N)
	privateKeys := make([]sigAlg.PrivateKey, N)

	for h := 0; h < N; h++ {
		pub, priv, err := sigAlg.GenerateKey(nil)
		if err != nil {
			return nil, nil
		}
		nodeIDs[h] = &NodeID{el.List[h], pub}
		privateKeys[h] = priv
	}

	chain := Chain{[]*Block{}, []byte("testBucket")}

	for i := 0; i < N; i++ {
		latencies := make(map[string]ConfirmedLatency)
		nbLatencies := 0
		for j := 0; j < N && nbLatencies < x; j++ {
			if i != j {
				nbLatencies++
				switch src {
				case random:
					latencies[string(nodeIDs[j].PublicKey)] = ConfirmedLatency{
						time.Duration((rand.Intn(300-20) + 20)),
						nil,
						time.Now(),
						nil,
					}
				case accurate:
					latencies[string(nodeIDs[j].PublicKey)] =
						ConfirmedLatency{
							time.Duration(10 * (i + j + 1)),
							nil,
							time.Now(),
							nil,
						}
				case variant:
					//adapt to percentage of distance
					latencies[string(nodeIDs[j].PublicKey)] = ConfirmedLatency{
						time.Duration(10 * (i + j + 1 + rand.Intn(5))),
						nil,
						time.Now(),
						nil,
					}
				case inaccurate:
					if i < nbLiars && (N-nbVictims) <= j {
						latencies[string(nodeIDs[j].PublicKey)] = ConfirmedLatency{
							time.Duration(((i * 10000) + j + 1)),
							nil,
							time.Now(),
							nil,
						}
					} else {
						if j < nbLiars && (N-nbVictims) <= i {
							latencies[string(nodeIDs[j].PublicKey)] = ConfirmedLatency{
								time.Duration(((j * 10000) + i + 1)),
								nil,
								time.Now(),
								nil,
							}
						} else {
							latencies[string(nodeIDs[j].PublicKey)] =
								ConfirmedLatency{
									time.Duration(10 * (i + j + 1)),
									nil,
									time.Now(),
									nil,
								}
						}
					}

				}

			}
		}
		chain.Blocks = append(chain.Blocks, &Block{nodeIDs[i], latencies})
	}

	return &chain, nodeIDs

}

var lettersToNumbers = map[string]int{
	"A": 0,
	"B": 1,
	"C": 2,
	"D": 3,
	"E": 4,
	"F": 5,
	"G": 6,
	"H": 7,
	"I": 8,
	"J": 9,
	"K": 10,
	"L": 11,
	"M": 12,
	"N": 13,
	"O": 14,
	"P": 15,
	"Q": 16,
	"R": 17,
	"S": 18,
	"T": 19,
	"U": 20,
	"V": 21,
	"W": 22,
	"X": 23,
	"Y": 24,
	"Z": 25,
}

var numbersToLetters = map[int]string{
	0:  "A",
	1:  "B",
	2:  "C",
	3:  "D",
	4:  "E",
	5:  "F",
	6:  "G",
	7:  "H",
	8:  "I",
	9:  "J",
	10: "K",
	11: "L",
	12: "M",
	13: "N",
	14: "O",
	15: "P",
	16: "Q",
	17: "R",
	18: "S",
	19: "T",
	20: "U",
	21: "V",
	22: "W",
	23: "X",
	24: "Y",
	25: "Z",
}

func simpleChain(nbNodes int) (*Chain, []sigAlg.PublicKey) {
	blocks := make([]*Block, nbNodes)
	nodes := make([]sigAlg.PublicKey, nbNodes)

	for i := 0; i < nbNodes; i++ {
		latencies := make(map[string]ConfirmedLatency)

		for j := 0; j < nbNodes; j++ {
			if j != i {
				latencies[numbersToLetters[j]] = ConfirmedLatency{time.Duration(10 * time.Nanosecond), nil, time.Now(), nil}
			}
		}

		block := &Block{
			ID: &NodeID{
				ServerID:  nil,
				PublicKey: sigAlg.PublicKey(numbersToLetters[i]),
			},
			Latencies: latencies,
		}

		blocks[i] = block

		nodes[i] = sigAlg.PublicKey([]byte(numbersToLetters[i]))

	}

	chain := &Chain{
		Blocks:     blocks,
		BucketName: []byte("TestBucket"),
	}

	return chain, nodes
}

func setLiarAndVictim(chain *Chain, liar string, victim string, latency int) {
	chain.Blocks[lettersToNumbers[liar]].Latencies[victim] = ConfirmedLatency{time.Duration(time.Duration(latency) * time.Nanosecond), nil, time.Now(), nil}
	chain.Blocks[lettersToNumbers[victim]].Latencies[liar] = ConfirmedLatency{time.Duration(time.Duration(latency) * time.Nanosecond), nil, time.Now(), nil}

}
