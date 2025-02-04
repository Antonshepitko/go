pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // ID в Jenkins для Docker Hub
        DOCKERHUB_REPO = "barongeddon/go" //Имя репозитория в Dockerhub
        REMOTE_SERVER = credentials('test-server-ip')  // Удаленный хост для разворачивания
        REMOTE_SSH_CREDENTIALS = credentials('test_ssh') // ID в Jenkins для SSH подключения
    }

    stages {
        stage('Checkout') {
            steps {
                echo 'Получаем код из GitHub...'
                // Укажите URL и ветку вашего репозитория
                git url: 'https://github.com/Antonshepitko/go.git', branch: 'master'
            }
        }
        stage('Build Docker Image') {
            steps {
                echo 'Собираем Docker-образ...'
                sh "docker build -t ${DOCKERHUB_REPO}:latest ."
            }
        }
        stage('Test Docker Container Locally') {
            steps {
                echo 'Запускаем временный контейнер для тестирования...'
                script {
                    // Запускаем контейнер в фоновом режиме
                    sh "docker run -d --name temp-container -p 8080:8080 ${DOCKERHUB_REPO}:latest"
                    // Ждем несколько секунд, чтобы контейнер запустился
                    sleep time: 5, unit: 'SECONDS'
                    // Выполняем тестовый запрос к healthcheck-эндпоинту
                    sh "curl --fail http://localhost:8080/health"
                    // Останавливаем и удаляем временный контейнер
                    sh "docker stop temp-container"
                    sh "docker rm temp-container"
                }
            }
        }
        stage('Push Image to Docker Hub') {
            steps {
                echo 'Публикуем образ в Docker Hub...'
                script {
                    sh """
                    echo ${DOCKERHUB_CREDENTIALS_PSW} | docker login -u ${DOCKERHUB_CREDENTIALS_USR} --password-stdin
                    """
                    // Отправляем образ в репозиторий
                    sh "docker push ${DOCKERHUB_REPO}:latest"
                }
            }
        }
        stage('Deploy on Remote Server') {
            steps {
                echo 'Деплоим образ на удалённом сервере...'
                sshagent (credentials: [REMOTE_SSH_CREDENTIALS]) {
                    script {
                        sh """
                        ssh -o StrictHostKeyChecking=no deployuser@${REMOTE_SERVER} 'docker stop my-go-service || true && docker rm my-go-service || true'
                        ssh -o StrictHostKeyChecking=no deployuser@${REMOTE_SERVER} 'docker pull ${DOCKERHUB_REPO}:latest'
                        ssh -o StrictHostKeyChecking=no deployuser@${REMOTE_SERVER} 'docker run -d --name my-go-service -p 8080:8080 ${DOCKERHUB_REPO}:latest'
                        """
                    }
                }
            }
        }
        stage('Test Remote Deployment') {
            steps {
                echo 'Проверяем доступность сервиса на удалённом сервере...'
                script {
                    // Ждем несколько секунд, чтобы контейнер успел запуститься
                    sleep time: 10, unit: 'SECONDS'
                    // Выполняем тестовый запрос к сервису на удалённом сервере
                    sh "curl --fail http://${REMOTE_SERVER}:8080/health"
                }
            }
        }
    }
    post {
        always {
            echo 'Пайплайн завершён'
        }
    }
}
