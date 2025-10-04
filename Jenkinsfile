pipeline {
    agent any
    environment {
        DOCKER_REGISTRY_USER = 'daiyanuthsa'
        DOCKER_REGISTRY_CREDENTIALS_ID = 'dockerhub-credentials'
        DEPLOY_SERVER_CREDENTIALS_ID = 'deploy-server-credentials' // ID untuk kredensial SSH password
        DEPLOY_SERVER_IP = '192.168.100.3'
        DEPLOY_SERVER_USER = 'root'
    }

    stages {
        stage('Build gRPC Image') {
            steps {
                script {
                    def imageName = "${DOCKER_REGISTRY_USER}/grpc-backend:latest"
                    withCredentials([usernamePassword(credentialsId: DOCKER_REGISTRY_CREDENTIALS_ID, usernameVariable: 'USER', passwordVariable: 'PASS')]) {
                        sh "echo ${PASS} | docker login -u ${USER} --password-stdin"
                    }
                    sh "docker build -t ${imageName} -f grpc.Dockerfile ."
                    sh "docker push ${imageName}"
                }
            }
        }
        
        // Stage ini akan berjalan kedua, setelah gRPC selesai
        stage('Build REST Image') {
            steps {
                script {
                    def imageName = "${DOCKER_REGISTRY_USER}/rest-uploader:latest"
                    // Tidak perlu login lagi jika sudah di stage sebelumnya
                    sh "docker build -t ${imageName} -f rest.Dockerfile ."
                    sh "docker push ${imageName}"
                }
            }
        }
        stage('Deploy Backend Services') {
            steps {
                withCredentials([usernamePassword(credentialsId: DEPLOY_SERVER_CREDENTIALS_ID, usernameVariable: 'SSH_USER', passwordVariable: 'SSH_PASS')]) {
                    sh "sshpass -p '${SSH_PASS}' ssh -o StrictHostKeyChecking=no ${SSH_USER}@${DEPLOY_SERVER_IP} 'cd /home/root/my-app && docker-compose pull backend rest-uploader && docker-compose up -d --no-deps backend rest-uploader && docker system prune -af'"
                }
            }
        }
    }
    post {
        always {
            sh "docker logout"
        }
    }
}