package sql

import (
	"log"
	"os/exec"
)

func tryRun(cmd string, args ...string) bool {
	if exec.Command(cmd, args...).Run() != nil {
		return false
	}
	return true
}

func hasPython() bool {
	return tryRun("python", "-V")
}

func hasTensorFlow() bool {
	return tryRun("python", "-c", "import tensorflow")
}

func hasMySQLConnector() bool {
	return tryRun("python", "-c", "import mysql.connector")
}

func hasDocker() bool {
	return tryRun("docker", "version")
}

func hasDockerImage(image string) bool {
	b, e := exec.Command("docker", "images", "-q", image).Output()
	if e != nil || len(b) == 0 {
		return false
	}
	return true
}

func tensorflowCmd() (cmd *exec.Cmd) {
	if hasPython() && hasTensorFlow() && hasMySQLConnector() {
		log.Printf("tensorflowCmd: run locally")
		cmd = exec.Command("python")
	} else if hasDocker() {
		log.Printf("tensorflowCmd: run in Docker container")
		const tfImg = "sqlflow/sqlflow"
		if !hasDockerImage(tfImg) {
			log.Printf("No local Docker image %s.  It will take a long time to pull.", tfImg)
		}
		cmd = exec.Command("docker", "run", "--rm", "--network=host", "-i", tfImg, "python")
	} else {
		log.Fatalf("No local TensorFlow or Docker.  No way to run TensorFlow programs")
	}
	return cmd
}
