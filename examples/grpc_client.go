// Exemple de client gRPC pour TwinPick
// Usage: go run examples/grpc_client.go
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/abroudoux/twinpick/api/proto"
)

func main() {
	// Connexion au serveur gRPC
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewTwinPickServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Test 1: Pick basique - juste les usernames, aucun filtre
	fmt.Println("=== Test 1: Pick basique (sans filtres) ===")
	pickResp, err := client.Pick(ctx, &pb.PickRequest{
		Usernames: []string{"abroudoux", "potatoze"},
	})
	if err != nil {
		log.Printf("Pick failed: %v", err)
	} else {
		fmt.Printf("Found %d films in common\n", pickResp.TotalCount)
		printFilms(pickResp.Films, 5)
	}

	// Test 2: Pick avec limite
	fmt.Println("\n=== Test 2: Pick avec limit=3 ===")
	pickResp, err = client.Pick(ctx, &pb.PickRequest{
		Usernames: []string{"abroudoux", "potatoze"},
		Filters: &pb.Filters{
			Limit: 3,
		},
	})
	if err != nil {
		log.Printf("Pick failed: %v", err)
	} else {
		fmt.Printf("Found %d films (limited to 3)\n", pickResp.TotalCount)
		printFilms(pickResp.Films, 10)
	}

	// Test 3: Pick avec filtre durée courte (<= 100 min)
	fmt.Println("\n=== Test 3: Pick films courts (duration=SHORT, <= 100 min) ===")
	pickResp, err = client.Pick(ctx, &pb.PickRequest{
		Usernames: []string{"abroudoux", "potatoze"},
		Filters: &pb.Filters{
			Duration: pb.Duration_DURATION_SHORT,
		},
	})
	if err != nil {
		log.Printf("Pick failed: %v", err)
	} else {
		fmt.Printf("Found %d short films in common\n", pickResp.TotalCount)
		printFilms(pickResp.Films, 5)
	}

	// Test 4: Pick avec filtre durée moyenne (<= 120 min) + limite
	fmt.Println("\n=== Test 4: Pick films moyens (duration=MEDIUM) + limit=5 ===")
	pickResp, err = client.Pick(ctx, &pb.PickRequest{
		Usernames: []string{"abroudoux", "potatoze"},
		Filters: &pb.Filters{
			Limit:    5,
			Duration: pb.Duration_DURATION_MEDIUM,
		},
	})
	if err != nil {
		log.Printf("Pick failed: %v", err)
	} else {
		fmt.Printf("Found %d medium-length films\n", pickResp.TotalCount)
		printFilms(pickResp.Films, 10)
	}

	// Test 5: Pick avec filtre genre
	fmt.Println("\n=== Test 5: Pick avec genre=horror ===")
	pickResp, err = client.Pick(ctx, &pb.PickRequest{
		Usernames: []string{"abroudoux", "potatoze"},
		ScrapperFilters: &pb.ScrapperFilters{
			Genres: []string{"horror"},
		},
	})
	if err != nil {
		log.Printf("Pick failed: %v", err)
	} else {
		fmt.Printf("Found %d horror films in common\n", pickResp.TotalCount)
		printFilms(pickResp.Films, 5)
	}

	// Test 6: Pick avec tous les filtres combinés
	fmt.Println("\n=== Test 6: Pick combiné (genre=thriller, duration=MEDIUM, limit=3) ===")
	pickResp, err = client.Pick(ctx, &pb.PickRequest{
		Usernames: []string{"abroudoux", "potatoze"},
		Filters: &pb.Filters{
			Limit:    3,
			Duration: pb.Duration_DURATION_MEDIUM,
		},
		ScrapperFilters: &pb.ScrapperFilters{
			Genres: []string{"thriller"},
		},
	})
	if err != nil {
		log.Printf("Pick failed: %v", err)
	} else {
		fmt.Printf("Found %d thriller films <= 120min\n", pickResp.TotalCount)
		printFilms(pickResp.Films, 10)
	}

	// Test 7: PickStream - Recevoir les films un par un (streaming)
	fmt.Println("\n=== Test 7: PickStream (Server Streaming) ===")
	stream, err := client.PickStream(ctx, &pb.PickRequest{
		Usernames: []string{"abroudoux", "potatoze"},
		Filters: &pb.Filters{
			Limit: 5,
		},
	})
	if err != nil {
		log.Printf("PickStream failed: %v", err)
	} else {
		fmt.Println("Receiving films via stream:")
		count := 0
		for {
			film, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Stream error: %v", err)
				break
			}
			count++
			fmt.Printf("  [%d] %s (%d) - %d min - %v\n",
				count, film.Title, film.Year, film.Duration, film.Directors)
		}
		fmt.Printf("Stream completed: received %d films\n", count)
	}
}

func printFilms(films []*pb.Film, max int) {
	for i, film := range films {
		if i >= max {
			fmt.Printf("  ... and %d more\n", len(films)-max)
			break
		}
		directors := "unknown"
		if len(film.Directors) > 0 {
			directors = film.Directors[0]
			if len(film.Directors) > 1 {
				directors += fmt.Sprintf(" (+%d)", len(film.Directors)-1)
			}
		}
		fmt.Printf("  - %s (%d) | %d min | %s\n",
			film.Title, film.Year, film.Duration, directors)
	}
}
