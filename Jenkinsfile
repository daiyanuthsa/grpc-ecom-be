pipeline {
    agent none

    stages {

        
        stage('Provision and Build on GCP') {
            agent { label 'built-in' }

            environment {
                TERRAFORM_DIR = 'terraform/gcp_builder'
            }
            options {
                // shallow clone untuk menghemat ruang di Jenkins Master
                skipDefaultCheckout true 
            }

            stages {
                stage('Checkout Terraform Scripts'){
                    steps{
                        cleanWs()

                        // Lakukan checkout eksplisit yang ringan di sini
                        checkout([
                            $class: 'GitSCM',
                            branches: [[name: '*/main']],
                            userRemoteConfigs: [[url: 'https://github.com/daiyanuthsa/grpc-ecom-be.git']],
                            extensions: [
                                [$class: 'CloneOption', shallow: true, noTags: true, depth: 1],
                                [$class: 'SubmoduleOption', disable: true]
                            ]
                        ])
                    }
                }
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
                                sh 'terraform apply -auto-approve -var="gcp_project_id=myexperiment-project"'
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

                                // Perintah build di dalam VM sekarang lebih sederhana
                                sh """
                                    ssh -o StrictHostKeyChecking=no -i \$SSH_KEY \$SSH_USER@${vmIp} '''
                                        echo \$DOCKER_PASS | docker login -u \$DOCKER_USER --password-stdin
                                        
                                        git clone https://github.com/daiyanuthsa/grpc-ecom-be.git
                                        cd grpc-ecom-be

                                        # Build & Push gRPC Image (tanpa buildx)
                                        docker build -t ${DOCKER_REGISTRY_USER}/grpc-backend:latest -f grpc.Dockerfile .
                                        docker push ${DOCKER_REGISTRY_USER}/grpc-backend:latest
                                        
                                        # Build & Push REST Image (tanpa buildx)
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
                        steps {
                            echo 'Tearing down the infrastructure...'
                            dir(TERRAFORM_DIR) {
                                withCredentials([file(credentialsId: 'gcp-service-account-key', variable: 'GOOGLE_APPLICATION_CREDENTIALS')]) {
                                    sh 'terraform destroy -auto-approve -var="gcp_project_id=nama-proyek-gcp-anda"'
                                }
                            }
                        }
                }
            }
        }

        stage('Deploy Services on Local Server') {
            agent { label 'built-in' }
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
            // Selalu bersihkan workspace di akhir untuk menjaga kebersihan
            cleanWs()
        }
    }
}