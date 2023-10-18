package main

type HoneycombBoard struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Style        string `json:"style"`
	ColumnLayout string `json:"column_layout"`
	Queries      []struct {
		Caption           string        `json:"caption"`
		GraphSettings     GraphSettings `json:"graph_settings"`
		QueryStyle        string        `json:"query_style"`
		Dataset           string        `json:"dataset"`
		QueryId           string        `json:"query_id"`
		QueryAnnotationId string        `json:"query_annotation_id"`
	} `json:"queries"`
	Links struct {
		BoardURL string `json:"board_url"`
	} `json:"links"`
}

type HoneycombQuery struct {
	Id                string               `json:"id"`
	Calculations      []QueryVisualization `json:"calculations"`
	Filters           []QueryFilter        `json:"filters"`
	FilterCombination string               `json:"filter_combination"`
	Breakdowns        []string             `json:"breakdowns"`
	Orders            []QueryOrder         `json:"orders"`
	Limit             int                  `json:"limit"`
	Havings           []QueryHaving        `json:"havings"`
	Granularity       int                  `json:"granularity"`
	StartTime         int                  `json:"start_time"`
	EndTime           int                  `json:"end_time"`
	TimeRange         int                  `json:"time_range"`
}

type HoneycombQueryAnnotation struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	QueryId     string `json:"query_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type HoneycombDerivedColumn struct {
	Id          string `json:"id"`
	Alias       string `json:"alias"`
	Expression  string `json:"expression"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type QueryVisualization struct {
	Op     string `json:"op"`
	Column string `json:"column"`
}

type QueryFilter struct {
	Op         string `json:"op"`
	Column     string `json:"column"`
	Value      any    `json:"value"`
	JoinColumn string `json:"join_column,omitempty"`
}

type QueryFilterSet struct {
	Combination string        `json:"combination"`
	Filters     []QueryFilter `json:"filters"`
}

type QueryOrder struct {
	Column string `json:"column"`
	Op     string `json:"op"`
	Order  string `json:"order"`
}

type QueryHaving struct {
	CalculateOp string `json:"calculate_op"`
	Column      string `json:"column"`
	Op          string `json:"op"`
	Value       any    `json:"value"`
	JoinColumn  string `json:"join_column,omitempty"`
}

type GraphSettings struct {
	OmitMissingValues bool `json:"omit_missing_values,omitempty"`
	StackedGraphs     bool `json:"stacked_graphs,omitempty"`
	LogScale          bool `json:"log_scale,omitempty"`
	UTCXAxis          bool `json:"utc_xaxis,omitempty"`
	OverlaidCharts    bool `json:"overlaid_charts,omitempty"`
	HideMarkers       bool `json:"hide_markers,omitempty"`
}

type BoardTemplate struct {
	PK             string          `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	ColumnStyle    string          `json:"columnStyle"`
	Graphic        int             `json:"graphic"`
	QueryTemplates []QueryTemplate `json:"queryTemplates"`
	Variables      []VariableSpec  `json:"variables"`
}

type QueryTemplate struct {
	Name             string        `json:"name"`
	Description      string        `json:"description"`
	ShortDescription string        `json:"shortDescription"`
	Style            string        `json:"style"`
	GraphSettings    GraphSettings `json:"graphSettings"`
	QuerySpec        QuerySpec     `json:"querySpec"`
}

type QuerySpec struct {
	Id                        string               `json:"id,omitempty"`
	StartTime                 int                  `json:"start_time,omitempty"`
	EndTime                   int                  `json:"end_time,omitempty"`
	DesiredGranularitySeconds int                  `json:"desired_granularity_seconds,omitempty"`
	Aggregates                []QueryVisualization `json:"aggregates,omitempty"`
	FilterSet                 *QueryFilterSet      `json:"filter_set,omitempty"`
	Groups                    []string             `json:"groups,omitempty"`
	Orders                    []QueryOrder         `json:"orders,omitempty"`
	Limit                     int                  `json:"limit,omitempty"`
	Havings                   []QueryHaving        `json:"havings,omitempty"`
}

type VariableSpec struct {
	Name           string `json:"name" yaml:"name"`
	ValueProviders []struct {
		Kind  string `json:"kind" yaml:"kind"`
		Value string `json:"value" yaml:"value"`
	} `json:"valueProviders" yaml:"valueProviders"`
}

type VariablesDefinition struct {
	Variables []VariableSpec `json:"variables" yaml:"variables"`
}

var (
	filterOp_display = map[string]string{
		"FilterOp_EQUAL":         "=",
		"FilterOp_NOT_EQUAL":     "!=",
		"FilterOp_GT":            ">",
		"FilterOp_GTE":           ">=",
		"FilterOp_LT":            "<",
		"FilterOp_LTE":           "<=",
		"FilterOp_PREFIX":        "starts-with",
		"FilterOp_NOT_PREFIX":    "does-not-start-with",
		"FilterOp_HAS_VALUE":     "exists",
		"FilterOp_NOT_HAS_VALUE": "does-not-exist",
		"FilterOp_CONTAINS":      "contains",
		"FilterOp_NOT_CONTAINS":  "does-not-contain",
		"FilterOp_IN_RESULT":     "in-result",
		"FilterOp_JOIN_RESULT":   "join-result",
		"FilterOp_IN":            "in",
		"FilterOp_NOT_IN":        "not-in",
	}
	filterOp_fromDisplay = (func() map[string]string {
		mp := make(map[string]string, len(filterOp_display))
		for k, v := range filterOp_display {
			mp[v] = k
		}
		return mp
	}())
)

func FilterOpFromString(str string) string {
	if op, ok := filterOp_fromDisplay[str]; ok {
		return op
	}
	return "FilterOp_UNKNOWN"
}
