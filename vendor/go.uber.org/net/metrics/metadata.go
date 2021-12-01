// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package metrics

import (
	"errors"
	"fmt"
	"sort"

	promproto "github.com/prometheus/client_model/go"
)

// Match the Prometheus error text.
var errInconsistentCardinality = errors.New("inconsistent tag cardinality")

// metadata stores our internal representation of metric specs. Adding this
// layer of indirection between the user-facing specs and the metric
// constructors serves two purposes: it centralizes the logic for calculating
// a variety of derived values, and it lets the remainder of the package
// assume that all user-supplied data has already been fully validated.
type metadata struct {
	Name, Help  *string // proto wants pointers
	Dims        string
	DisablePush bool

	constTagPairs []*promproto.LabelPair
	varTagNames   []string // unscrubbed
}

func newMetadata(o Spec) (metadata, error) {
	// TODO: Consider checking for duplicate tags with Bloom filters, allocating
	// maps only if we suspect a duplicate.
	sortedConstNames := make([]string, 0, len(o.ConstTags))
	for k := range o.ConstTags {
		sortedConstNames = append(sortedConstNames, k)
	}
	sort.Strings(sortedConstNames)

	constNameSet := make(map[string]struct{}, len(sortedConstNames))
	sortedScrubbedConstNames := make([]string, len(sortedConstNames))
	sortedScrubbedConstVals := make([]string, len(sortedConstNames))
	for i, name := range sortedConstNames {
		scrubbedName := scrubName(name)
		if _, ok := constNameSet[scrubbedName]; ok {
			return metadata{}, fmt.Errorf("duplicate constant tag name %q", scrubbedName)
		}
		constNameSet[scrubbedName] = struct{}{}
		sortedScrubbedConstNames[i] = scrubbedName
		sortedScrubbedConstVals[i] = scrubTagValue(o.ConstTags[name])
	}

	varNameSet := make(map[string]struct{}, len(o.VarTags))
	sortedScrubbedVarNames := make([]string, len(o.VarTags))
	for i, name := range o.VarTags {
		scrubbedName := scrubName(name)
		if _, ok := varNameSet[scrubbedName]; ok {
			return metadata{}, fmt.Errorf("duplicate variable tag name %q", scrubbedName)
		}
		if _, ok := constNameSet[scrubbedName]; ok {
			return metadata{}, fmt.Errorf("variable tag name %q is also a constant tag name", scrubbedName)
		}
		varNameSet[scrubbedName] = struct{}{}
		sortedScrubbedVarNames[i] = scrubbedName
	}
	sort.Strings(sortedScrubbedVarNames)

	var pairs []*promproto.LabelPair
	if len(sortedScrubbedConstNames) > 0 {
		pairs = make([]*promproto.LabelPair, 0, len(sortedScrubbedConstNames))
		for i := range sortedScrubbedConstNames {
			pairs = append(pairs, &promproto.LabelPair{
				Name:  &sortedScrubbedConstNames[i],
				Value: &sortedScrubbedConstVals[i],
			})
		}
	}
	scrubbedName := scrubName(o.Name)
	return metadata{
		Name:          &scrubbedName,
		Help:          &o.Help,
		Dims:          makeDims(scrubbedName, sortedScrubbedConstNames, sortedScrubbedVarNames),
		DisablePush:   o.DisablePush,
		constTagPairs: pairs,
		varTagNames:   o.VarTags, // preserve user-defined order
	}, nil
}

// MergeTags merges variable and constant tags.
func (m metadata) MergeTags(variableTagPairs []string) []*promproto.LabelPair {
	if len(variableTagPairs) == 0 {
		return m.constTagPairs
	}
	n := len(m.constTagPairs) + len(m.varTagNames)
	pairs := make([]*promproto.LabelPair, 0, n)
	pairs = append(pairs, m.constTagPairs...)
	for i := range m.varTagNames { // user-supplied order was preserved
		name := scrubName(m.varTagNames[i])
		val := scrubTagValue(variableTagPairs[i*2+1])
		pairs = append(pairs, &promproto.LabelPair{
			Name:  &name,
			Value: &val,
		})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].GetName() < pairs[j].GetName()
	})
	return pairs
}

// ValidateVariableTags checks that the user-supplied variable tag names
// and values match the spec supplied at vector creation.
func (m metadata) ValidateVariableTags(variableTagPairs []string) error {
	if len(variableTagPairs) != 2*len(m.varTagNames) {
		return errInconsistentCardinality
	}
	for i, expected := range m.varTagNames { // user-supplied order was preserved
		if expected != variableTagPairs[i*2] {
			return fmt.Errorf(
				"variable tag #%d doesn't match vector definition: expected %s, got %s",
				i,
				expected,
				variableTagPairs[i*2],
			)
		}
	}
	return nil
}

// writeID writes the metric's ID to the supplied digester. Since we only use
// IDs as map keys, we can save an allocation on metric creation by not
// allocating a string for each ID.
func (m metadata) writeID(d *digester) {
	d.add("", *m.Name)
	for _, pair := range m.constTagPairs {
		d.add("", pair.GetName())
		d.add("", pair.GetValue())
	}
}

// makeDims creates a string representation of the metric's dimensions (name,
// constant tag names, and variable tag names). It's used as a map value,
// so we can't avoid this allocation.
func makeDims(name string, constNames, sortedVarNames []string) string {
	d := newDigester()
	d.add("", name)
	for _, n := range constNames {
		d.add("", n)
	}
	for _, n := range sortedVarNames {
		// To make sure that we can tell whether a given tag name is constant or
		// variable, prefix variable tag names with a character that's otherwise
		// forbidden.
		d.add("$", n)
	}
	dims := string(d.digest())
	d.free()
	return dims
}
