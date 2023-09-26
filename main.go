package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"unicode"
)

var VERSION = "UNKNOWN"

var (
	honeycombApiKey string
	boardId         string
	graphic         = 1
	sequenceNumber  = 99999
	printVersion    = false
)

func main() {
	err := validateOptions()
	if err != nil {
		fmt.Println()
		fmt.Println(err)
		os.Exit(1)
	}

	if printVersion {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	b, err := loadHoneycombBoard()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	tpl, err := generateBoardTemplate(b)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	fmt.Println(tpl)
}

func loadHoneycombBoard() (*HoneycombBoardWithDetails, error) {

	hnyClient := NewHoneycombClient(honeycombApiKey)

	board, err := hnyClient.GetBoard(boardId)
	if err != nil {
		return nil, err
	}

	queries := make([]HoneycombQuery, 0, len(board.Queries))
	queryAnnotations := make([]HoneycombQueryAnnotation, 0, len(board.Queries))
	for _, q := range board.Queries {
		query, err := hnyClient.GetQuery(q.Dataset, q.QueryId)
		if err != nil {
			return nil, err
		}
		queries = append(queries, *query)

		queryAnnotation, err := hnyClient.GetQueryAnnotation(q.Dataset, q.QueryAnnotationId)
		if err != nil {
			return nil, err
		}
		queryAnnotations = append(queryAnnotations, *queryAnnotation)
	}

	return &HoneycombBoardWithDetails{
		Board:            board,
		Queries:          queries,
		QueryAnnotations: queryAnnotations,
	}, nil

}

func generateBoardTemplate(b *HoneycombBoardWithDetails) (string, error) {
	codeName := firstLetterToLower(strings.ReplaceAll(b.Board.Name, " ", ""))

	tpl := ""

	tpl += "package templates\n"
	tpl += "\n"
	tpl += "import (\n"
	tpl += "	\"github.com/honeycombio/hound/api\"\n"
	tpl += "	\"github.com/honeycombio/hound/types\"\n"
	tpl += ")\n"
	tpl += "\n"

	tpl += "func " + codeName + "BoardTemplate() BoardTemplate {\n"
	tpl += "\tqueryTemplates := []QueryTemplate{}\n"
	tpl += "\n"

	for i, bq := range b.Board.Queries {
		qs := fmt.Sprint("qs", i+1)
		tpl += "\t" + qs + " := api.QuerySpec{\n"
		tpl += "\t\tAggregates: []*api.QuerySpec_Aggregate{\n"
		for _, c := range b.Queries[i].Calculations {
			tpl += "\t\t\t{Op: api.AggregateOp_" + c.Op + ", Column: \"" + c.Column + "\"},\n"
		}
		tpl += "\t\t},\n"
		tpl += "\t\tGroups: []string{\n"
		for _, g := range b.Queries[i].Breakdowns {
			tpl += "\t\t\t\"" + g + "\",\n"
		}
		tpl += "\t\t},\n"
		tpl += "\t}\n"

		qt := fmt.Sprint("qt", i+1)
		tpl += "\t" + qt + " := QueryTemplate{\n"

		tpl += "\t\tName: \"" + b.QueryAnnotations[i].Name + "\",\n"
		tpl += "\t\tShortDescription: \"" + bq.Caption + "\",\n"
		tpl += "\t\tDescription: \"" + b.QueryAnnotations[i].Description + "\",\n"
		tpl += "\t\tQuerySpec: " + qs + ",\n"
		tpl += "\t\tStyle: types.BoardQueryStyle" + firstLetterToUpper(bq.QueryStyle) + ",\n"
		tpl += "\t\tGraphSettings: &types.GraphSettings{\n"
		tpl += "\t\t\tOmitMissingValues: " + fmt.Sprint(bq.GraphSettings.OmitMissingValues) + ",\n"
		tpl += "\t\t\tUseStackedGraphs: " + fmt.Sprint(bq.GraphSettings.StackedGraphs) + ",\n"
		tpl += "\t\t\tUseLogScale: " + fmt.Sprint(bq.GraphSettings.LogScale) + ",\n"
		tpl += "\t\t\tUseUTCAxis: " + fmt.Sprint(bq.GraphSettings.UtcAxis) + ",\n"
		tpl += "\t\t\tHideMarkers: " + fmt.Sprint(bq.GraphSettings.HideMarkers) + ",\n"
		tpl += "\t\t\tPreferOverlaidCharts: " + fmt.Sprint(bq.GraphSettings.OverlaidCharts) + ",\n"
		tpl += "\t\t},\n"
		tpl += "\t\tAutoFilter: true,\n"
		tpl += "\t}\n"
		tpl += "\tqueryTemplates = append(queryTemplates, " + qt + ")\n"
		tpl += "\n"
	}

	tpl += "\treturn BoardTemplate{\n"
	tpl += "\t\tPK: ToBoardTemplatePK(" + fmt.Sprint(sequenceNumber) + "),\n"
	tpl += "\t\tName: \"" + b.Board.Name + "\",\n"
	tpl += "\t\tDescription: \"" + b.Board.Description + "\",\n"
	tpl += "\t\tGraphic: " + fmt.Sprint(graphic) + ",\n"
	tpl += "\t\tQueryTemplates: queryTemplates,\n"
	tpl += "\t\tColumnStyle: types.BoardManyColumns,\n"
	tpl += "\t}\n"
	tpl += "}\n"

	return tpl, nil
}

func firstLetterToLower(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])

	return string(r)
}

func firstLetterToUpper(s string) string {
	if len(s) == 0 {
		return s
	}

	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])

	return string(r)
}

func validateOptions() error {
	flag.StringVar(&honeycombApiKey, "honeycomb-api-key", lookupEnvOrString("HONEYCOMB_API_KEY", honeycombApiKey), "Honeycomb API Key")
	flag.StringVar(&boardId, "board", "", "Honeycomb Board ID")
	flag.IntVar(&graphic, "graphic", graphic, "Graphic # to use")
	flag.IntVar(&sequenceNumber, "sequence-number", sequenceNumber, "Sequence number to use")
	flag.BoolVar(&printVersion, "version", false, "Print version")
	flag.Parse()

	if printVersion {
		return nil
	}

	if honeycombApiKey == "" {
		printUsage()
		return fmt.Errorf("missing: Honeycomb API Key")
	}

	if boardId == "" {
		printUsage()
		return fmt.Errorf("missing: Honeycomb Board ID")
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage: hny-btgen [options]")
	flag.PrintDefaults()
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
