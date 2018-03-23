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
                sh 'mkdir -p testOutput'
                sh 'rm -rf testOutput/*'
                sh 'mkdir -p ${workDir}'
                sh 'rm -rf ${workDir}'
                sh 'ln -s ${WORKSPACE} ${workDir}'
                //git 'https://github.com/individuwill/mcast.git'
            }
        }

        stage('Build') {
            steps {
                sh 'go build'
            }
        }
        
        stage('Multicast Code Test') {
            steps {
                sh 'go test -v github.com/individuwill/mcast/multicast 2>&1 | tee testOutput/multicast.gotest'
            }
        }

        stage('CLI Code Test') {
            steps {
                sh 'go test -v github.com/individuwill/mcast 2>&1 | tee testOutput/cli.gotest'
            }
        }
 
        stage('Package') {
            steps {
                sh './build.sh'
                sh 'rm -f binaries.zip'
                // requires pipeline-utility-steps plugin
                zip zipFile: 'binaries.zip', archive: true, dir: 'binaries'
            }
        }
    }
}
