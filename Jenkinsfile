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
                    branches: [[name: '*/PR-*'], [name: '*/develop']],
                )
                // submodules
                sh "git submodule update --init --recursive"
                // make our workspace dir
                sh "rm -rf workspace && mkdir workspace"
                // env files
                sh 'cp -a config/local_docker/. workspace/'
                sh 'chmod -R 777 workspace'
            }
        }
        stage('create build image') {
            steps{
                script {
                    docker.withRegistry('https://792926601177.dkr.ecr.us-east-2.amazonaws.com', 'ecr:us-east-2:builder') {
                        docker.build(
                            registry + ":e2e-pr-$BUILD_NUMBER",
                            "-f docker/e2e.Dockerfile ."
                        )
                    }

                }
            }
        }
        stage('e2e-test') {
            parallel {
                stage('transactions') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$BUILD_NUMBER ../sbin/run_e2e_tests.sh transactions"
                    }
                }
                stage('deposit') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$BUILD_NUMBER ../sbin/run_e2e_tests.sh deposit"
                    }
                }
                stage('erc20') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$BUILD_NUMBER ../sbin/run_e2e_tests.sh erc20"
                    }
                }
                stage('withdraw') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$BUILD_NUMBER ../sbin/run_e2e_tests.sh withdraw"
                    }
                }
            }
        }
        stage('publish images') {
            when {
              branch "develop"
            }
            steps {
                script {
                    docker.withRegistry('https://792926601177.dkr.ecr.us-east-2.amazonaws.com', 'ecr:us-east-2:builder') {
                        docker.image(registry + ":e2e-pr-$BUILD_NUMBER").push()
                    }
                }
            }
        }
    }
}
