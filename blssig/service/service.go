// Package service implements a BLSCoSi service for which clients can connect to
// and then sign messages.
package service

import (
	"errors"
	"time"

	"github.com/dedis/student_19_proof-of-loc/blssig/protocol"
	"go.dedis.ch/cothority/v3/messaging"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
	uuid "gopkg.in/satori/go.uuid.v1"
)

// This file contains all the code to run a BLSCoSi service. It is used to reply to
// client request for signing something using BLSCoSiService.
// This is an updated version of the CoSi Service (dedis/cothority/cosi/service), which only does simple signing

// ServiceName is the name to refer to BLSCoSiService
const ServiceName = "BLSCoSiService"

func init() {
	log.Lvl3("Service: init")
	onet.RegisterNewService(ServiceName, newBLSCoSiService)
	network.RegisterMessage(&SignatureRequest{})
	network.RegisterMessage(&SignatureResponse{})
	network.RegisterMessage(&PropagationFunction{})
}

// BLSCoSiService is the service that handles collective signing operations
type BLSCoSiService struct {
	*onet.ServiceProcessor
	propagationFunction messaging.PropagationFunc
	propTimeout         time.Duration
	propagatedSignature []byte
}

// SignatureRequest is what the BLSCosi service is expected to receive from clients.
type SignatureRequest struct {
	Message []byte
	Roster  *onet.Roster
}

// SignatureResponse is what the BLSCosi service will reply to clients.
type SignatureResponse struct {
	Signature  []byte
	Propagated []byte
}

// PropagationFunction sends the complete signature to all members of the Cothority
type PropagationFunction struct {
	Signature []byte
}

// SignatureRequest treats external request to this service.
func (s *BLSCoSiService) SignatureRequest(req *SignatureRequest) (network.Message, error) {
	log.Lvl3("Service: SignatureRequest")
	if req.Roster.ID.IsNil() {
		req.Roster.ID = onet.RosterID(uuid.NewV4())
	}

	_, root := req.Roster.Search(s.ServerIdentity().ID)
	if root == nil {
		return nil, errors.New("Couldn't find a serverIdentity in Roster")
	}

	tree := req.Roster.GenerateNaryTreeWithRoot(2, root)
	tni := s.NewTreeNodeInstance(tree, tree.Root, protocol.Name)
	pi, err := protocol.NewDefaultProtocol(tni)
	if err != nil {
		return nil, errors.New("Couldn't make new protocol: " + err.Error())
	}
	s.RegisterProtocolInstance(pi)

	//Set message and start signing
	protocolInstance := pi.(*protocol.SimpleBLSCoSi)
	protocolInstance.Message = req.Message

	log.Lvl3("BLSCosi Service starting up root protocol")

	// start the protocol
	log.Lvl3("Cosi Service starting up root protocol")
	go pi.Dispatch()

	if err = pi.Start(); err != nil {
		return nil, err
	}

	sig := <-protocolInstance.FinalSignature

	// We propagate the signature to all nodes
	err = s.startPropagation(s.propagationFunction, req.Roster, &PropagationFunction{sig})
	if err != nil {
		return nil, err
	}

	propagated := s.propagatedSignature

	return &SignatureResponse{sig, propagated}, nil

}

// NewProtocol is called on all nodes of a Tree (except the root, since it is
// the one starting the protocol) so it's the Service that will be called to
// generate the PI on all others node.
func (s *BLSCoSiService) NewProtocol(tn *onet.TreeNodeInstance, conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	log.Lvl3("Service: NewProtocol")
	pi, err := protocol.NewDefaultProtocol(tn)
	return pi, err
}

func newBLSCoSiService(c *onet.Context) (onet.Service, error) {
	log.Lvl3("Service: newBLSCoSiService")
	s := &BLSCoSiService{
		ServiceProcessor: onet.NewServiceProcessor(c),
	}
	err := s.RegisterHandler(s.SignatureRequest)
	if err != nil {
		log.Error(err, "Couldn't register message:")
		return nil, err
	}

	s.propagationFunction, err = messaging.NewPropagationFunc(c, "propagateSignature", s.propagateFuncHandler, -1)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// SetPropTimeout is used to set the propagation timeout.
func (s *BLSCoSiService) SetPropTimeout(t time.Duration) {
	s.propTimeout = t
}

func (s *BLSCoSiService) startPropagation(propagate messaging.PropagationFunc, ro *onet.Roster, msg network.Message) error {

	replies, err := propagate(ro, msg, s.propTimeout)
	if err != nil {
		return err
	}

	if replies != len(ro.List) {
		log.Lvl1(s.ServerIdentity(), "Only got", replies, "out of", len(ro.List))
	}

	return nil
}

// propagateForwardLinkHandler will update the latest block with
// the new forward link and the new block when given
func (s *BLSCoSiService) propagateFuncHandler(msg network.Message) {
	s.propagatedSignature = msg.(*PropagationFunction).Signature
}
