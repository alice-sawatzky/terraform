package terraform

import (
	"log"

	"github.com/hashicorp/terraform/dag"
	"github.com/hashicorp/terraform/states"
)

// StateTransformer is a GraphTransformer that adds the elements of
// the state to the graph.
//
// This transform is used for example by the DestroyPlanGraphBuilder to ensure
// that only resources that are in the state are represented in the graph.
type StateTransformer struct {
	Concrete ConcreteResourceInstanceNodeFunc

	State *states.State
}

func (t *StateTransformer) Transform(g *Graph) error {
	if t.State.HasResources() {
		// Nothing to do!
		return nil
	}

	log.Printf("[TRACE] StateTransformer: starting")
	for _, ms := range t.State.Modules {
		moduleAddr := ms.Addr
		log.Printf("[TRACE] StateTransformer: Module: %v", moduleAddr)

		for _, rs := range ms.Resources {
			resourceAddr := rs.Addr.Absolute(moduleAddr)

			for key, is := range rs.Instances {
				addr := resourceAddr.Instance(key)
				log.Printf("[TRACE] StateTransformer: Resource %q: %#v", addr, rs)

				abstract := NewNodeAbstractResourceInstance(addr)
				var node dag.Vertex = abstract
				if f := t.Concrete; f != nil {
					node = f(abstract)
				}

				g.Add(node)
			}
		}
	}

	return nil
}
