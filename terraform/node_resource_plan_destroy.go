package terraform

import (
	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/plans"
	"github.com/hashicorp/terraform/states"
)

// NodePlanDestroyableResourceInstance represents a resource that is ready
// to be planned for destruction.
type NodePlanDestroyableResourceInstance struct {
	*NodeAbstractResourceInstance
}

var (
	_ GraphNodeSubPath              = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeReferenceable        = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeReferencer           = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeDestroyer            = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeResource             = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeResourceInstance     = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeAttachResourceConfig = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeAttachResourceState  = (*NodePlanDestroyableResourceInstance)(nil)
	_ GraphNodeEvalable             = (*NodePlanDestroyableResourceInstance)(nil)
)

// GraphNodeDestroyer
func (n *NodePlanDestroyableResourceInstance) DestroyAddr() *addrs.AbsResourceInstance {
	addr := n.ResourceInstanceAddr()
	return &addr
}

// GraphNodeEvalable
func (n *NodePlanDestroyableResourceInstance) EvalTree() EvalNode {
	addr := n.ResourceInstanceAddr()

	// State still uses legacy-style internal ids, so we need to shim to get
	// a suitable key to use.
	stateId := NewLegacyResourceInstanceAddress(addr).stateId()

	// Declare a bunch of variables that are used for state during
	// evaluation. These are written to by address in the EvalNodes we
	// declare below.
	var change *plans.ResourceInstanceChange
	var state *states.ResourceInstanceObject

	return &EvalSequence{
		Nodes: []EvalNode{
			&EvalReadState{
				Name:   stateId,
				Output: &state,
			},
			&EvalDiffDestroy{
				Addr:   addr.Resource,
				State:  &state,
				Output: &change,
			},
			&EvalCheckPreventDestroy{
				Addr:   addr.Resource,
				Config: n.Config,
				Change: &change,
			},
			&EvalWriteDiff{
				Name:   stateId,
				Change: &change,
			},
		},
	}
}
