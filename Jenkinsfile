pipeline {
    agent any

    // Параметры сборки: имя сервиса, выбор репозитория, тег Docker-образа
    parameters {
        string(name: 'SERVICE_NAME', defaultValue: 'my-go-service', description: 'Имя сервиса, который будет развёрнут')
        choice(name: 'GIT_REPO', choices: '''https://github.com/Antonshepitko/go.git
https://github.com/Antonshepitko/another.git''', description: 'Выберите репозиторий для сборки')
        string(name: 'DOCKER_IMAGE_TAG', defaultValue: 'latest', description: 'Тег Docker-образа')
    }

    // Переменные окружения. Здесь вы задаёте базовое имя репозитория в Docker Hub,
    // формируете полное имя образа, а также пути и данные для удалённого деплоя.
    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub')
        DOCKERHUB_REPO = "barongeddon"
        FULL_DOCKER_IMAGE = "${DOCKERHUB_REPO}/${params.SERVICE_NAME}:${params.DOCKER_IMAGE_TAG}"
        REMOTE_DEPLOY_DIR = "/home/deployuser/deploy/${params.SERVICE_NAME}"
        REMOTE_USER = "deployuser"

        REMOTE_SERVER = credentials('test-server-ip')
    }

    stages {

        stage('Checkout') {
            steps {
                echo "${DOCKERHUB_CREDENTIALS_USR}"
                echo "Клонируем репозиторий: ${params.GIT_REPO}"
                // Клонируем выбранный репозиторий. Здесь используется публичный URL.
                git url: "${params.GIT_REPO}", branch: 'master'
            }
        }

        stage('Build Docker Image') {
            steps {
                echo "Собираем Docker-образ: ${env.FULL_DOCKER_IMAGE}"
                // Выполняем сборку Docker-образа. Dockerfile должен быть в корне репозитория.
                sh "docker build -t ${env.FULL_DOCKER_IMAGE} ."
            }
        }

        stage('Test Docker Container Locally') {
            steps {
                script {
                    echo "Запускаем временный контейнер для тестирования..."
                    // Запускаем контейнер локально, пробрасывая порт 8081 на 8080 внутри контейнера
                    sh "docker run -d --name temp_${params.SERVICE_NAME} -p 8081:8080 ${env.FULL_DOCKER_IMAGE}"
                    // Ждём несколько секунд для старта
                    sh "sleep 10"
                    // Тестируем эндпоинт /health. Если приложение отвечает, тест считается успешным.
                    sh "curl --fail http://localhost:8081/health"
                    // Останавливаем и удаляем тестовый контейнер
                    sh "docker stop temp_${params.SERVICE_NAME}"
                    sh "docker rm temp_${params.SERVICE_NAME}"
                }
            }
        }

        stage('Push Image to Docker Hub') {
            steps {
                script {
                    echo "Публикуем Docker-образ ${env.FULL_DOCKER_IMAGE} в Docker Hub..."
                    // Логинимся в Docker Hub. Учётные данные должны быть настроены в Jenkins с ID "dockerhub".
                    sh """
                      docker login -u ${DOCKERHUB_CREDENTIALS_USR} -p ${DOCKERHUB_CREDENTIALS_PSW}
                      docker push ${env.FULL_DOCKER_IMAGE}
                    """
                }
            }
        }

        stage('Deploy on Remote Server') {
            steps {
                // Используем sshagent с учётными данными SSH (ID "test_ssh"), чтобы подключиться к удалённому серверу.
                sshagent (credentials: ['test_ssh']) {
                    script {
                        // Формируем команду для удалённого сервера.
                        // Команда делает следующее:
                        // - Удаляет старую директорию (если есть) и клонирует свежий репозиторий.
                        // - Переходит в каталог и собирает Docker-образ (необязательно, если образ уже запушен в Docker Hub).
                        // - На самом деле, мы будем брать образ из Docker Hub.
                        // - Останавливает и удаляет старый контейнер, если он запущен.
                        // - Запускает новый контейнер с пробросом порта 8080.
                        def remoteCmd = """
                          docker pull ${env.FULL_DOCKER_IMAGE} &&
                          docker stop ${params.SERVICE_NAME} || true &&
                          docker rm ${params.SERVICE_NAME} || true &&
                          docker run -d --name ${params.SERVICE_NAME} -p 8080:8080 ${env.FULL_DOCKER_IMAGE}
                        """
                        echo "Разворачиваем на удалённом сервере командой: ${remoteCmd}"
                        // Выполняем команду по SSH на удалённом сервере.
                        sh "ssh -o StrictHostKeyChecking=no ${REMOTE_USER}@${REMOTE_SERVER} '${remoteCmd}'"
                    }
                }
            }
        }

        stage('Test Remote Deployment') {
            steps {
                script {
                    echo "Проверяем доступность сервиса на удалённом сервере..."
                    // Ждём несколько секунд, чтобы контейнер успел запуститься.
                    sh "sleep 10"
                    // Выполняем тестовый запрос к сервису на удалённом сервере через SSH.
                    sh "ssh -o StrictHostKeyChecking=no ${REMOTE_USER}@${REMOTE_SERVER} 'curl --fail http://localhost:8080/health'"
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
