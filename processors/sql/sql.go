//go:generate bitfanDoc
package sqlprocessor

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/ShowMax/go-fqdn"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vjeantet/bitfan/commons"
	"github.com/vjeantet/bitfan/processors"
	"strconv"
	"strings"
	"text/template"
	"time"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// GOLANG driver class to load, for example, "mysql".
	// @ExampleLS driver => "mysql"
	Driver string `mapstructure:"driver" validate:"required"`

	// Send an event row by row or one event with all results
	// possible values "row", "result"
	// @Default "row"
	// @Enum "row","result"
	EventBy string `mapstructure:"event_by"`

	// SQL Statement
	//
	// When there are more than 1 statement, only data from the last one will generate events.
	// @ExampleLS statement => "SELECT * FROM mytable"
	// @Type Location
	Statement string `mapstructure:"statement" validate:"required"`
	Offset    int64
	Time      string `mapstructure:"time_field"`

	// Set an interval when this processor is used as a input
	// @ExampleLS interval => "10"
	// @Type Interval
	Interval string `mapstructure:"interval" `

	// @ExampleLS connection_string => "username:password@tcp(192.168.1.2:3306)/mydatabase?charset=utf8"
	ConnectionString string `mapstructure:"connection_string" validate:"required"`

	// You can set variable to be used in Statements by using ${var}.
	// each reference will be replaced by the value of the variable found in Statement's content
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`

	// Define the target field for placing the retrieved data. If this setting is omitted,
	// the data will be stored in the "data" field
	// Set the value to "." to store value to the root (top level) of the event
	// @ExampleLS target => "data"
	// @Default "data"
	Target string `mapstructure:"target"`
}

type processor struct {
	processors.Base
	db            *sql.DB
	opt           *options
	q             chan bool
	host          string
	StatementTmpl *template.Template
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		EventBy: "row",
		Target:  "data",
	}

	p.opt = &defaults
	p.host = fqdn.Get()

	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.Interval == "" {
		p.Logger.Warningln("No interval set")
	}

	loc, err := commons.NewLocation(p.opt.Statement, p.ConfigWorkingLocation)
	if err != nil {
		return err
	}

	content, _, err := loc.ContentWithOptions(p.opt.Var)
	if err != nil {
		return err
	}
	p.opt.Statement = string(content)

	p.StatementTmpl, err = template.New("statement").Parse(p.opt.Statement)
	if err != nil {
		p.Logger.Errorf("sql Statement tpl error : %v", err)
		return err
	}

	p.db, err = sql.Open(p.opt.Driver, p.opt.ConnectionString)
	if err != nil {
		return err
	}

	return p.db.Ping()
}

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func getGoValue(sqlType, val string) interface{} {
	switch strings.ToLower(sqlType) {
	case "tinyint", "smallint", "mediumint", "int", "integer", "bigint":
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	case "date", "time", "year", "datetime", "timestamp":
		return val
	case "varchar", "char", "text", "tinytext", "mediumtext", "longtext":
		return val
	case "tinyblob", "blob", "mediumblob", "longblob":
		return val
	case "float":
		i, _ := strconv.ParseFloat(val, 32)
		return i
	case "double", "decimal":
		i, _ := strconv.ParseFloat(val, 64)
		return i
	default:
		return val
	}
}

func (p *processor) getsqlType() map[string]string {
	maps := map[string]string{}

	sql := "select data_type, column_name from information_schema.columns where table_schema=? and table_name=?"
	str := strings.Split(p.opt.ConnectionString, "/")[1]
	strArr := strings.Split(strings.TrimSpace(p.opt.Statement), " ")
	rows, err := p.db.Query(sql, str[:strings.Index(str, "?")], strArr[len(strArr)-1])
	if err != nil {
		fmt.Println(err)
		return maps
	}
	columns, err := rows.Columns()
	if err != nil {
		fmt.Println(err)
		return maps
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println(err)
			return maps
		}
		key, val := "", ""
		for i, col := range values {
			if col != nil {
				colStr := string(col.([]byte))
				switch strings.ToLower(columns[i]) {
				case "data_type":
					val = colStr
				case "column_name":
					key = colStr
				}
			}
		}

		maps[key] = val
	}
	return maps
}

var batchSize = 200

func (p *processor) Receive(e processors.IPacket) error {
	buff := bytes.NewBufferString("")
	p.StatementTmpl.Execute(buff, e.Fields())

	statements := strings.Trim(buff.String(), ";")
	reqs := strings.Split(statements, ";")
	if len(reqs) > 1 {
		for _, r := range reqs[:len(reqs)-1] {
			p.Logger.Debugf("db.Exec - %s", r)
			p.db.Exec(r)
		}
	}

	p.Logger.Debugf("db.Query - %s", reqs[len(reqs)-1])

	maps := p.getsqlType()
	total := 0

	rows, err := p.db.Query(reqs[len(reqs)-1] + fmt.Sprint(" limit ", p.opt.Offset, ",", batchSize))
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	var records []map[string]interface{}
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		total++
		record := make(map[string]interface{})
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}
		for i, col := range values {
			if col != nil {
				// fmt.Printf("\n%s: type= %s\n", columns[i], reflect.TypeOf(col))
				switch t := col.(type) {
				default:
					fmt.Printf("Unexpected type %T\n", t)
				case bool:
					record[columns[i]] = col.(bool)
				case int:
					record[columns[i]] = col.(int)
				case int64:
					record[columns[i]] = col.(int64)
				case float64:
					record[columns[i]] = col.(float64)
				case string:
					record[columns[i]] = col.(string)
				case []byte: // -- all cases go HERE!
					record[columns[i]] = getGoValue(maps[columns[i]], string(col.([]byte)))
					//case time.Time:
					// record[columns[i]] = col.(string)
				}
			}
		}

		if p.opt.Time != "" {
			if _, ok := record[p.opt.Time].(int64); ok {
				record["__time"] = record[p.opt.Time]
			} else if timeStr, ok := record[p.opt.Time].(string); ok {
				if timeStr == "" {
					record["__time"] = time.Now().Unix() * 1000
				} else {
					loc, _ := time.LoadLocation("Asia/Shanghai")
					format := ""
					switch maps[p.opt.Time] {
					case "time":
						timeStr = fmt.Sprint(time.Now().Format("20060102"), " ", timeStr)
						format = "20060102 15:04:05"
					case "year":
						format = "2006"
					case "date":
						format = "2006-01-02"
					case "timestamp", "datetime":
						format = "2006-01-02 15:04:05"
					}
					timeObj, err := time.ParseInLocation(format, timeStr, loc)
					if err != nil {
						fmt.Println("date trans err:", err)
						timeObj = time.Now()
					}
					record["__time"] = timeObj.Unix() * 1000
				}
			}
		} else {
			record["__time"] = time.Now().Unix() * 1000
		}

		if p.opt.EventBy == "row" {
			var e2 processors.IPacket = e.Clone()
			e2.Fields().SetValueForPath(p.host, "host")
			if len(p.opt.Var) > 0 {
				e2.Fields().SetValueForPath(p.opt.Var, "var")
			}

			if p.opt.Target == "." {
				for k, v := range record {
					e2.Fields().SetValueForPath(v, k)
				}
			} else {
				e2.Fields().SetValueForPath(record, p.opt.Target)
			}

			p.opt.ProcessCommonOptions(e2.Fields())
			p.Send(e2)
		} else {
			records = append(records, record)
		}
	}
	rows.Close()

	if total == 0 {
		return nil
	}

	if p.opt.EventBy != "row" {
		e.Fields().SetValueForPath(p.host, "host")
		if len(p.opt.Var) > 0 {
			e.Fields().SetValueForPath(p.opt.Var, "var")
		}
		e.Fields().SetValueForPath(records, p.opt.Target)

		p.opt.ProcessCommonOptions(e.Fields())
		p.Send(e)
	}
	p.opt.Offset += int64(total)
	if total == batchSize {
		return fmt.Errorf("goon")
	}
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.db.Close()
	return nil
}
