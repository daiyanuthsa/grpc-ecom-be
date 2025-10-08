pipeline {
    agent { label 'built-in' }

    options {
        skipDefaultCheckout()
    }

    environment {
        TERRAFORM_DIR = 'terraform/gcp_builder'
        GCP_PROJECT_ID = 'myexperiment-project'
        GOOGLE_APPLICATION_CREDENTIALS = credentials('gcp-service-account-key') // ✅ Kredensial global
    }

    stages {

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

        stage('Provision VM with Terraform') {
            steps {
                withCredentials([
                    string(credentialsId: 'gcp-ssh-public-key', variable: 'GCP_SSH_PUBLIC_KEY')
                ]) {
                    dir(TERRAFORM_DIR) {
                        sh '''
                            terraform init
                            terraform apply -auto-approve \
                                -var="gcp_project_id=${GCP_PROJECT_ID}" \
                                -var="ssh_public_key_content=${GCP_SSH_PUBLIC_KEY}"
                        '''
                    }
                }
            }
        }

        stage('Build & Push Docker Images on GCP VM') {
            steps {
                script {
                    // ✅ Tidak lagi error: GOOGLE_APPLICATION_CREDENTIALS sudah aktif global
                    def vmIp = dir(TERRAFORM_DIR) {
                        sh(script: 'terraform output -raw instance_ip', returnStdout: true).trim()
                    }

                    withCredentials([
                        usernamePassword(credentialsId: 'dockerhub-credentials', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS'),
                        // string(credentialsId: 'dockerhub-username', variable: 'DOCKER_USER'),
                        sshUserPrivateKey(credentialsId: 'gcp-ssh-key', keyFileVariable: 'SSH_KEY', usernameVariable: 'SSH_USER')
                    ]) {

                        stage('Waiting for SSH to become ready'){
                            sh """
                                echo '⏳ Waiting for SSH to become ready...'
                                for i in {1..15}; do
                                    if ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} 'echo SSH OK' 2>/dev/null; then
                                        echo '✅ SSH is ready!'
                                        break
                                    fi
                                    echo "SSH not ready yet (attempt \$i)..."
                                    sleep 10
                                done
                            """
                        }
                        // 🔹 Step 1: Docker Login
                        stage('Remote: Docker Login') {
                            sh """
                                ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} <<'EOF'
                                set -e
                                echo "\$DOCKER_PASS" | docker login -u "\$DOCKER_USER" --password-stdin
                                docker info | grep Username || true
                                EOF
                            """
                        }
                        // 🔹 Step 2: Clone repository
                        stage('Remote: Git Clone') {
                            sh """
                                ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} <<'EOF'
                                set -e
                                if [ ! -d grpc-ecom-be ]; then
                                    git clone --recurse-submodules https://github.com/daiyanuthsa/grpc-ecom-be.git
                                else
                                    echo "Repo already exists, pulling latest..."
                                    cd grpc-ecom-be && git pull
                                fi
                                EOF
                            """
                        }

                        // 🔹 Step 3: Build gRPC backend
                        stage('Remote: Build gRPC Backend') {
                            sh """
                                ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} <<'EOF'
                                set -e
                                cd grpc-ecom-be
                                docker build -t "\$DOCKER_USER/grpc-backend:latest" -f grpc.Dockerfile .
                                EOF
                            """
                        }

                        // 🔹 Step 4: Push gRPC backend
                        stage('Remote: Push gRPC Backend') {
                            sh """
                                ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} <<'EOF'
                                set -e
                                docker push "\$DOCKER_USER/grpc-backend:latest"
                                EOF
                            """
                        }

                        // 🔹 Step 5: Build REST uploader
                        stage('Remote: Build REST Uploader') {
                            sh """
                                ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} <<'EOF'
                                set -e
                                cd grpc-ecom-be
                                docker build -t "\$DOCKER_USER/rest-uploader:latest" -f rest.Dockerfile .
                                EOF
                            """
                        }

                        // 🔹 Step 6: Push REST uploader
                        stage('Remote: Push REST Uploader') {
                            sh """
                                ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} <<'EOF'
                                set -e
                                docker push "\$DOCKER_USER/rest-uploader:latest"
                                docker logout
                                EOF
                            """
                        }
                    }
                }
            }

            post {
                always {
                    echo 'Destroying VM...'
                    dir(TERRAFORM_DIR) {
                        withCredentials([
                            string(credentialsId: 'gcp-ssh-public-key', variable: 'GCP_SSH_PUBLIC_KEY')
                        ]) {
                            sh '''
                                terraform destroy -auto-approve \
                                    -var="gcp_project_id=${GCP_PROJECT_ID}" \
                                    -var="ssh_public_key_content=${GCP_SSH_PUBLIC_KEY}"
                            '''
                        }
                    }
                }
            }
        }

        stage('Deploy to Local Server') {
            steps {
                withCredentials([
                    usernamePassword(credentialsId: 'deploy-server-credentials', usernameVariable: 'SSH_USER', passwordVariable: 'SSH_PASS'),
                    string(credentialsId: 'deploy-server-ip', variable: 'DEPLOY_SERVER_IP')
                ]) {
                    sh '''
                        sshpass -p $SSH_PASS ssh -o StrictHostKeyChecking=no $SSH_USER@$DEPLOY_SERVER_IP '
                            cd /home/root/grpc-ecom &&
                            docker compose pull backend rest-uploader &&
                            docker compose up -d --no-deps backend rest-uploader &&
                            docker system prune -af
                        '
                    '''
                }
            }
        }
    }

    post {
        success {
            echo '✅ Build & deployment succeeded.'
        }
        failure {
            echo '❌ Build failed. Infrastructure already cleaned up.'
        }
        always {
            cleanWs()
        }
    }
}
