void setBuildStatus(String message, String state) {
  step([
      $class: "GitHubCommitStatusSetter",
      reposSource: [$class: "ManuallyEnteredRepositorySource", url: "https://github.com/individuwill/mcast"],
      contextSource: [$class: "ManuallyEnteredCommitContextSource", context: "ci/jenkins/build-status"],
      errorHandlers: [[$class: "ChangingBuildStatusErrorHandler", result: "UNSTABLE"]],
      statusResultSource: [ $class: "ConditionalStatusResultSource", results: [[$class: "AnyBuildResult", message: message, state: state]] ]
  ]);
}

pipeline {
    // requires docker plugin and docker installed on node
    agent {
        docker {
            image 'golang'
            // Will need to build and place into director the tool from:
            // https://github.com/jstemmer/go-junit-report
            args '-v /usr/local/bin/go-junit-report:/usr/local/bin/go-junit-report'
        }
    }

    environment {
        workDir = "/go/src/github.com/individuwill/mcast"
        GOCACHE = "${env.WORKSPACE}/.cache"
        testOutDir = 'testOutput'
    }

    stages {
        stage('Prepare') {
            steps {
                sh 'printenv'
                sh 'go version'
                sh 'mkdir -p ${testOutDir}'
                sh 'rm -rf ${testOutDir}/*'
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
        
        stage('Code Tests') {
            parallel {
                stage('Multicast Code Test') {
                    steps {
                        sh 'go test -v github.com/individuwill/mcast/multicast 2>&1 | tee ${testOutDir}/multicast.gotest'
                        sh 'cat ${testOutDir}/multicast.gotest | go-junit-report --package-name "multicast" > ${testOutDir}/multicast.xml'
                    }
                }

                stage('CLI Code Test') {
                    steps {
                        sh 'go test -v github.com/individuwill/mcast 2>&1 | tee ${testOutDir}/cli.gotest'
                        sh 'cat ${testOutDir}/cli.gotest | go-junit-report --package-name "cli" > ${testOutDir}/cli.xml'
                    }
                }
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

    post {
        always {
            junit 'testOutput/*.xml'
        }

        success {
            setBuildStatus('Build complete', 'SUCCESS')
        }

        failure {
            setBuildStatus('Build failure', 'FAILURE')
        }
    }
}
