package pkg

import (
	"com.github/davidkleiven/tripleworks/models"
	"github.com/google/uuid"
	"gonum.org/v1/gonum/graph/path"
	"gonum.org/v1/gonum/graph/simple"
)

type EquipmentConnector struct {
	terminals []models.Terminal
	idMap     map[uuid.UUID]int64
	graph     *simple.UndirectedGraph
}

func (e *EquipmentConnector) IsConnected(mrid1, mrid2 uuid.UUID, maxStep int) bool {
	n1, ok := e.idMap[mrid1]
	if !ok {
		return false
	}

	n2, ok := e.idMap[mrid2]
	if !ok {
		return false
	}

	node1 := e.graph.Node(n1)
	sp, ok := path.BellmanFordFrom(node1, e.graph)
	if !ok {
		return false
	}

	distance := sp.WeightTo(n2)
	return distance < float64(maxStep)
}

func (e *EquipmentConnector) GetTerminal(mrid uuid.UUID) *models.Terminal {
	for _, terminal := range e.terminals {
		if terminal.ConductingEquipmentMrid == mrid {
			return &terminal
		}
	}
	return nil
}

type ConnectParams struct {
	CreateSeqNo1       int // If created, this is the produced sequence number
	CreateSeqNo2       int // If created, use this sequence number
	Mrid1              uuid.UUID
	Mrid2              uuid.UUID
	ReportingGroupMrid uuid.UUID
	VoltageLevel       models.VoltageLevel
}

func (e *EquipmentConnector) Connect(params *ConnectParams) *ConnectionResult {
	result := ConnectionResult{
		Terminals:         []models.Terminal{},
		BusNameMarkers:    []models.BusNameMarker{},
		ConnectivityNodes: []models.ConnectivityNode{},
	}

	var (
		cn1Mrid uuid.UUID
		cn2Mrid uuid.UUID
	)

	term1 := e.GetTerminal(params.Mrid1)
	name := params.VoltageLevel.Name
	if term1 != nil {
		cn1Mrid = term1.ConnectivityNodeMrid
	} else {
		cn := CreateConnectivityNode(name)
		cn1Mrid = cn.Mrid
		result.ConnectivityNodes = append(result.ConnectivityNodes, cn)

		bnm := CreateBusNameMarker(name, params.ReportingGroupMrid)
		result.BusNameMarkers = append(result.BusNameMarkers, bnm)

		newTerm := CreateTerminal(cn.Mrid, params.Mrid1, bnm, params.CreateSeqNo1)
		result.Terminals = append(result.Terminals, newTerm)
	}

	term2 := e.GetTerminal(params.Mrid2)
	if term2 != nil {
		cn2Mrid = term2.ConnectivityNodeMrid
	} else {
		cn := CreateConnectivityNode(name)
		cn2Mrid = cn.Mrid
		result.ConnectivityNodes = append(result.ConnectivityNodes, cn)

		bnm := CreateBusNameMarker(name, params.ReportingGroupMrid)
		result.BusNameMarkers = append(result.BusNameMarkers, bnm)

		newTerm := CreateTerminal(cn.Mrid, params.Mrid2, bnm, params.CreateSeqNo2)
		result.Terminals = append(result.Terminals, newTerm)
	}

	// Create switch
	result.Switch = CreateSwitch(name, &params.VoltageLevel)
	bnm1 := CreateBusNameMarker(name, params.ReportingGroupMrid)
	result.BusNameMarkers = append(result.BusNameMarkers, bnm1)

	switchTerm1 := CreateTerminal(cn1Mrid, result.Switch.Mrid, bnm1, 1)

	bnm2 := CreateBusNameMarker(name, params.ReportingGroupMrid)
	result.BusNameMarkers = append(result.BusNameMarkers, bnm2)

	switchTerm2 := CreateTerminal(cn2Mrid, result.Switch.Mrid, bnm2, 2)
	result.Terminals = append(result.Terminals, switchTerm1, switchTerm2)
	return &result
}

func (e *EquipmentConnector) AddTerminals(terminals ...models.Terminal) {
	e.terminals = append(e.terminals, terminals...)
	for _, terminal := range terminals {
		nodeId1, ok := e.idMap[terminal.ConnectivityNodeMrid]
		if !ok {
			node := e.graph.NewNode()
			nodeId1 = node.ID()
			e.idMap[terminal.ConnectivityNodeMrid] = nodeId1
			e.graph.AddNode(node)
		}

		nodeId2, ok := e.idMap[terminal.ConductingEquipmentMrid]
		if !ok {
			node := e.graph.NewNode()
			nodeId2 = node.ID()
			e.idMap[terminal.ConductingEquipmentMrid] = nodeId2
			e.graph.AddNode(node)
		}
		edge := e.graph.NewEdge(e.graph.Node(nodeId1), e.graph.Node(nodeId2))
		e.graph.SetEdge(edge)
	}
}

type ConnectionResult struct {
	Terminals         []models.Terminal
	BusNameMarkers    []models.BusNameMarker
	ConnectivityNodes []models.ConnectivityNode
	Switch            models.Switch
}

func NewEmptyConnector() *EquipmentConnector {
	return &EquipmentConnector{terminals: []models.Terminal{}, idMap: make(map[uuid.UUID]int64), graph: simple.NewUndirectedGraph()}
}
