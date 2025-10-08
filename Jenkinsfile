pipeline {
    agent { label 'built-in' } // Jalankan semua orkestrasi di master


    options {
        skipDefaultCheckout()  // ⬅️ Matikan auto-checkout bawaan
    }

    environment {
        TERRAFORM_DIR = 'terraform/gcp_builder'
    }

    stages {
        stage('Clean Workspace') {
            steps {
                cleanWs()
            }
        }

        stage('Checkout Source') {
            steps {
                git(
                    url: 'https://github.com/daiyanuthsa/grpc-ecom-be.git',
                    branch: 'main',
                    changelog: false,
                    poll: false
                )
            }
        }

        stage('Provision and Build on GCP') {
            stages {
                stage('Initialize Terraform') {
                    steps {
                        dir(TERRAFORM_DIR) {
                            sh 'terraform init'
                        }
                    }
                }

                stage('Apply Terraform to Create VM') {
                    steps {
                        dir(TERRAFORM_DIR) {
                            withCredentials([file(credentialsId: 'gcp-service-account-key', variable: 'GOOGLE_APPLICATION_CREDENTIALS')]) {
                                sh 'terraform apply -auto-approve -var="gcp_project_id=nama-proyek-gcp-anda"'
                            }
                        }
                    }
                }

                stage('Build Images on Dynamic VM') {
                    steps {
                        script {
                            def vmIp = dir(TERRAFORM_DIR) {
                                sh(script: 'terraform output -raw instance_ip', returnStdout: true).trim()
                            }
                            
                            withCredentials([
                                usernamePassword(credentialsId: 'dockerhub-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS'),
                                string(credentialsId: 'dockerhub-username', variable: 'DOCKER_REGISTRY_USER'),
                                sshUserPrivateKey(credentialsId: 'gcp-ssh-key', keyFileVariable: 'SSH_KEY', usernameVariable: 'SSH_USER')
                            ]) {
                                sh "sleep 30"
                                sh """
                                    ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} '''
                                        echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin
                                        
                                        # Clone ulang di VM untuk memastikan kebersihan lingkungan build
                                        git clone --recurse-submodules https://github.com/daiyanuthsa/grpc-ecom-be.git
                                        cd grpc-ecom-be

                                        docker build -t ${DOCKER_REGISTRY_USER}/grpc-backend:latest -f grpc.Dockerfile .
                                        docker push ${DOCKER_REGISTRY_USER}/grpc-backend:latest
                                        
                                        docker build -t ${DOCKER_REGISTRY_USER}/rest-uploader:latest -f rest.Dockerfile .
                                        docker push ${DOCKER_REGISTRY_USER}/rest-uploader:latest

                                        docker logout
                                    '''
                                """
                            }
                        }
                    }
                }
            }
            post {
                always {
                        echo 'Tearing down the infrastructure...'
                        dir(TERRAFORM_DIR) {
                            withCredentials([file(credentialsId: 'gcp-service-account-key', variable: 'GOOGLE_APPLICATION_CREDENTIALS')]) {
                                sh 'terraform destroy -auto-approve -var="gcp_project_id=nama-proyek-gcp-anda"'
                            }
                    }
                }
            }
        }

        stage('Deploy Services on Local Server') {
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
    
    post {
        always {
            // Selalu bersihkan workspace di akhir
            cleanWs()
        }
    }
}