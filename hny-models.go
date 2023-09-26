package main

type HoneycombBoard struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Style        string `json:"style"`
	ColumnLayout string `json:"column_layout"`
	Queries      []struct {
		Caption       string `json:"caption"`
		GraphSettings struct {
			HideMarkers       bool `json:"hide_markers"`
			LogScale          bool `json:"log_scale"`
			OmitMissingValues bool `json:"omit_missing_values"`
			StackedGraphs     bool `json:"stacked_graphs"`
			UtcAxis           bool `json:"utc_axis"`
			OverlaidCharts    bool `json:"overlaid_charts"`
		} `json:"graph_settings"`
		QueryStyle        string `json:"query_style"`
		Dataset           string `json:"dataset"`
		QueryId           string `json:"query_id"`
		QueryAnnotationId string `json:"query_annotation_id"`
	} `json:"queries"`
	Links struct {
		BoardURL string `json:"board_url"`
	} `json:"links"`
}

type HoneycombQuery struct {
	Id           string `json:"id"`
	Calculations []struct {
		Op     string `json:"op"`
		Column string `json:"column"`
	} `json:"calculations"`
	Filters []struct {
		Op     string      `json:"op"`
		Column string      `json:"column"`
		Value  interface{} `json:"value"`
	} `json:"filters"`
	FilterCombination string   `json:"filter_combination"`
	Breakdowns        []string `json:"breakdowns"`
	Orders            []struct {
		Column string `json:"column"`
		Op     string `json:"op"`
		Order  string `json:"order"`
	} `json:"orders"`
	Limit   int `json:"limit"`
	Havings []struct {
		CalculateOp string `json:"calculate_op"`
		Column      string `json:"column"`
		Op          string `json:"op"`
		Value       int    `json:"value"`
	} `json:"havings"`
	Granularity int `json:"granularity"`
	StartTime   int `json:"start_time"`
	EndTime     int `json:"end_time"`
	TimeRange   int `json:"time_range"`
}

type HoneycombQueryAnnotation struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	QueryId     string `json:"query_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type HoneycombBoardWithDetails struct {
	Board            *HoneycombBoard
	Queries          []HoneycombQuery
	QueryAnnotations []HoneycombQueryAnnotation
}
