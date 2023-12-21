pipeline {
    agent any
    stages {
        stage('E2E Tests') {
            parallel {
                stage('transactions') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest ../sbin/run_e2e_tests.sh transactions"
                    }
                }
                stage('deposit') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest ../sbin/run_e2e_tests.sh deposit"
                    }
                }
                stage('erc20') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest ../sbin/run_e2e_tests.sh erc20"
                    }
                }
                stage('withdraw') {
                    steps {
                        sh "docker run -w /specular/workspace 792926601177.dkr.ecr.us-east-2.amazonaws.com/specular-platform:e2e-latest ../sbin/run_e2e_tests.sh withdraw"
                    }
                }
            }
        }
    }
}
