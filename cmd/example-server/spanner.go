package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/spanner"
	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	"cloud.google.com/go/spanner/spansql"
	"go.einride.tech/iam/iamexample"
	"go.einride.tech/iam/iamspanner"
	"google.golang.org/api/option"
	databasepb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"google.golang.org/grpc"
)

func newSpannerDatabase(ctx context.Context, conn *grpc.ClientConn) (*spanner.Client, error) {
	instanceAdminClient, err := instance.NewInstanceAdminClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	const projectID = "example"
	instanceID := fmt.Sprintf("example-instance-%d", time.Now().UnixNano())
	createInstanceOp, err := instanceAdminClient.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectID),
		InstanceId: instanceID,
		Instance: &instancepb.Instance{
			DisplayName: "Example",
			NodeCount:   1,
		},
	})
	if err != nil {
		return nil, err
	}
	createdInstance, err := createInstanceOp.Wait(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("created", createdInstance)
	databaseAdminClient, err := database.NewDatabaseAdminClient(ctx, option.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	databaseID := fmt.Sprintf("db%d", time.Now().UnixNano())
	var statements []string
	for _, schema := range []string{iamspanner.SQLSchema(), iamexample.SQLSchema()} {
		ddl, err := spansql.ParseDDL("example", schema)
		if err != nil {
			return nil, err
		}
		for _, ddlStmt := range ddl.List {
			statements = append(statements, ddlStmt.SQL())
		}
	}
	createDatabaseOp, err := databaseAdminClient.CreateDatabase(ctx, &databasepb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID),
		CreateStatement: fmt.Sprintf("CREATE DATABASE %s", databaseID),
		ExtraStatements: statements,
	})
	if err != nil {
		return nil, err
	}
	createdDatabase, err := createDatabaseOp.Wait(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("created", createdDatabase)
	return spanner.NewClient(ctx, createdDatabase.Name, option.WithGRPCConn(conn))
}
