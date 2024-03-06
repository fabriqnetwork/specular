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
        stage('create build image for pr') {
            when {
                not {
                    branch 'develop'
                }
            }
            steps{
                script {
                    docker.withRegistry('https://792926601177.dkr.ecr.us-east-2.amazonaws.com', 'ecr:us-east-2:builder') {
                        docker.build(
                            registry + ":e2e-pr-$GIT_COMMIT",
                            "-f docker/e2e.Dockerfile ."
                        )
                    }

                }
            }
        }
        stage('create build image for devnet') {
            when {
              branch 'develop'
            }
            steps{
                script {
                    docker.withRegistry('https://792926601177.dkr.ecr.us-east-2.amazonaws.com', 'ecr:us-east-2:builder') {
                        docker.build(
                            registry + ":e2e-$GIT_COMMIT",
                            "-f docker/e2e.Dockerfile ."
                        ).tag("e2e-latest").push()
                    }

                }
            }
        }
        stage('e2e-test') {
            when {
                not {
                    branch 'develop'
                }
            }
            parallel {
                stage('transactions') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$GIT_COMMIT ../sbin/run_e2e_tests.sh transactions"
                    }
                }
                stage('deposit') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$GIT_COMMIT ../sbin/run_e2e_tests.sh deposit"
                    }
                }
                stage('erc20') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$GIT_COMMIT ../sbin/run_e2e_tests.sh erc20"
                    }
                }
                stage('withdraw') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-pr-$GIT_COMMIT ../sbin/run_e2e_tests.sh withdraw"
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
                        docker.build(registry + ":$GIT_COMMIT", "-f docker/specular.Dockerfile .").push()
                    }
                }
            }
        }
        stage('upgrade helm') {
            when {
              branch "develop"
            }
          steps {
            withCredentials([[
                $class: 'AmazonWebServicesCredentialsBinding',
                credentialsId: "builder",
                accessKeyVariable: 'AWS_ACCESS_KEY_ID',
                secretKeyVariable: 'AWS_SECRET_ACCESS_KEY'
            ]]) {
              sh 'echo $AWS_ACCESS_KEY_ID'
              cd "charts/specular"
              sh "aws eks update-kubeconfig --name specular-staging-eks"
              sh "helm upgrade specular . -n specular --set image.tag=$GIT_COMMIT"
            }
          }
        }
    }
}
