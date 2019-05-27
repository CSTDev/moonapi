pipeline {
    agent { 
            label 'golang-builder' 
        }

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                make build
            }
        }
        stage('Test') {
            container('golang') {
                steps {
                    echo 'Testing..'
                    GO111MODULE=on; go test ./...
                }
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
            }
        }
    }
}