pipeline {
    agent {
        docker { image golang:1.16-alpine }
    }
    stages {
        stage('Test') {
            steps {
                sh 'go --version'
            }
        }
    }
}