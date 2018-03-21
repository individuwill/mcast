pipeline {
    agent { docker { image 'golang' } }
    stages {
        stage('build') {
            steps {
                sh 'go version'
                sh 'go build'
            }
        }
        stage('Test') {

        }
        stage('Deploy') {

        }
    }
    post {
        always {

        }
        success {

        }
        failure {

        }
        unstable {

        }
        changed {

        }
    }
}