package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
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
	variablesFile   string
	outputFile      string
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

	bt, err := convertHoneycombBoardToTemplate()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if variablesFile != "" {
		variables, err := loadVariables()
		if err != nil {
			fmt.Println(err)
			os.Exit(3)
		}

		bt.Variables = variables
	}

	tpl, err := generateTemplateGoCode(bt)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	if outputFile != "" {
		fmt.Println("Writing template to:", outputFile)

		f, err := os.Create(outputFile)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		defer f.Close()

		_, err = f.WriteString(tpl)
		if err != nil {
			fmt.Println(err)
			os.Exit(5)
		}

		fmt.Println("Done")
	} else {
		fmt.Println(tpl)
	}
}

func convertHoneycombBoardToTemplate() (*BoardTemplate, error) {

	if outputFile != "" {
		fmt.Println("Loading Honeycomb Board:", boardId)
	}

	hnyClient := NewHoneycombClient(honeycombApiKey)

	board, err := hnyClient.GetBoard(boardId)
	if err != nil {
		return nil, err
	}

	queryTemplates := make([]QueryTemplate, 0, len(board.Queries))
	for _, q := range board.Queries {
		query, err := hnyClient.GetQuery(q.Dataset, q.QueryId)
		if err != nil {
			return nil, err
		}

		queryAnnotation, err := hnyClient.GetQueryAnnotation(q.Dataset, q.QueryAnnotationId)
		if err != nil {
			return nil, err
		}

		queryTemplate := QueryTemplate{
			Name:             queryAnnotation.Name,
			ShortDescription: q.Caption,
			Description:      queryAnnotation.Description,
			Style:            q.QueryStyle,
			GraphSettings:    q.GraphSettings,
			QuerySpec: QuerySpec{
				Id:                        query.Id,
				StartTime:                 query.StartTime,
				EndTime:                   query.EndTime,
				DesiredGranularitySeconds: query.Granularity,
				Aggregates:                query.Calculations,
				FilterSet: &QueryFilterSet{
					Combination: query.FilterCombination,
					Filters:     query.Filters,
				},
				Groups:  query.Breakdowns,
				Orders:  query.Orders,
				Limit:   query.Limit,
				Havings: query.Havings,
			},
		}

		queryTemplates = append(queryTemplates, queryTemplate)
	}

	return &BoardTemplate{
		PK:             fmt.Sprint(sequenceNumber),
		Name:           board.Name,
		Description:    board.Description,
		Graphic:        graphic,
		ColumnStyle:    board.ColumnLayout,
		QueryTemplates: queryTemplates,
	}, nil
}

