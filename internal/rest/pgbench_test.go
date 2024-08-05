package rest

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestPgBench(t *testing.T) {
	fmt.Println("Starting Docker containers...")
	if err := runDockerCompose("up", "-d"); err != nil {
		t.Fatalf("Error starting Docker containers: %v", err)
	}

	defer func() {
		fmt.Println("Cleaning up Docker containers...")
		if err := runDockerCompose("down", "-v"); err != nil {
			fmt.Printf("Error cleaning up Docker containers: %v", err)
		}
	}()

	// Give services some time to start up
	fmt.Println("Waiting for services to start...")
	time.Sleep(10 * time.Second)

	if err := runPgBenchMigrations(); err != nil {
		t.Fatalf("Error running migrations: %v", err)
	}

	// Initialize the pgbench database
	fmt.Println("Initializing pgbench database...")
	if err := runPgBenchCommand("-i", "-s", "10", "-h", "localhost", "-U", "testuser", "-d", "testdb"); err != nil {
		t.Fatalf("Error initializing pgbench database: %v", err)
	}

	// Run the pgbench benchmark
	fmt.Println("Running pgbench benchmark...")
	if err := runPgBenchCommand("-c", "10", "-j", "2", "-T", "60", "-h", "localhost", "-U", "testuser", "-d", "testdb"); err != nil {
		t.Fatalf("Error running pgbench benchmark: %v", err)
	}
}

func runDockerCompose(args ...string) error {
	fmt.Printf("Running docker-compose command: %s\n", strings.Join(args, " "))
	cmd := exec.Command("docker-compose", append([]string{"-f", "../../docker-compose.test.yml"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runPgBenchCommand(args ...string) error {
	// Use docker-compose exec to run pgbench inside the db container
	fmt.Printf("Running pgbench command: %s\n", strings.Join(args, " "))
	cmd := exec.Command("docker-compose", append([]string{"exec", "-T", "db", "pgbench"}, args...)...)
	cmd.Env = append(os.Environ(), "PGPASSWORD=testpass")

	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	if err := cmd.Run(); err != nil {
		fmt.Printf("pgbench error output: %s\n", errBuf.String())
		return fmt.Errorf("pgbench command failed: %w", err)
	}

	fmt.Println("pgbench standard output:")
	fmt.Println(outBuf.String())

	return nil
}

func runPgBenchMigrations() error {
	// Run migrations using docker-compose exec
	fmt.Println("Running database migrations...")
	cmd := exec.Command("migrate", "-path", "../database/migrations", "-database", "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
