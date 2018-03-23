pipeline {
    // requires docker plugin and docker installed on node
    agent { docker { image 'golang' } }

    environment {
        workDir = "/go/src/github.com/individuwill/mcast"
        GOCACHE = "${env.WORKSPACE}/.cache"
    }

    stages {
        stage('Prepare') {
            steps {
                sh 'printenv'
                sh 'go version'
                sh 'mkdir -p ${workDir}'
                sh 'rm -rf ${workDir}'
                sh 'ln -s ${WORKSPACE} ${workDir}'
                git 'https://github.com/individuwill/mcast.git'
            }
        }

        stage('Build') {
            steps {
                sh 'go build'
                sh './build.sh'
            }
        }
        
        stage('Package') {
            steps {
                // requires pipeline-utility-steps plugin
                zip zipFile: 'binaries.zip', archive: true, dir: 'binaries'
            }
        }
    }
}