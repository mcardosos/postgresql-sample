package main

import (
	"fmt"
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/Azure/go-autorest/autorest/utils"
	"github.com/mcardosos/azure-sdk-for-go/arm/postgresql"
	"github.com/mcardosos/azure-sdk-for-go/arm/resources/resources"
)

// This example requires that the following environment vars are set:
//
// AZURE_TENANT_ID: contains your Azure Active Directory tenant ID or domain
// AZURE_CLIENT_ID: contains your Azure Active Directory Application Client ID
// AZURE_CLIENT_SECRET: contains your Azure Active Directory Application Secret
// AZURE_SUBSCRIPTION_ID: contains your Azure Subscription ID
//
var (
	groupClient  resources.GroupsClient
	serverClient postgresql.ServersClient
	groupName    = "postgresql-sample"
	location     = "westus"
	serverName   = "golangrocks"
)

func init() {
	authorizer, err := utils.GetAuthorizer(azure.PublicCloud)
	onErrorFail(err, "GetAuthorizer failed")

	subscriptionID := utils.GetEnvVarOrExit("AZURE_SUBSCRIPTION_ID")
	createClients(subscriptionID, authorizer)
}

func main() {
	fmt.Println("Create resource group...")
	resourceGroupParameters := resources.Group{
		Location: &location,
	}
	_, err := groupClient.CreateOrUpdate(groupName, resourceGroupParameters)
	onErrorFail(err, "CreateOrUpdate resource group failed")
	defer groupClient.Delete(groupName, nil)

	fmt.Println("Create server...")
	server := postgresql.ServerForCreate{
		Location: to.StringPtr("westus"),
		Properties: &postgresql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         to.StringPtr("notadmin"),
			AdministratorLoginPassword: to.StringPtr("Pa$$w0rd1975"),
			StorageMB:                  to.Int64Ptr(51200),
		},
	}
	_, errChan := serverClient.Create(groupName, serverName, server, nil)
	onErrorFail(<-errChan, "Create failed")
}

func createClients(subscriptionID string, authorizer *autorest.BearerAuthorizer) {
	groupClient = resources.NewGroupsClient(subscriptionID)
	groupClient.Authorizer = authorizer

	serverClient = postgresql.NewServersClient(subscriptionID)
	serverClient.Authorizer = authorizer
}

// onErrorFail prints a failure message and exits the program if err is not nil.
func onErrorFail(err error, message string, a ...interface{}) {
	if err != nil {
		fmt.Printf("%s: %s\n", fmt.Sprintf(message, a), err)
		os.Exit(1)
	}
}
