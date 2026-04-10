package migrations

import (
	"bytes"
	"text/template"
)

const createLatestViewTmpl = `
CREATE VIEW v_{{.TableName}}_latest AS
SELECT sub.*
FROM (
    SELECT 
        t.*,
        ROW_NUMBER() OVER (
            PARTITION BY t.mrid
            ORDER BY c.created_at DESC
        ) AS row_num
    FROM {{.TableName}} t
    JOIN commits c 
        ON t.commit_id = c.id
) sub
WHERE sub.row_num = 1;
`

func MustGetViewSql(name string) string {
	data := struct {
		TableName string
	}{
		TableName: name,
	}
	templ, err := template.New("view").Parse(createLatestViewTmpl)
	if err != nil {
		panic("Could not parse template: " + err.Error())
	}

	var out bytes.Buffer

	if err := templ.Execute(&out, data); err != nil {
		panic("Could not execute template: " + err.Error())
	}
	return out.String()
}
