/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package internal

import (
	"crypto/tls"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func GetCredentials(secureConnection string) grpc.DialOption {

	switch secureConnection {

	case "true":
		log.Println("Using secure gRPC connection")
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true, // Adjust based on your security requirements
		}
		return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig))
	case "false":
		log.Println("Using insecure gRPC connection")
		return grpc.WithTransportCredentials(insecure.NewCredentials())
	default:
		log.Fatalf("Invalid SECURE_CONNECTION value: %s. Expected 'true' or 'false'", secureConnection)
		return nil // This will never be reached since log.Fatalf exits the program
	}
}
