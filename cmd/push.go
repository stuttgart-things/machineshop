/*
Copyright © 2024 Patrick Hermann patrick.hermann@sva.de
*/

package cmd

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	ipservice "github.com/stuttgart-things/clusterbook/ipservice"

	"github.com/stuttgart-things/machineshop/internal"
	sthingsCli "github.com/stuttgart-things/sthingsCli"
	"google.golang.org/grpc"

	"github.com/spf13/cobra"
)

type Message struct {
	Title           string `json:"title,omitempty"`           // if empty: info
	Message         string `json:"info,omitempty"`            // if empty: title
	Severity        string `json:"severity,omitempty"`        // default: info
	Author          string `json:"author,omitempty"`          // default: unknown
	Timestamp       string `json:"timestamp,omitempty"`       // generate timestamp func
	System          string `json:"system,omitempty"`          // default: unknown
	Tags            string `json:"tags,omitempty"`            // empty
	AssigneeAddress string `json:"assigneeaddress,omitempty"` // empty
	AssigneeName    string `json:"assigneename,omitempty"`    // empty
	Artifacts       string `json:"artifacts,omitempty"`       // empty
	Url             string `json:"url,omitempty"`             // empty
}

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
				messageBody := Message{
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

				rendered := RenderBody(homeRunBodyData, messageBody)
				fmt.Println(rendered)

				// CREATE HTTP-Request
				req, err := http.NewRequest("POST", destination, bytes.NewBuffer([]byte(rendered)))
				if err != nil {
					fmt.Println("faiulure at creating requests:", err)
					return
				}

				// ADD HEADER
				req.Header.Set("Content-Type", contentType)
				req.Header.Set("X-Auth-Token", token)

				// CREATE HTTP-Client + SEND REQUEST
				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					fmt.Println("error at sending request:", err)
					return
				}
				defer resp.Body.Close()

				// READ THE ANSWER
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println("error reading answer:", err)
					return
				}

				fmt.Println("answer status:", resp.Status)
				fmt.Println("answer body:", string(body))

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
}

func RenderBody(templateData string, object interface{}) string {

	tmpl, err := template.New("template").Parse(templateData)
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, object)

	if err != nil {
		fmt.Println(err)
	}

	return buf.String()

}
