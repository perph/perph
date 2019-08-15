package configure

import (
	context "context"
	fmt "fmt"
	"time"

	grpc "google.golang.org/grpc"
)

func ConfigureAgent(agentAddress string, jobID string) (string, error) {
	opts := grpc.WithInsecure()
	conn, err := grpc.Dial(agentAddress, opts)
	defer conn.Close()

	if err != nil {
		return "", err
	}

	client := NewPerphAgentClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	fmt.Println(client)

	configuredVersion, err := client.Configure(ctx, &NewVersion{VersionID: jobID})

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return configuredVersion.VersionID, nil
}
