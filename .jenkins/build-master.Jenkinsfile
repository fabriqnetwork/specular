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
                    branches: [[name: '*/rhoop/*']]
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

                // login to ecr
                sh "aws ecr get-login-password --region us-east-2 | docker login --username AWS --password-stdin 792926601177.dkr.ecr.us-east-2.amazonaws.com"
            }
        }
        stage('create build image and push') {
            steps{
                script {
                    dockerImage = docker.build(
                        registry + ":build-v0.0.$BUILD_NUMBER",
                        "-f docker/build-master.Dockerfile ."
                    )
                    dockerImage.push()
                    dockerImage.push("latest")
                }
            }
        }
    }
}
