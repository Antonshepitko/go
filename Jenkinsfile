pipeline {
    agent any
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub-credentials')
        DOCKER_IMAGE = "antonshepitko/health-service"
        KUBE_CONFIG = credentials('kube-config')
    }
    stages {
        stage('Checkout') {
            steps {
                git url: 'https://github.com/Antonshepitko/go', branch: 'main'
            }
        }
        stage('Build Docker Image') {
            steps {
                sh 'docker build -t ${DOCKER_IMAGE}:${BUILD_NUMBER} .'
                sh 'docker tag ${DOCKER_IMAGE}:${BUILD_NUMBER} ${DOCKER_IMAGE}:latest'
            }
        }
        stage('Push to Docker Hub') {
            steps {
                sh 'echo $DOCKERHUB_CREDENTIALS_PSW | docker login -u $DOCKERHUB_CREDENTIALS_USR --password-stdin'
                sh 'docker push ${DOCKER_IMAGE}:${BUILD_NUMBER}'
                sh 'docker push ${DOCKER_IMAGE}:latest'
            }
        }
        stage('Deploy to Kubernetes') {
            steps {
                sh 'kubectl --kubeconfig=$KUBE_CONFIG apply -f health-service.yaml'
            }
        }
    }
    post {
        always {
            sh 'docker logout'
        }
    }
}