package logic

import (
	"context"
	"github.com/tal-tech/cds/galaxy/internal/svc"
	"github.com/tal-tech/cds/galaxy/internal/types"
	"github.com/tal-tech/cds/tools/mongodbx"
	"github.com/tal-tech/cds/tools/mysqlx"
	"strings"

	"github.com/tal-tech/go-zero/core/logx"
)

type GenerateCreateSqlLogic struct {
	ctx context.Context
	logx.Logger
}

func NewGenerateCreateSqlLogic(ctx context.Context, svcCtx *svc.ServiceContext) GenerateCreateSqlLogic {
	return GenerateCreateSqlLogic{
		ctx:    ctx,
		Logger: logx.WithContext(ctx),
	}
	// TODO need set model here from svc
}

func (l *GenerateCreateSqlLogic) GenerateCreateSql(req types.GenerateCreateSqlRequest) (*types.GenerateCreateSqlResponse, error) {
	var failedtables, reasons []string
	sqls := make([]string, 0, len(req.Table))
	qks := make([]string, 0, len(req.Table))
	for _, v := range req.Table {
		var qk string
		var sql []string
		var e error
		if strings.HasPrefix(req.Dsn, "mongodb://") {
			sql, qk, e = mongodbx.ToClickhouseTable(req.Dsn, req.Database, v, "")
			if e != nil {
				failedtables = append(failedtables, v)
				reasons = append(reasons, e.Error())
				continue
			}
		} else {
			sql, qk, e = mysqlx.ToClickhouseTable(req.Dsn, req.Database, v, "")
			if e != nil {
				failedtables = append(failedtables, v)
				reasons = append(reasons, e.Error())
				continue
			}
		}
		sqls = append(sqls, strings.Join(sql, ";\n"))
		qks = append(qks, qk)
	}
	return &types.GenerateCreateSqlResponse{
		Sql:           sqls,
		QueryKey:      qks,
		FailedReasons: reasons,
		FailedTables:  failedtables,
	}, nil
}