func loadVariables() ([]VariableSpec, error) {
	if outputFile != "" {
		fmt.Println("Loading Variables from:", variablesFile)
	}

	var variables VariablesDefinition
	file, err := os.Open(variablesFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if strings.HasSuffix(variablesFile, ".json") {
		err = json.NewDecoder(file).Decode(&variables)
	} else if strings.HasSuffix(variablesFile, ".yaml") || strings.HasSuffix(variablesFile, ".yml") {
		err = yaml.NewDecoder(file).Decode(&variables)
	} else {
		return nil, fmt.Errorf("unsupported file type: %s", variablesFile)
	}

	if err != nil {
		return nil, err
	}

	return variables.Variables, nil

}

func generateTemplateGoCode(bt *BoardTemplate) (string, error) {

	if outputFile != "" {
		fmt.Println("Generating Template Go Code")
	}

	codeName := firstLetterToLower(strings.ReplaceAll(bt.Name, " ", ""))

	tpl := ""

	tpl += "package templates\n"
	tpl += "\n"
	tpl += "import (\n"
	tpl += "	\"github.com/honeycombio/hound/api\"\n"
	tpl += "	\"github.com/honeycombio/hound/types\"\n"
	tpl += ")\n"
	tpl += "\n"

	tpl += "func " + codeName + "BoardTemplate() BoardTemplate {\n"
	tpl += "\tqueryTemplates := make([]QueryTemplate, 0, " + fmt.Sprint(len(bt.QueryTemplates)) + ")\n"
	tpl += "\n"

	for i, qt := range bt.QueryTemplates {
		// QuerySpec
		tpl += "\tqs" + fmt.Sprint(i+1) + " := api.QuerySpec{\n"

		// QuerySpec Aggregates
		if len(qt.QuerySpec.Aggregates) > 0 {
			tpl += "\t\tAggregates: []*api.QuerySpec_Aggregate{\n"
			for _, a := range qt.QuerySpec.Aggregates {
				tpl += "\t\t\t{\n"
				tpl += "\t\t\t\tOp: api.AggregateOp_" + a.Op + ",\n"
				tpl += "\t\t\t\tColumn: " + columnOrVariableName(a.Column, bt.Variables) + ",\n"
				tpl += "\t\t\t},\n"
			}
			tpl += "\t\t},\n"
		}

		// QuerySpec FilterSet
		if len(qt.QuerySpec.FilterSet.Filters) > 0 {
			tpl += "\t\tFilterSet: &api.QuerySpec_FilterSet{\n"
			tpl += "\t\t\tFilters: []*api.QuerySpec_Filter{\n"
			for _, f := range qt.QuerySpec.FilterSet.Filters {
				tpl += "\t\t\t\t{\n"
				tpl += "\t\t\t\t\tOp: api.FilterOpFromString(\"" + f.Op + "\"),\n"
				tpl += "\t\t\t\t\tColumn: " + columnOrVariableName(f.Column, bt.Variables) + ",\n"
				tpl += "\t\t\t\t\tValue: \"" + fmt.Sprint(f.Value) + "\",\n"
				tpl += "\t\t\t\t},\n"
			}
			tpl += "\t\t\t},\n"
			if qt.QuerySpec.FilterSet.Combination != "" {
				tpl += "\t\t\tCombination: api.FilterCombinationOp_" + qt.QuerySpec.FilterSet.Combination + ",\n"
			}
			tpl += "\t\t},\n"
		}

		// QuerySpec Groups
		if len(qt.QuerySpec.Groups) > 0 {
			tpl += "\t\tGroups: []string{\n"
			for _, g := range qt.QuerySpec.Groups {
				tpl += "\t\t\t" + columnOrVariableName(g, bt.Variables) + ",\n"
			}
			tpl += "\t\t},\n"
		}

		// QuerySpec Orders
		if len(qt.QuerySpec.Orders) > 0 {
			tpl += "\t\tOrders: []*api.QuerySpec_Order{\n"
			for _, o := range qt.QuerySpec.Orders {
				tpl += "\t\t\t{\n"
				tpl += "\t\t\t\tColumn: " + columnOrVariableName(o.Column, bt.Variables) + ",\n"
				tpl += "\t\t\t\tOp: api.AggregateOp_" + o.Op + ",\n"
				if o.Order == "descending" {
					tpl += "\t\t\t\tDescending: true,\n"
				} else {
					tpl += "\t\t\t\tDescending: false,\n"
				}
				tpl += "\t\t\t},\n"
			}
			tpl += "\t\t},\n"
		}

		// QuerySpec Limit
		if qt.QuerySpec.Limit > 0 {
			tpl += "\t\tLimit: " + fmt.Sprint(qt.QuerySpec.Limit) + ",\n"
		}

		// QuerySpec Havings
		if len(qt.QuerySpec.Havings) > 0 {
			tpl += "\t\tHavings: []*api.QuerySpec_Having{\n"
			for _, h := range qt.QuerySpec.Havings {
				tpl += "\t\t\t{\n"
				tpl += "\t\t\t\tAggregateOp: api.AggregateOp_" + h.CalculateOp + ",\n"
				tpl += "\t\t\t\tColumn: " + columnOrVariableName(h.Column, bt.Variables) + ",\n"
				tpl += "\t\t\t\tOp: api.FilterOp_" + h.Op + ",\n"
				tpl += "\t\t\t\tValue: \"" + fmt.Sprint(h.Value) + "\",\n"
				tpl += "\t\t\t\tJoinColumn: \"" + h.JoinColumn + "\",\n"
				tpl += "\t\t\t},\n"
			}
			tpl += "\t\t},\n"
		}

		// QuerySpec (close)
		tpl += "\t}\n"

		tpl += "\tqt" + fmt.Sprint(i+1) + " := QueryTemplate{\n"

		tpl += "\t\tName: \"" + qt.Name + "\",\n"
		tpl += "\t\tShortDescription: \"" + qt.ShortDescription + "\",\n"
		tpl += "\t\tDescription: \"" + qt.Description + "\",\n"
		tpl += "\t\tQuerySpec: qs" + fmt.Sprint(i+1) + ",\n"
		tpl += "\t\tStyle: types.BoardQueryStyle" + firstLetterToUpper(qt.Style) + ",\n"
		tpl += "\t\tGraphSettings: types.GraphSettings{\n"
		tpl += "\t\t\tOmitMissingValues: " + fmt.Sprint(qt.GraphSettings.OmitMissingValues) + ",\n"
		tpl += "\t\t\tUseStackedGraphs: " + fmt.Sprint(qt.GraphSettings.StackedGraphs) + ",\n"
		tpl += "\t\t\tUseLogScale: " + fmt.Sprint(qt.GraphSettings.LogScale) + ",\n"
		tpl += "\t\t\tUseUTCXAxis: " + fmt.Sprint(qt.GraphSettings.UTCXAxis) + ",\n"
		tpl += "\t\t\tHideMarkers: " + fmt.Sprint(qt.GraphSettings.HideMarkers) + ",\n"
		tpl += "\t\t\tPreferOverlaidCharts: " + fmt.Sprint(qt.GraphSettings.OverlaidCharts) + ",\n"
		tpl += "\t\t},\n"
		tpl += "\t\tAutoFilter: true,\n"
		tpl += "\t}\n"
		tpl += "\tqueryTemplates = append(queryTemplates, qt" + fmt.Sprint(i+1) + ")\n"
		tpl += "\n"
	}

	tpl += "\treturn BoardTemplate{\n"
	tpl += "\t\tPK: ToBoardTemplatePK(" + bt.PK + "),\n"
	tpl += "\t\tName: \"" + bt.Name + "\",\n"
	tpl += "\t\tDescription: \"" + bt.Description + "\",\n"
	tpl += "\t\tGraphic: " + fmt.Sprint(bt.Graphic) + ",\n"
	tpl += "\t\tQueryTemplates: queryTemplates,\n"
	tpl += "\t\tColumnStyle: types.BoardManyColumns,\n"
	tpl += "\t\tVariables: []VariableSpec{\n"
	for _, v := range bt.Variables {
		tpl += "\t\t\t{\n"
		tpl += "\t\t\t\tName: VariableName(\"" + v.Name + "\"),\n"
		tpl += "\t\t\t\tValueProviders: []ValueProvider{\n"
		for _, vp := range v.ValueProviders {
			tpl += "\t\t\t\t\t{\n"
			tpl += "\t\t\t\t\t\tKind: Column_" + vp.Kind + ",\n"
			tpl += "\t\t\t\t\t\tValue: \"" + vp.Value + "\",\n"
			tpl += "\t\t\t\t\t},\n"
		}
		tpl += "\t\t\t\t},\n"
		tpl += "\t\t\t},\n"
	}
	tpl += "\t\t},\n"
	tpl += "\t}\n"
	tpl += "}\n"

	return tpl, nil
}

func columnOrVariableName(column string, variables []VariableSpec) string {
	for _, v := range variables {
		if column == v.Name {
			return "VariableName(\"" + column + "\")"
		}
	}
	return "\"" + column + "\""
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
	flag.StringVar(&variablesFile, "variables", "", "Variables definition file to use")
	flag.StringVar(&outputFile, "out", "", "Output template fo file")
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
