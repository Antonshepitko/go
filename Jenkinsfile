pipeline {
    agent any

    environment {
        // Учётные данные для Docker Hub (тип "Username with password")
        DOCKERHUB_CREDENTIALS = credentials('dockerhub')
        // Базовое имя репозитория Docker Hub (например, имя вашей учётной записи)
        DOCKERHUB_REPO = "barongeddon"
        // Имя сервиса, который будет развёрнут (жёстко задано)
        SERVICE_NAME = "my-go-service"
        // URL Git-репозитория для данного сервиса
        GIT_REPO = "https://github.com/Antonshepitko/go.git"
        // Тег Docker-образа
        DOCKER_IMAGE_TAG = "latest"
        // Полное имя Docker-образа, которое будет собираться и пушиться
        FULL_DOCKER_IMAGE = "${DOCKERHUB_REPO}/${SERVICE_NAME}:${DOCKER_IMAGE_TAG}"
        // Путь на удалённом сервере, куда будет развёрнут сервис
        REMOTE_DEPLOY_DIR = "/home/deployuser/deploy/${SERVICE_NAME}"
        // Имя пользователя для SSH-подключения к удалённому серверу
        REMOTE_USER = "deployuser"
        // Удалённый сервер (если это секрет, то используйте тип Secret text; здесь предполагается, что значение – IP или доменное имя)
        REMOTE_SERVER = credentials('test-server-ip')
    }

    stages {

        stage('Checkout') {
            steps {
                echo "Клонируем репозиторий: ${GIT_REPO}"
                // Клонирование репозитория из GitHub (ветка master)
                git url: "${GIT_REPO}", branch: "master"
            }
        }

        stage('Build Docker Image') {
            steps {
                echo "Собираем Docker-образ: ${FULL_DOCKER_IMAGE}"
                // Сборка Docker-образа (Dockerfile должен находиться в корне репозитория)
                sh "docker build -t ${FULL_DOCKER_IMAGE} ."
            }
        }

        stage('Test Docker Container Locally') {
            steps {
                script {
                    echo "Запускаем временный контейнер для тестирования..."
                    // Запускаем контейнер с пробросом порта 8081 (локально, внутри контейнера порт 8080)
                    sh "docker run -d --name temp_${SERVICE_NAME} -p 8081:8080 ${FULL_DOCKER_IMAGE}"
                    // Ждём несколько секунд, чтобы контейнер успел подняться
                    sh "sleep 10"
                    // Тестируем эндпоинт /health (ожидаем, что приложение вернёт, например, "ok")
                    sh "curl --fail http://localhost:8081/health"
                    // Останавливаем и удаляем тестовый контейнер
                    sh "docker stop temp_${SERVICE_NAME}"
                    sh "docker rm temp_${SERVICE_NAME}"
                }
            }
        }

        stage('Push Image to Docker Hub') {
            steps {
                script {
                    echo "Публикуем Docker-образ ${FULL_DOCKER_IMAGE} в Docker Hub..."
                    // Логинимся в Docker Hub и пушим образ.
                    // Учётные данные из DOCKERHUB_CREDENTIALS автоматически создают переменные:
                    // DOCKERHUB_CREDENTIALS_USR и DOCKERHUB_CREDENTIALS_PSW
                    sh """
                        docker login -u ${DOCKERHUB_CREDENTIALS_USR} -p ${DOCKERHUB_CREDENTIALS_PSW}
                        docker push ${FULL_DOCKER_IMAGE}
                    """
                }
            }
        }

        stage('Deploy on Remote Server') {
            steps {
                // Используем sshagent с учётными данными SSH (ID "test_ssh") для подключения к удалённому серверу
                sshagent (credentials: ['test_ssh']) {
                    script {
                        // Формируем команду для удалённого сервера:
                        // Здесь мы делаем pull образа из Docker Hub, останавливаем и удаляем старый контейнер (если есть)
                        // и запускаем новый контейнер с пробросом порта 8080.
                        def remoteCmd = """
                            docker pull ${FULL_DOCKER_IMAGE} &&
                            docker stop ${SERVICE_NAME} || true &&
                            docker rm ${SERVICE_NAME} || true &&
                            docker run -d --name ${SERVICE_NAME} -p 8080:8080 ${FULL_DOCKER_IMAGE}
                        """
                        echo "Разворачиваем на удалённом сервере командой: ${remoteCmd}"
                        // Выполняем команду по SSH на удалённом сервере.
                        sh "ssh -o StrictHostKeyChecking=no ${REMOTE_USER}@${REMOTE_SERVER} '${remoteCmd}'"
                        sh "exit"
                    }
                }
            }
        }

        stage('Test Remote Deployment') {
            steps {
                script {
                    echo "Проверяем доступность сервиса на удалённом сервере..."
                    // Ждём, чтобы контейнер успел запуститься
                    sh "sleep 10"
                    // Тестируем эндпоинт /health через SSH (запрос к localhost на удалённом сервере)
                    sh "curl -fail 147.45.60.20:8080/health"
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
