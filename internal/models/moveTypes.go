package models

type Move struct {
        Name            string
        Input           string
        Headers			[]string
        FrameData 		map[string]Cell
        FrameDataGrids	[]FrameDataGrid
        FrameDataRows	[]FrameDataRow
        Notes           []string
        Description     string
}

type Cell struct {
        Value           string
        Tooltip         string
}

type FrameDataRow struct {
		Cells		[]Cell
}

type FrameDataGrid struct {
	Headers []string
	Rows	[]FrameDataRow
}
