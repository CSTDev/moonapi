pipeline {
    agent any

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh "make build"
            }
        }
        stage('Test') {
            container('golang') {
                steps {
                    echo 'Testing..'
                    sh "GO111MODULE=on; go test ./..."
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