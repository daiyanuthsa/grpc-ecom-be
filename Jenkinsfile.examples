pipeline {
    // Tidak ada agent global, kita tentukan per-stage
    agent none

    stages {
        // TAHAP INI DIJALANKAN DI VM TENCENT
        stage('Build & Push Images on Remote Agent') {
            agent { label 'tencent-vm' } // Pastikan label ini sesuai

            stages {
                stage('Login to Docker') {
                    steps {
                        withCredentials([
                            usernamePassword(credentialsId: 'dockerhub-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')
                        ]) {
                            sh 'echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin'
                        }
                    }
                }
                stage('Build and Push gRPC Image') {
                    steps {
                        withCredentials([string(credentialsId: 'dockerhub-username', variable: 'DOCKER_REGISTRY_USER')]) {
                            script {
                                def imageName = "${DOCKER_REGISTRY_USER}/grpc-backend:latest"
                                // Gunakan buildx untuk target arm64
                                sh "docker buildx build --platform linux/arm64 -t ${imageName} -f grpc.Dockerfile --push ."
                            }
                        }
                    }
                }
                stage('Build and Push REST Image') {
                    steps {
                        withCredentials([string(credentialsId: 'dockerhub-username', variable: 'DOCKER_REGISTRY_USER')]) {
                            script {
                                def imageName = "${DOCKER_REGISTRY_USER}/rest-uploader:latest"
                                // Gunakan buildx untuk target arm64
                                sh "docker buildx build --platform linux/arm64 -t ${imageName} -f rest.Dockerfile --push ."
                            }
                        }
                    }
                }
            }

            post {
                always {
                    sh 'docker logout'
                }
            }
        }

        // TAHAP INI DIJALANKAN DI JENKINS MASTER (STB ANDA)
        stage('Deploy Services on Local Server') {
            agent { label 'built-in' } // Menggunakan label yang benar untuk Master
            steps {
                withCredentials([
                    usernamePassword(credentialsId: 'deploy-server-credentials', usernameVariable: 'SSH_USER', passwordVariable: 'SSH_PASS'),
                    string(credentialsId: 'deploy-server-ip', variable: 'DEPLOY_SERVER_IP')
                ]) {
                    sh 'sshpass -p $SSH_PASS ssh -o StrictHostKeyChecking=no $SSH_USER@$DEPLOY_SERVER_IP \'cd /home/root/grpc-ecom && docker compose pull backend rest-uploader && docker compose up -d --no-deps backend rest-uploader && docker system prune -af\''
                }
            }
        }
    }
}