pipeline {
    agent { label 'built-in' } // Jalankan semua orkestrasi di master


    options {
        skipDefaultCheckout()  // ⬅️ Matikan auto-checkout bawaan
    }

    environment {
        TERRAFORM_DIR = 'terraform/gcp_builder'
        GCP_PROJECT_ID = 'myexperiment-project' 
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
             steps {
                withCredentials([
                    file(credentialsId: 'gcp-service-account-key', variable: 'GOOGLE_APPLICATION_CREDENTIALS'),
                    string(credentialsId: 'gcp-ssh-public-key', variable: 'GCP_SSH_PUBLIC_KEY')
                ]) {
                    // Initialize Terraform
                    dir(TERRAFORM_DIR) {
                        sh 'terraform init'
                    }

                    // Apply Terraform
                    dir(TERRAFORM_DIR) {
                        sh "terraform apply -auto-approve -var=\"gcp_project_id=${GCP_PROJECT_ID}\" -var=\"ssh_public_key_content=${GCP_SSH_PUBLIC_KEY}\""
                    }

                    // Build on Dynamic VM
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
            post {
                always {
                        echo 'Tearing down the infrastructure...'
                        dir(TERRAFORM_DIR) {
                            withCredentials([file(credentialsId: 'gcp-service-account-key', variable: 'GOOGLE_APPLICATION_CREDENTIALS')]) {
                                sh "terraform destroy -auto-approve -var=\"gcp_project_id=${GCP_PROJECT_ID}\""
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
        success {
            echo 'Destroying resources after successful build...'
            dir(TERRAFORM_DIR) {
                withCredentials([
                    file(credentialsId: 'gcp-service-account-key', variable: 'GOOGLE_APPLICATION_CREDENTIALS'),
                    string(credentialsId: 'gcp-ssh-public-key', variable: 'GCP_SSH_PUBLIC_KEY')
                    ]) {
                    sh '''
                        terraform init -reconfigure
                        terraform destroy -auto-approve -var="gcp_project_id=${GCP_PROJECT_ID}" -var="ssh_public_key_content=${GCP_SSH_PUBLIC_KEY}"
                    '''
                }
            }
        }
        failure {
            echo 'Skipping destroy because build failed before VM was provisioned.'
        }
    }

}