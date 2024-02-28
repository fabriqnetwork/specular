pipeline {
    agent any
    environment{
        registry = "792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform"
    }

    stages {
        stage('prepare workspace') {
            steps {

                // checkout git
                checkout scmGit(
                    userRemoteConfigs: [
                        [ credentialsId: 'jenkins-specular', url: 'github.com:SpecularL2/specular.git']
                    ],
                    branches: [[name: '*/PR-*']]
                )

                // submodules
                sh "git submodule update --init --recursive"

                // make our workspace dir
                script {
                    if(!fileExists("workspace")) {
                        fileOperations([folderCreateOperation('workspace')])

                    }
                }


                // env files
                fileOperations([fileCopyOperation(
                        excludes: '',
                        flattenFiles: false,
                        includes: 'config/local_docker/[.?]',
                        targetLocation: "workspace"
                )])

            }
        }
        // stage('create build image') {
        //     steps{
        //         script {
        //             docker.withRegistry('https://792926601177.dkr.ecr.us-east-2.amazonaws.com', 'ecr:us-east-2:builder') {
        //                 docker.build(
        //                     registry + ":e2e-pr-10",
        //                     "-f docker/e2e.Dockerfile ."
        //                 )
        //             }
        //         }
        //     }
        // }
        stage('e2e-test') {
            parallel {
                stage('transactions') {
                    steps {
                      script {
                        docker.image("792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-10").inside("-w /specular/workspace") {
                          c -> sh "cd /specular/workspace && ls -la && ../sbin/run_e2e_tests.sh transactions"
                        }
                      }

                    }
                }
                stage('deposit') {
                    steps {
                      script {
                        docker.image("792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-10").inside("-w /specular/workspace") {
                          c -> sh "cd /specular/workspace && ../sbin/run_e2e_tests.sh deposit"
                        }
                      }
                    }
                }
                stage('erc20') {
                    steps {
                      script {
                        docker.image("792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-10").inside("-w /specular/workspace") {
                          c -> sh "cd /specular/workspace && ../sbin/run_e2e_tests.sh erc20"
                        }
                      }
                    }
                }
                stage('withdraw') {
                    steps {
                      script {
                        docker.image("792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-10").inside("-w /specular/workspace") {
                          c -> sh "cd /specular/workspace && ../sbin/run_e2e_tests.sh withdraw"
                        }
                      }
                    }
                }
            }
        }

    }
}
