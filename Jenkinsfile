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
            sh 'echo testing...'
        }
        stage('Deploy') {
            sh 'echo Deploying...'
        }
    }
    post {
        always {
            echo 'always run after'
        }
        success {
            echo 'only run after on success'
        }
        failure {
            echo 'only run after on failure'
        }
        unstable {
            echo 'only run after if marked unstable'
        }
        changed {
            echo 'This will run only if the state of the Pipeline has changed'
            echo 'For example, if the Pipeline was previously failing but is now successful'
        }
    }
}