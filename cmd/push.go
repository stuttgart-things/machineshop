/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"math/rand"

	homerun "github.com/stuttgart-things/homerun-library"
	sthingsBase "github.com/stuttgart-things/sthingsBase"

	ipservice "github.com/stuttgart-things/clusterbook/ipservice"

	"github.com/stuttgart-things/machineshop/internal"
	"github.com/stuttgart-things/machineshop/surveys"

	sthingsCli "github.com/stuttgart-things/sthingsCli"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

var (
	contentType = "application/json"
	// url           = "https://homerun.homerun-dev.sthings-vsphere.labul.sva.de/generic"
	minioLocation = "us-east-1"
	colors        = map[string]string{
		"red":    "#FF0000",
		"green":  "#00FF00",
		"blue":   "#0000FF",
		"orange": "#DF813D",
		"purple": "#726EAD",
		// Add more colors as needed
	}
	token = "IhrGeheimerToken"

	homeRunBodyData = `{
		"Title": "{{ .Title }}",
		"Message": "{{ .Message }}",
		"Severity": "{{ .Severity }}",
		"Author": "{{ .Author }}",
		"Timestamp": "{{ .Timestamp }}",
		"System": "{{ .System }}",
		"Tags": "{{ .Tags }}",
		"AssigneeAddress": "{{ .AssigneeAddress }}",
		"AssigneeName": "{{ .AssigneeName }}",
		"Artifacts": "{{ .Artifacts }}",
		"Url": "{{ .Url }}"
	}`
)

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "push things",
	Long:  `push things to external systems`,
	Run: func(cmd *cobra.Command, args []string) {

		// FLAGS
		target, _ := cmd.LocalFlags().GetString("target")
		source, _ := cmd.LocalFlags().GetString("source")
		color, _ := cmd.LocalFlags().GetString("color")
		title, _ := cmd.LocalFlags().GetString("title")
		body, _ := cmd.LocalFlags().GetString("message")
		author, _ := cmd.LocalFlags().GetString("author")
		severity, _ := cmd.LocalFlags().GetString("severity")
		system, _ := cmd.LocalFlags().GetString("system")
		destination, _ := cmd.LocalFlags().GetString("destination")
		tags, _ := cmd.LocalFlags().GetString("tags")
		assignee, _ := cmd.LocalFlags().GetString("assignee")
		assigneeUrl, _ := cmd.LocalFlags().GetString("assigneeUrl")
		artifacts, _ := cmd.LocalFlags().GetString("artifacts")
		url, _ := cmd.LocalFlags().GetString("url")
		insecure, _ := cmd.LocalFlags().GetBool("insecure")

		if destination != "" {

			switch target {

			case "ips":

				log.Info("⚡️ CONNECTING TO CLUSTERBOOK ⚡️")
				log.Info("CLUSTERBOOK SERVER: ", destination)
				log.Info("IPs: ", artifacts)
				log.Info("CLUSTER: ", assignee)

				clusterBookServer := destination //"clusterbook.rke2.sthings-vsphere.labul.sva.de:443"
				secureConnection := "true"

				// SELECT CREDENTIALS BASED ON SECURECONNECTION
				conn, err := grpc.NewClient(clusterBookServer, internal.GetCredentials(secureConnection))
				if err != nil {
					log.Fatalf("DID NOT CONNECT: %v", err)
				}
				defer conn.Close()

				c := ipservice.NewIpServiceClient(conn)

				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				clusterReq := &ipservice.ClusterRequest{
					IpAddressRange: artifacts,
					ClusterName:    assignee,
				}

				clusterRes, err := c.SetClusterInfo(ctx, clusterReq)
				if err != nil {
					log.Fatalf("COULD NOT SET CLUSTER INFO: %v", err)
				}

				log.Printf("CLUSTER STATUS: %s", clusterRes.Status)

			case "homerun-demo":
				fmt.Println("PUSHING TO HOMERUN-DEMO")

				// MAP FOR VALUES
				values := make(map[string]interface{})

				// READ PROFILE FOR MAY EXISTING PRE-SURVEY
				preSurveyValues := surveys.RunSurvey(source, "preSurvey")
				values = sthingsBase.MergeMaps(preSurveyValues, values)

				// READ PROFILE
				var demo surveys.HomerunDemo
				err := surveys.ReadProfileFile(source, &demo)
				if err != nil {
					log.Fatalf("Failed to read profile: %v", err)
				}

				// RENDER ALIASES + MERGE w/ VALUES
				aliases := surveys.RenderAliases(demo.Aliases, values)
				values = sthingsBase.MergeMaps(aliases, values)

				// LOOP OVER AUTHORS
				authorNames := []string{}
				authorAddresses := make(map[string]string)

				for _, author := range demo.Authors {
					authors := strings.Split(author, ":")
					authorNames = append(authorNames, authors[0])
					authorAddresses[authors[0]] = authors[1]
				}

				values["authors"] = authorNames
				values["authorAddresses"] = authorAddresses
				values["allUsecases"] = demo.Usecases
				values["messageTemplates"] = demo.MessageTemplates
				values["allObjects"] = demo.Objects
				values["allArtifacts"] = demo.Artifacts
				values["allUrls"] = demo.Urls

				// READ TEMPLATE TO RENDER
				funcMap := template.FuncMap{
					"hello": func(input string) string {
						return input + " Hello from push"
					},
					"random": func(input []string) string {
						return input[rand.Intn(len(input))]
					},
					"timestamp": func() string {
						dt := time.Now()
						return dt.Format("01-02-2006 15:04:05")
					},
					"getValueFromStringMap": func(key string, keyValues map[string]string) string {
						return keyValues[key]
					},
					"textBlock": func(severity string, verb string, template map[string]map[string]string) string {
						return template[severity][verb]
					},
					"randomUsecase": func(system string, usecases map[string][]string) string {
						systemUsecases := usecases[system]

						allUsecaseNames := []string{}

						// LOOP OVER USECASES
						for _, usecase := range systemUsecases {
							usecasesplit := strings.Split(usecase, ":")
							allUsecaseNames = append(allUsecaseNames, usecasesplit[0])
						}

						return allUsecaseNames[rand.Intn(len(allUsecaseNames))]
					},
					"getUsecaseVerb": func(usecaseName string, system string, usecases map[string][]string) string {
						verb := "Verb not found for usecase: " + usecaseName
						systemUsecases := usecases[system]

						// LOOP OVER USECASES
						for _, usecase := range systemUsecases {
							usecasesplit := strings.Split(usecase, ":")
							if usecasesplit[0] == usecaseName {
								verb = usecasesplit[1]
								break // Exit the loop once the verb is found
							}
						}

						return verb
					},
					"getObject": func(usecaseName string, objects map[string][]string) string {
						allObjects := objects[usecaseName]
						return allObjects[rand.Intn(len(allObjects))]
					},
					"getArtifact": func(usecaseName string, objects map[string][]string) string {
						allObjects := objects[usecaseName]
						return allObjects[rand.Intn(len(allObjects))]
					},
				}

				rendered, err := surveys.RenderTemplateInlineWithFunctions(funcMap, demo.BodyTemplate, values)

				if err != nil {
					log.Fatalf("FAILED TO RENDER TEMPLATE: %v", err)
				}

				fmt.Println(string(rendered))

				answer, resp := internal.SendToHomerun(destination, token, rendered, insecure)

				log.Info("ANSWER STATUS: ", resp.Status)
				log.Info("ANSWER BODY: ", string(answer))

			case "homerun":

				// GETTING TIMESTAMP
				dt := time.Now()
				timestamp := dt.Format("01-02-2006 15:04:05")

				if author == "" {
					author = "machineshop"
				}

				if system == "" {
					system = "machineshop"
				}

				if title == "" || body == "" {
					log.Error("TITILE AND/OR BODY EMPTY. EXITING")
					os.Exit(3)
				}

				log.Info("PUSHING TO HOMERUN")

				messageBody := homerun.Message{
					Title:           title,
					Message:         body,
					Severity:        severity,
					Author:          author,
					Timestamp:       timestamp,
					System:          system,
					Tags:            tags,
					AssigneeAddress: assigneeUrl,
					AssigneeName:    assignee,
					Artifacts:       artifacts,
					Url:             url,
				}

				rendered := homerun.RenderBody(homeRunBodyData, messageBody)
				fmt.Println(rendered)

				answer, resp := internal.SendToHomerun(destination, token, []byte(rendered), insecure)

				log.Info("ANSWER STATUS: ", resp.Status)
				log.Info("ANSWER BODY: ", string(answer))

			case "teams":

				if colors[color] == "" {
					target = "orange"
				}

				log.Info("PUSHING TO MS TEAMS")
				log.Info("MESSAGE: ", source)
				log.Info("MS TEAMS URL: ", destination)
				log.Info("COLOR: ", color)

				webhook := sthingsCli.MsTeamsWebhook{Title: "Machineshop", Text: source, Color: colors[color], Url: destination}
				sthingsCli.SendWebhookToTeams(webhook)

			case "minio":

				log.Info("PUSHING TO MINIO S3")
				log.Info("MINIO URL: ", os.Getenv("MINIO_ADDR"))
				log.Info("SOURCE: ", source)
				log.Info("TARGET: ", destination)

				// VERIFY IF SOURCE FILE IS EXISTING
				internal.ValidateSourceFile(source)

				clientCreated, minioClient := sthingsCli.CreateMinioClient()

				if !clientCreated {
					log.Error("MINIO CLIENT CAN NOT BE CREATED")
					os.Exit(3)

				} else {
					log.Info("MINIO CLIENT CREATED")

					destination := strings.Split(destination, ":")
					bucket := destination[0]
					objectName := destination[1]

					log.Info("BUCKET: ", bucket)
					log.Info("OBJECT: ", objectName)

					sthingsCli.CreateMinioBucket(minioClient, bucket, minioLocation)
					uploaded, fileSize := sthingsCli.UploadObjectToMinioBucket(minioClient, bucket, source, objectName)

					if uploaded {
						log.Info("SUCCESSFULLY UPLOADED OF SIZE: ", fileSize)
					}
				}

			}

		} else {
			log.Error("DESTINATION PATH SEEMS SO BE EMPTY")
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("source", "", "source file path")
	pushCmd.Flags().String("destination", "", "destination path")
	pushCmd.Flags().String("target", "minio", "push target")
	pushCmd.Flags().String("color", "orange", "color for webhook message")
	pushCmd.Flags().String("title", "", "title of homerun message")
	pushCmd.Flags().String("message", "", "homerun message body")
	pushCmd.Flags().String("severity", "info", "homerun message severity")
	pushCmd.Flags().String("author", "machineShop", "homerun message author")
	pushCmd.Flags().String("system", "", "homerun message system")
	pushCmd.Flags().String("tags", "", "homerun message tags")
	pushCmd.Flags().String("assignee", "", "homerun message assignee")
	pushCmd.Flags().String("assigneeUrl", "", "homerun message assignee url")
	pushCmd.Flags().String("artifacts", "", "homerun artifacts")
	pushCmd.Flags().String("url", "", "homerun message url/link")
	pushCmd.Flags().Bool("insecure", true, "insecure connection")
}

func Hello(input string) string {

	return input + " Hello from push"
}
