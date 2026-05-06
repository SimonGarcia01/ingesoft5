pipeline {
    agent any
    environment {
        SONAR_HOST_URL = 'http://sonarqube:9000'
        SONAR_PROJECT_KEY = 'my-app'
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Build & Test') {
            steps {
                sh 'chmod +x mvnw'
                sh './mvnw -B clean package'
            }
        }

        stage('Static Analysis (SonarQube)') {
            steps {
                withCredentials([string(credentialsId: 'sonar-token', variable: 'SONAR_TOKEN')]) {
                    sh './mvnw -B sonar:sonar -Dsonar.projectKey=$SONAR_PROJECT_KEY -Dsonar.host.url=$SONAR_HOST_URL -Dsonar.login=$SONAR_TOKEN -Dsonar.qualitygate.wait=true -Dsonar.qualitygate.timeout=300'
                }
            }
        }

        stage('Docker Build') {
            steps {
                sh 'for i in $(seq 1 20); do docker version && break || sleep 2; done'
                sh 'docker build -t mi-app:latest .'
            }
        }

        stage('Container Security Scan (Trivy)') {
            steps {
                                sh '''
TRIVY_DIR=.trivy
mkdir -p "$TRIVY_DIR"
if [ ! -x "$TRIVY_DIR/trivy" ]; then
    curl -sSfL "https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh" | sh -s -- -b "$TRIVY_DIR"
fi
"$TRIVY_DIR/trivy" image --exit-code 0 --severity CRITICAL --no-progress mi-app:latest
'''
            }
        }

        stage('Deploy') {
            when { branch 'main' }
            steps {
                sh 'docker rm -f mi-app || true'
                sh 'docker run -d --name mi-app -p 80:80 -e SERVER_PORT=80 mi-app:latest'
            }
        }
    }

    post {
        always {
            script {
                def hasReports = sh(
                    script: 'ls -1 target/surefire-reports/*.xml >/dev/null 2>&1',
                    returnStatus: true
                ) == 0
                if (hasReports) {
                    junit 'target/surefire-reports/*.xml'
                }
            }
            archiveArtifacts artifacts: 'target/*.jar', fingerprint: true, allowEmptyArchive: true
            echo 'Limpiando entorno...'
            script {
                try {
                    cleanWs()
                } catch (err) {
                    echo 'cleanWs no esta disponible'
                }
            }
        }
    }
}