pipeline {
  agent any
  stages {
    stage('docker compose') {
      steps {
        sh 'docker-compose build'
      }
    }

    stage('docker down') {
      steps {
        sh 'docker stack down lovelove'
        sh '''
          until [ -z "$(docker service ls --filter label=com.docker.stack.namespace=lovelove -q)" ] || [ "$limit" -lt 0 ]; do
            sleep 1;
          done

          until [ -z "$(docker network ls --filter label=com.docker.stack.namespace=lovelove -q)" ] || [ "$limit" -lt 0 ]; do
            sleep 1;
          done
        '''
      }
    }

    stage('docker stack') {
      steps {
        sh 'docker stack up -c docker-compose.yml -c docker-compose.yml lovelove'
      }
    }

  }
}