pipeline {
    agent any

    environment {
        REGISTRY       = "docker.io/afriyanpm"
        IMAGE_NAME     = "Withdraw-Service-Test"
        IMAGE_TAG      = "${env.GIT_COMMIT.take(7)}"
        VPS_HOST       = "103.150.227.159"
        COMPOSE_PATH   = "/opt/Withdraw-Service-Test/docker-compose.yml"
    }

    stages {
        stage('Checkout') {
            steps {
                git branch: 'master', url: 'https://github.com/afryn123/Withdraw-Service-Test.git'
            }
        }

        stage('Build Docker Image') {
            steps {
                sh """
                docker build -t $REGISTRY/$IMAGE_NAME:$IMAGE_TAG .
                docker tag $REGISTRY/$IMAGE_NAME:$IMAGE_TAG $REGISTRY/$IMAGE_NAME:latest
                """
            }
        }

        stage('Push to Registry') {
            steps {
                withCredentials([usernamePassword(
                    credentialsId: 'docker-hub',
                    usernameVariable: 'DOCKER_USER',
                    passwordVariable: 'DOCKER_PASS'
                )]) {
                    sh """
                    echo $DOCKER_PASS | docker login -u $DOCKER_USER --password-stdin
                    docker push $REGISTRY/$IMAGE_NAME:$IMAGE_TAG
                    docker push $REGISTRY/$IMAGE_NAME:latest
                    docker logout
                    """
                }
            }
        }

        stage('Deploy to VPS via docker-compose') {
            steps {
                sshagent(['my-vps-ssh']) {
                    sh """
                    ssh -o StrictHostKeyChecking=no ubuntu@$VPS_HOST '
                        cd /opt/my-go-app &&
                        docker-compose pull &&
                        docker-compose up -d
                    '
                    """
                }
            }
        }
    }

    post {
        success {
            echo "✅ Deployed successfully using docker-compose"
        }
        failure {
            echo "❌ Deployment failed"
        }
        always {
            cleanWs()
        }
    }
}
