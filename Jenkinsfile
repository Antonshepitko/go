pipeline {
    agent any

    environment {
        DOCKERHUB_CREDENTIALS = credentials('dockerhub') // ID в Jenkins для Docker Hub
        DOCKERHUB_REPO = "barongeddon/go" //Имя репозитория в Dockerhub
        REMOTE_SERVER = credentials('test-server-ip')  // Удаленный хост для разворачивания
        REMOTE_SSH_CREDENTIALS = credentials('test_ssh') // ID в Jenkins для SSH подключения
        REMOTE_DEPLOY_DIR = "/home/deployuser/deploy/myapp"
        REMOTE_USER = "deployuser"
    }

    stages {

        stage('Checkout') {
            steps {
                echo 'Получаем код из GitHub...'
                // Клонирование репозитория. Обратите внимание: здесь используется публичный URL, без авторизации.
                git url: 'https://github.com/Antonshepitko/go.git', branch: 'master'
            }
        }

        stage('Deploy on Remote Server') {
            steps {
                // Блок sshagent использует заранее настроенные SSH-учётные данные в Jenkins.
                sshagent (credentials: ['test_ssh']) {
                    script {
                        // Формируем команду, которая будет выполнена на удалённом сервере.
                        // Команда проверяет: если директория существует, то обновляет код, иначе — клонирует репозиторий.
                        // Затем переходит в директорию, строит Docker-образ, останавливает старый контейнер (если есть) и запускает новый.
                        def remoteCmd = """
                            if [ -d '${REMOTE_DEPLOY_DIR}' ]; then
                                cd ${REMOTE_DEPLOY_DIR} && git pull;
                            else
                                git clone https://github.com/Antonshepitko/go.git ${REMOTE_DEPLOY_DIR};
                            fi;
                            cd ${REMOTE_DEPLOY_DIR} &&
                            docker build -t ${DOCKERHUB_REPO}:latest . &&
                            docker stop my-go-service || true &&
                            docker rm my-go-service || true &&
                            docker run -d --name my-go-service -p 8080:8080 ${DOCKERHUB_REPO}:latest
                        """
                        // Поскольку наш Jenkins работает на Windows, для выполнения SSH-команды используем bat.
                        // В этом случае команда ssh должна быть доступна в PATH (например, из Git for Windows).
                        bat "ssh -o StrictHostKeyChecking=no ${REMOTE_USER}@${REMOTE_SERVER} \"${remoteCmd}\""
                    }
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
