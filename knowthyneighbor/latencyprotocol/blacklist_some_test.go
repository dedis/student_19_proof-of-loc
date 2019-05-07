package latencyprotocol

import (
	"testing"

	"go.dedis.ch/onet/v3/log"
)

//Problem: none of the tests create chains that blacklist nodes even with all latencies given
//Solution: this is the test for coordinated, proving detection does not work in this case
//Create other file with uncoherent connections instead, where it should work.
func TestBlacklist1Liar1Victim(t *testing.T) {
	N := 4

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N3", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}

func TestBlacklist2Liars3Victims(t *testing.T) {
	N := 8

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N2", 25)
	setLiarAndVictim(chain, "N1", "N2", 25)
	setLiarAndVictim(chain, "N0", "N3", 25)
	setLiarAndVictim(chain, "N1", "N3", 25)
	setLiarAndVictim(chain, "N0", "N4", 25)
	setLiarAndVictim(chain, "N1", "N4", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}

func TestBlacklist3Liars3Victims(t *testing.T) {
	N := 9

	chain, nodeIDs := simpleChain(N)

	//Liars: N0, N1, N2

	setLiarAndVictim(chain, "N0", "N3", 25)
	setLiarAndVictim(chain, "N0", "N4", 25)
	setLiarAndVictim(chain, "N0", "N5", 25)
	setLiarAndVictim(chain, "N1", "N3", 25)
	setLiarAndVictim(chain, "N1", "N4", 25)
	setLiarAndVictim(chain, "N1", "N5", 25)
	setLiarAndVictim(chain, "N2", "N3", 25)
	setLiarAndVictim(chain, "N2", "N4", 25)
	setLiarAndVictim(chain, "N2", "N5", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}

func TestBlacklist15Nodes4Liars5Victims(t *testing.T) {
	N := 15

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N10", 25)
	setLiarAndVictim(chain, "N1", "N10", 25)
	setLiarAndVictim(chain, "N2", "N10", 25)
	setLiarAndVictim(chain, "N3", "N10", 25)

	setLiarAndVictim(chain, "N0", "N6", 25)
	setLiarAndVictim(chain, "N1", "N6", 25)
	setLiarAndVictim(chain, "N2", "N6", 25)
	setLiarAndVictim(chain, "N3", "N6", 25)

	setLiarAndVictim(chain, "N0", "N7", 25)
	setLiarAndVictim(chain, "N1", "N7", 25)
	setLiarAndVictim(chain, "N2", "N7", 25)
	setLiarAndVictim(chain, "N3", "N7", 25)

	setLiarAndVictim(chain, "N0", "N8", 25)
	setLiarAndVictim(chain, "N1", "N8", 25)
	setLiarAndVictim(chain, "N2", "N8", 25)
	setLiarAndVictim(chain, "N3", "N8", 25)

	setLiarAndVictim(chain, "N0", "N9", 25)
	setLiarAndVictim(chain, "N1", "N9", 25)
	setLiarAndVictim(chain, "N2", "N9", 25)
	setLiarAndVictim(chain, "N3", "N9", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}

func TestBlacklist15Nodes3Liars5Victims(t *testing.T) {
	N := 15

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N10", 25)
	setLiarAndVictim(chain, "N1", "N10", 25)
	setLiarAndVictim(chain, "N2", "N10", 25)

	setLiarAndVictim(chain, "N0", "N6", 25)
	setLiarAndVictim(chain, "N1", "N6", 25)
	setLiarAndVictim(chain, "N2", "N6", 25)

	setLiarAndVictim(chain, "N0", "N7", 25)
	setLiarAndVictim(chain, "N1", "N7", 25)
	setLiarAndVictim(chain, "N2", "N7", 25)

	setLiarAndVictim(chain, "N0", "N8", 25)
	setLiarAndVictim(chain, "N1", "N8", 25)
	setLiarAndVictim(chain, "N2", "N8", 25)

	setLiarAndVictim(chain, "N0", "N9", 25)
	setLiarAndVictim(chain, "N1", "N9", 25)
	setLiarAndVictim(chain, "N2", "N9", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}

func TestBlacklist15Nodes4Liars3Victims(t *testing.T) {
	N := 15

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N10", 25)
	setLiarAndVictim(chain, "N1", "N10", 25)
	setLiarAndVictim(chain, "N2", "N10", 25)
	setLiarAndVictim(chain, "N3", "N10", 25)

	setLiarAndVictim(chain, "N0", "N6", 25)
	setLiarAndVictim(chain, "N1", "N6", 25)
	setLiarAndVictim(chain, "N2", "N6", 25)
	setLiarAndVictim(chain, "N3", "N6", 25)

	setLiarAndVictim(chain, "N0", "N7", 25)
	setLiarAndVictim(chain, "N1", "N7", 25)
	setLiarAndVictim(chain, "N2", "N7", 25)
	setLiarAndVictim(chain, "N3", "N7", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}

func TestBlacklist15Nodes5Liars3Victims(t *testing.T) {
	N := 15

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N10", 25)
	setLiarAndVictim(chain, "N1", "N10", 25)
	setLiarAndVictim(chain, "N2", "N10", 25)
	setLiarAndVictim(chain, "N3", "N10", 25)
	setLiarAndVictim(chain, "N4", "N10", 25)

	setLiarAndVictim(chain, "N0", "N6", 25)
	setLiarAndVictim(chain, "N1", "N6", 25)
	setLiarAndVictim(chain, "N2", "N6", 25)
	setLiarAndVictim(chain, "N3", "N6", 25)
	setLiarAndVictim(chain, "N4", "N6", 25)

	setLiarAndVictim(chain, "N0", "N7", 25)
	setLiarAndVictim(chain, "N1", "N7", 25)
	setLiarAndVictim(chain, "N2", "N7", 25)
	setLiarAndVictim(chain, "N3", "N7", 25)
	setLiarAndVictim(chain, "N4", "N7", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}
func TestBlacklist15Nodes5Liars1Victim(t *testing.T) {
	N := 15

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N10", 25)
	setLiarAndVictim(chain, "N1", "N10", 25)
	setLiarAndVictim(chain, "N2", "N10", 25)
	setLiarAndVictim(chain, "N3", "N10", 25)
	setLiarAndVictim(chain, "N4", "N10", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}

func TestBlacklist15Nodes4Liars8Victims(t *testing.T) {
	N := 15

	chain, nodeIDs := simpleChain(N)

	setLiarAndVictim(chain, "N0", "N10", 25)
	setLiarAndVictim(chain, "N1", "N10", 25)
	setLiarAndVictim(chain, "N2", "N10", 25)
	setLiarAndVictim(chain, "N3", "N10", 25)

	setLiarAndVictim(chain, "N0", "N6", 25)
	setLiarAndVictim(chain, "N1", "N6", 25)
	setLiarAndVictim(chain, "N2", "N6", 25)
	setLiarAndVictim(chain, "N3", "N6", 25)

	setLiarAndVictim(chain, "N0", "N7", 25)
	setLiarAndVictim(chain, "N1", "N7", 25)
	setLiarAndVictim(chain, "N2", "N7", 25)
	setLiarAndVictim(chain, "N3", "N7", 25)

	setLiarAndVictim(chain, "N0", "N8", 25)
	setLiarAndVictim(chain, "N1", "N8", 25)
	setLiarAndVictim(chain, "N2", "N8", 25)
	setLiarAndVictim(chain, "N3", "N8", 25)

	setLiarAndVictim(chain, "N0", "N9", 25)
	setLiarAndVictim(chain, "N1", "N9", 25)
	setLiarAndVictim(chain, "N2", "N9", 25)
	setLiarAndVictim(chain, "N3", "N9", 25)

	setLiarAndVictim(chain, "N0", "N13", 25)
	setLiarAndVictim(chain, "N1", "N13", 25)
	setLiarAndVictim(chain, "N2", "N13", 25)
	setLiarAndVictim(chain, "N3", "N13", 25)

	setLiarAndVictim(chain, "N0", "N11", 25)
	setLiarAndVictim(chain, "N1", "N11", 25)
	setLiarAndVictim(chain, "N2", "N11", 25)
	setLiarAndVictim(chain, "N3", "N11", 25)

	setLiarAndVictim(chain, "N0", "N12", 25)
	setLiarAndVictim(chain, "N1", "N12", 25)
	setLiarAndVictim(chain, "N2", "N12", 25)
	setLiarAndVictim(chain, "N3", "N12", 25)

	log.Print(checkBlacklistWithRemovedLatencies(chain, nodeIDs))

}
